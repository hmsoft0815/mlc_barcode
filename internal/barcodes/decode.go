package barcodes

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"sort"
	"strings"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/aztec"
	"github.com/makiuchi-d/gozxing/datamatrix"
	multiqr "github.com/makiuchi-d/gozxing/multi/qrcode"
	"github.com/makiuchi-d/gozxing/oned"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/mlcmcp/mlc_barcode/internal/pdf417decode"
)

// Decoded is one barcode found in an image.
type Decoded struct {
	Type   BarcodeType   // our name, e.g. qr or ean13; the reader's name for formats we do not generate
	Text   string        // content as read (EAN/UPC including the check digit)
	Points []image.Point // corner or finder points in image coordinates
}

// DecodeOptions tune Decode for where the image comes from.
type DecodeOptions struct {
	// Live marks a camera frame: another one follows a moment later, so
	// the expensive search for codes on bottles (decodeCurved) is cut to
	// the most likely shapes of the largest candidate.
	Live bool
	// Surface is what the user said the code sits on (the scan view's
	// shape switch); SurfaceAuto guesses.
	Surface Surface
	// Family is what kind of code the user scans (the scan view's code
	// switch): only those readers run. A camera frame with the 2D readers
	// alone takes a tenth of the time of all readers.
	Family Family
	// Sharpen applies an unsharp mask first — for the front camera, whose
	// fixed focus blurs a code held close. Measured on a blurred
	// DataMatrix it reads one blur step more (sigma 2.0 instead of 1.5 at
	// 120 px), no more: it helps at the margin only.
	Sharpen bool
}

// Family narrows the readers to one kind of code.
type Family int

const (
	FamilyAll    Family = iota
	FamilySquare        // QR, DataMatrix, Aztec (and codes on bottles)
	FamilyWide          // 1D codes and PDF417
)

// ParseFamily maps "square" and "wide" to a Family; anything else is
// FamilyAll.
func ParseFamily(s string) Family {
	switch s {
	case "square":
		return FamilySquare
	case "wide":
		return FamilyWide
	}
	return FamilyAll
}

// firstReaders is the reader set of the first pass and the fallbacks.
func (f Family) firstReaders() readerSet {
	switch f {
	case FamilySquare:
		return squareReaders
	case FamilyWide:
		return wideReaders
	}
	return allReaders
}

// Surface is a hint for the shape under the code. On a bottle the user
// can judge how wide the code is compared with the bottle, q; the code
// then covers s = 2·asin(q) of the circumference (q = ½ → 60°).
type Surface int

const (
	SurfaceAuto   Surface = iota
	SurfaceFlat           // a box or label: no search for curved codes
	SurfaceBottle         // the code covers about half the bottle's width (60°)
	SurfaceTube           // the code wraps far around a thin tube (about 100°)
)

// ParseSurface maps "auto", "flat", "bottle", "tube" to a Surface; anything
// else is SurfaceAuto.
func ParseSurface(s string) Surface {
	switch s {
	case "flat":
		return SurfaceFlat
	case "bottle":
		return SurfaceBottle
	case "tube":
		return SurfaceTube
	}
	return SurfaceAuto
}

// Decode finds every barcode in img, all ten symbologies we generate
// (PDF417 through the ported ZXing reader in internal/pdf417decode). When
// nothing is found it returns an InputError with code ErrNothingFound.
func Decode(img image.Image) ([]Decoded, error) {
	return DecodeWith(img, DecodeOptions{})
}

// DecodeWith is Decode with options.
func DecodeWith(img image.Image, opts DecodeOptions) ([]Decoded, error) {
	found, err := decodeAll(img, opts)
	return dropInside2D(found), err
}

// dropInside2D removes 1D results that lie inside a 2D code read from the
// same image. A row through a QR code's modules can match a short 1D
// pattern: UPC-E, with its weak check, was read from the modules of an
// event QR code — in every version before this check.
func dropInside2D(found []Decoded) []Decoded {
	var areas []image.Rectangle
	for _, d := range found {
		if !isLinear(d.Type) && len(d.Points) > 0 {
			r := image.Rectangle{Min: d.Points[0], Max: d.Points[0]}
			for _, p := range d.Points[1:] {
				r = r.Union(image.Rectangle{Min: p, Max: p.Add(image.Pt(1, 1))})
			}
			// QR points are finder centres, not corners: widen a little.
			m := max(r.Dx(), r.Dy())/5 + 2
			areas = append(areas, r.Inset(-m))
		}
	}
	if len(areas) == 0 {
		return found
	}
	out := found[:0]
	for _, d := range found {
		if isLinear(d.Type) && len(d.Points) > 0 && insideAny(d.Points, areas) {
			continue
		}
		out = append(out, d)
	}
	return out
}

func insideAny(pts []image.Point, areas []image.Rectangle) bool {
	for _, a := range areas {
		inside := true
		for _, p := range pts {
			inside = inside && p.In(a)
		}
		if inside {
			return true
		}
	}
	return false
}

func decodeAll(img image.Image, opts DecodeOptions) ([]Decoded, error) {
	gray := toGray(img)
	b := img.Bounds()
	base := frame{ox: float64(b.Min.X), oy: float64(b.Min.Y), f: 1}

	// Phone photos have far more pixels than any code needs; decoding them
	// at full size took seconds. Work on a reduced copy first.
	work, wf := gray, base
	if k := (max(gray.Rect.Dx(), gray.Rect.Dy()) + maxWorkSide - 1) / maxWorkSide; k > 1 {
		work, wf = shrink(gray, k), base.scaled(1/float64(k))
	}
	if opts.Sharpen {
		work = unsharp(work, 2, 1.5)
	}
	if found := decodeScene(work, wf, opts.Family); len(found) > 0 {
		return found, nil
	}
	// Attempts from cheap to expensive; the first that finds anything wins.
	// Doubling helps dense codes on few pixels, inverting helps light codes
	// on dark ground, unwrapping codes on bottles and tubes, full
	// resolution small codes in a large photo.
	padded, pf := withQuietZone(work), wf.crop(image.Pt(-quietZone, -quietZone))
	set := opts.Family.firstReaders()
	if found := decodeOnce(scale2x(padded), pf.scaled(2), set); len(found) > 0 {
		return found, nil
	}
	if found := decodeOnce(invert(padded), pf, set); len(found) > 0 {
		return found, nil
	}
	if opts.Family != FamilyWide {
		if found := decodeCurved(work, wf, opts); len(found) > 0 {
			return found, nil
		}
	}
	if work != gray {
		if found := decodeScene(gray, base, opts.Family); len(found) > 0 {
			return found, nil
		}
	}
	return nil, inputError(ErrNothingFound,
		"no barcode found in the image — check that the code is sharp, fully visible and not too small",
		nil)
}

// maxWorkSide is the long side Decode first works at. A 12 MP photo is
// reduced to a quarter of its pixels; a code that fills a tenth of its
// width keeps about 200 px.
const maxWorkSide = 2000

// frame maps a point p of the image a reader saw to the caller's image:
// p*f + (ox, oy).
type frame struct{ ox, oy, f float64 }

// crop is the frame of a part of the image that starts at at.
func (fr frame) crop(at image.Point) frame {
	return frame{fr.ox + float64(at.X)*fr.f, fr.oy + float64(at.Y)*fr.f, fr.f}
}

// scaled is the frame of the image enlarged by s.
func (fr frame) scaled(s float64) frame { return frame{fr.ox, fr.oy, fr.f / s} }

// inverse maps a point of the caller's image back into the frame's image.
func (fr frame) inverse(p image.Point) image.Point {
	return image.Pt(int((float64(p.X)-fr.ox)/fr.f), int((float64(p.Y)-fr.oy)/fr.f))
}

func (fr frame) point(x, y float64) image.Point {
	return image.Pt(int(x*fr.f+fr.ox), int(y*fr.f+fr.oy))
}

type namedReader struct {
	format gozxing.BarcodeFormat
	reader gozxing.Reader
}

// readerSet selects what decodeOnce runs.
type readerSet int

const (
	allReaders    readerSet = iota
	matrixReaders           // DataMatrix and Aztec, see decodeScene
	linearReaders           // the 1D readers, see decodeLinearRest
	curvedReaders           // DataMatrix, Aztec and QR, see decodeCurved
	squareReaders           // QR (multi), DataMatrix, Aztec: FamilySquare
	wideReaders             // the 1D readers and PDF417: FamilyWide
)

func readers(set readerSet) []namedReader {
	matrix := []namedReader{
		{gozxing.BarcodeFormat_DATA_MATRIX, datamatrix.NewDataMatrixReader()},
		{gozxing.BarcodeFormat_AZTEC, aztec.NewAztecReader()},
	}
	linear := []namedReader{
		{gozxing.BarcodeFormat_EAN_13, oned.NewMultiFormatUPCEANReader(nil)},
		{gozxing.BarcodeFormat_CODE_128, oned.NewCode128Reader()},
		{gozxing.BarcodeFormat_CODE_39, oned.NewCode39Reader()},
		{gozxing.BarcodeFormat_ITF, oned.NewITFReader()},
	}
	switch set {
	case matrixReaders, curvedReaders, squareReaders:
		return matrix
	case linearReaders, wideReaders:
		return linear
	}
	return append(matrix, linear...)
}

// decodeScene decodes the whole image, then looks for the codes a single
// pass misses. The DataMatrix and Aztec detectors search outwards from the
// centre of what they are given and report one code each: they get a window
// around every spot that looks like a 2D code. The 1D readers report the
// first code they meet: each found one is painted over and the image read
// again. QR and PDF417 readers report all codes in one pass.
func decodeScene(img *image.Gray, fr frame, family Family) []Decoded {
	found := decodeOnce(withQuietZone(img), fr.crop(image.Pt(-quietZone, -quietZone)), family.firstReaders())
	if family == FamilyWide {
		return decodeLinearRest(img, fr, found)
	}
	for _, box := range matrixCandidates(img) {
		// The box covers the busy part; a margin gives the detector the
		// code's edge and quiet zone.
		r := box.Inset(-(max(box.Dx(), box.Dy())/3 + 8)).Intersect(img.Rect)
		if alreadyFound(found, r, fr, matrixReaders) {
			continue
		}
		found = merge(found, decodePart(img, r, fr, matrixReaders))
	}
	if family == FamilySquare {
		return found
	}
	return decodeLinearRest(img, fr, found)
}

// maxLinearCodes bounds the paint-and-read-again loop of decodeLinearRest.
const maxLinearCodes = 12

// decodeLinearRest finds further 1D codes: the ones already found are
// painted white on a copy of img and the 1D readers run again, until a run
// brings nothing new.
func decodeLinearRest(img *image.Gray, fr frame, found []Decoded) []Decoded {
	var work *image.Gray
	masked := map[int]bool{}
	for range maxLinearCodes {
		painted := false
		for i, d := range found {
			if masked[i] || !isLinear(d.Type) || len(d.Points) < 2 {
				continue
			}
			if work == nil {
				work = image.NewGray(img.Rect)
				copy(work.Pix, img.Pix)
			}
			var pts []image.Point
			for _, p := range d.Points {
				pts = append(pts, fr.inverse(p))
			}
			paintOverLinear(work, pts)
			masked[i], painted = true, true
		}
		if !painted {
			break
		}
		n := len(found)
		found = merge(found, decodeOnce(withQuietZone(work), fr.crop(image.Pt(-quietZone, -quietZone)), linearReaders))
		if len(found) == n {
			break
		}
	}
	return found
}

// isLinear reports whether t is a 1D code — also formats we do not
// generate but the readers report, such as UPC_E.
func isLinear(t BarcodeType) bool {
	switch t {
	case TypeQR, TypeDataMatrix, TypeAztec, TypePDF417:
		return false
	}
	return true
}

// paintOverLinear whites out the 1D code whose result points (on the line
// that was read) are pts. The bars run across that line; the code reaches
// as far up and down (left and right for a turned code) as the pixels
// still match the line that was read.
func paintOverLinear(img *image.Gray, pts []image.Point) {
	r := image.Rectangle{Min: pts[0], Max: pts[0]}
	for _, p := range pts[1:] {
		r = r.Union(image.Rectangle{Min: p, Max: p.Add(image.Pt(1, 1))})
	}
	r = r.Intersect(img.Rect)
	if r.Empty() {
		return
	}
	pad := max(r.Dx(), r.Dy())/12 + 4
	horizontal := r.Dx() >= r.Dy()
	line := func(i int) []uint8 { // the pixels across the bars at offset i
		out := []uint8{}
		if horizontal {
			if y := r.Min.Y + i; y >= img.Rect.Min.Y && y < img.Rect.Max.Y {
				for x := r.Min.X; x < r.Max.X; x++ {
					out = append(out, img.Pix[img.PixOffset(x, y)])
				}
			}
		} else if x := r.Min.X + i; x >= img.Rect.Min.X && x < img.Rect.Max.X {
			for y := r.Min.Y; y < r.Max.Y; y++ {
				out = append(out, img.Pix[img.PixOffset(x, y)])
			}
		}
		return out
	}
	ref := line(0)
	lo, hi := 0, 0
	for lo-1 >= -img.Rect.Dy() && similar(ref, line(lo-1)) {
		lo--
	}
	for hi+1 <= img.Rect.Dy() && similar(ref, line(hi+1)) {
		hi++
	}
	if horizontal {
		r = image.Rect(r.Min.X-pad, r.Min.Y+lo-pad, r.Max.X+pad, r.Min.Y+hi+1+pad)
	} else {
		r = image.Rect(r.Min.X+lo-pad, r.Min.Y-pad, r.Min.X+hi+1+pad, r.Max.Y+pad)
	}
	r = r.Intersect(img.Rect)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		row := img.Pix[img.PixOffset(r.Min.X, y):][:r.Dx()]
		for i := range row {
			row[i] = 0xff
		}
	}
}

// similar reports whether two lines across the bars show the same code.
func similar(a, b []uint8) bool {
	if len(a) == 0 || len(a) != len(b) {
		return false
	}
	sum := 0
	for i := range a {
		sum += absDiff(a[i], b[i])
	}
	return sum/len(a) < lineMatch
}

// lineMatch is the mean grey difference up to which two lines count as
// the same bars.
const lineMatch = 40

// alreadyFound reports whether a code the readers of set return lies
// inside r already, so the window would only find it again.
func alreadyFound(found []Decoded, r image.Rectangle, fr frame, set readerSet) bool {
	caller := image.Rectangle{Min: fr.point(float64(r.Min.X), float64(r.Min.Y)), Max: fr.point(float64(r.Max.X), float64(r.Max.Y))}
	for _, d := range found {
		if (d.Type == TypeDataMatrix || d.Type == TypeAztec) != (set == matrixReaders) || len(d.Points) == 0 {
			continue
		}
		inside := true
		for _, p := range d.Points {
			inside = inside && p.In(caller)
		}
		if inside {
			return true
		}
	}
	return false
}

// Candidate search parameters. The image is divided into cells and the
// dark/light changes along rows and along columns counted in each. A 2D
// code changes often in both directions; plain ground, ruled lines and the
// bars of a 1D code change in one direction only.
const (
	cellsPerSide     = 160 // along the long side
	minCell          = 6
	edgeContrast     = 40   // grey step that counts as a dark/light change
	busyDensity      = 0.04 // changes per pixel, in each direction
	maxCandidates    = 40
	joinCells        = 2 // busy cells this far apart belong to one area
	minCandidateArea = 4 // cells
)

// matrixCandidates returns the bounding boxes of the areas that look like
// 2D codes, largest first.
func matrixCandidates(img *image.Gray) []image.Rectangle {
	b := img.Rect
	cell := max(max(b.Dx(), b.Dy())/cellsPerSide, minCell)
	cw, ch := b.Dx()/cell, b.Dy()/cell
	if cw < 2 || ch < 2 {
		return nil
	}
	busy := make([]bool, cw*ch)
	for cy := 0; cy < ch; cy++ {
		for cx := 0; cx < cw; cx++ {
			var hx, vy int
			for y := cy * cell; y < (cy+1)*cell; y++ {
				row := img.Pix[img.PixOffset(b.Min.X, b.Min.Y+y):]
				next := row
				if y+1 < b.Dy() {
					next = img.Pix[img.PixOffset(b.Min.X, b.Min.Y+y+1):]
				}
				for x := cx * cell; x < (cx+1)*cell; x++ {
					if x+1 < b.Dx() && absDiff(row[x], row[x+1]) > edgeContrast {
						hx++
					}
					if absDiff(row[x], next[x]) > edgeContrast {
						vy++
					}
				}
			}
			n := float64(cell * cell)
			busy[cy*cw+cx] = float64(hx)/n >= busyDensity && float64(vy)/n >= busyDensity
		}
	}
	return areaBoxes(busy, cw, ch, cell, b.Min)
}

// areaBoxes returns the pixel bounding boxes of the connected areas of
// marked cells, largest first. Inside a code some cells show no change (a
// large module, Aztec's bullseye, a wide bar); the marks are grown first so
// the parts of one code join.
func areaBoxes(marked []bool, cw, ch, cell int, origin image.Point) []image.Rectangle {
	marked = dilate(marked, cw, ch, joinCells)
	seen := make([]bool, len(marked))
	type area struct {
		r     image.Rectangle
		cells int
	}
	var areas []area
	var stack []int
	for i, m := range marked {
		if !m || seen[i] {
			continue
		}
		a := area{r: image.Rect(i%cw, i/cw, i%cw+1, i/cw+1)}
		seen[i] = true
		stack = append(stack[:0], i)
		for len(stack) > 0 {
			j := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			x, y := j%cw, j/cw
			a.cells++
			a.r = a.r.Union(image.Rect(x, y, x+1, y+1))
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					nx, ny := x+dx, y+dy
					if nx < 0 || ny < 0 || nx >= cw || ny >= ch {
						continue
					}
					if k := ny*cw + nx; marked[k] && !seen[k] {
						seen[k] = true
						stack = append(stack, k)
					}
				}
			}
		}
		if a.cells >= minCandidateArea {
			areas = append(areas, a)
		}
	}
	sort.Slice(areas, func(i, j int) bool { return areas[i].cells > areas[j].cells })
	var boxes []image.Rectangle
	for _, a := range areas[:min(len(areas), maxCandidates)] {
		boxes = append(boxes, image.Rect(a.r.Min.X*cell, a.r.Min.Y*cell, a.r.Max.X*cell, a.r.Max.Y*cell).Add(origin))
	}
	return boxes
}

// dilate marks every cell within r cells of a marked one.
func dilate(m []bool, w, h, r int) []bool {
	out := make([]bool, len(m))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if !m[y*w+x] {
				continue
			}
			for ny := max(0, y-r); ny <= min(h-1, y+r); ny++ {
				for nx := max(0, x-r); nx <= min(w-1, x+r); nx++ {
					out[ny*w+nx] = true
				}
			}
		}
	}
	return out
}

func absDiff(a, b uint8) int {
	if a > b {
		return int(a - b)
	}
	return int(b - a)
}

// decodePart decodes the rectangle r of img, with a quiet zone of its own.
func decodePart(img *image.Gray, r image.Rectangle, fr frame, set readerSet) []Decoded {
	part := img.SubImage(r).(*image.Gray)
	return decodeOnce(withQuietZone(part), fr.crop(r.Min.Sub(image.Pt(quietZone, quietZone))), set)
}

// merge appends the codes of more that are not in found yet.
func merge(found, more []Decoded) []Decoded {
	for _, m := range more {
		dup := false
		for _, f := range found {
			if f.Type == m.Type && f.Text == m.Text {
				dup = true
				break
			}
		}
		if !dup {
			found = append(found, m)
		}
	}
	return found
}

// pointMapper maps a result point of the image a reader saw to the
// caller's image: a frame for crops and scalings, an unwrapped cylinder
// (cylinder.go) for codes on bottles.
type pointMapper interface {
	point(x, y float64) image.Point
}

// decodeOnce runs the readers of set on one image and collects distinct
// results; allReaders adds QR (the multi reader, so several QR codes are
// all reported) and PDF417, curvedReaders adds the single QR reader.
// Result points are mapped through fr.
func decodeOnce(img image.Image, fr pointMapper, set readerSet) []Decoded {
	bmp, err := binaryBitmap(img, set)
	if err != nil {
		return nil
	}
	hints := map[gozxing.DecodeHintType]interface{}{gozxing.DecodeHintType_TRY_HARDER: true}

	var results []*gozxing.Result
	if set == allReaders || set == squareReaders {
		if qrs, err := multiqr.NewQRCodeMultiReader().DecodeMultiple(bmp, hints); err == nil {
			results = append(results, qrs...)
		} else if r, err := qrcode.NewQRCodeReader().Decode(bmp, hints); err == nil {
			results = append(results, r)
		}
	}
	if set == curvedReaders {
		if r, err := qrcode.NewQRCodeReader().Decode(bmp, hints); err == nil {
			results = append(results, r)
		}
	}
	for _, nr := range readers(set) {
		if r, err := nr.reader.Decode(bmp, hints); err == nil {
			results = append(results, r)
		}
	}
	if set == allReaders || set == wideReaders {
		if rs, err := pdf417decode.NewReader().DecodeMultiple(bmp, hints); err == nil {
			results = append(results, rs...)
		}
	}

	seen := map[string]bool{}
	var out []Decoded
	for _, r := range results {
		d := Decoded{Type: typeOfFormat(r.GetBarcodeFormat()), Text: r.GetText()}
		if d.Type == TypeEAN13 && strings.HasPrefix(d.Text, "0") {
			// A UPC-A is an EAN-13 with a leading 0, the same symbol; report
			// it the way it was most likely made, as the 12-digit UPC-A.
			d.Type, d.Text = TypeUPCA, d.Text[1:]
		}
		key := string(d.Type) + "\x00" + d.Text
		if seen[key] {
			continue
		}
		seen[key] = true
		for _, p := range r.GetResultPoints() {
			d.Points = append(d.Points, fr.point(p.GetX(), p.GetY()))
		}
		out = append(out, d)
	}
	return out
}

// binaryBitmap hands img to the readers. gozxing's generic conversion reads
// every pixel through image.At; a grey image is already the luminance
// plane and is passed as it is — but that source cannot rotate, which the
// 1D readers need for vertical codes, so it is used only for 2D reader
// sets.
func binaryBitmap(img image.Image, set readerSet) (*gozxing.BinaryBitmap, error) {
	if g, ok := img.(*image.Gray); ok && (set == matrixReaders || set == curvedReaders || set == squareReaders) &&
		g.Rect.Min == (image.Point{}) && g.Stride == g.Rect.Dx() {
		w, h := g.Rect.Dx(), g.Rect.Dy()
		if src, err := gozxing.NewPlanarYUVLuminanceSource(g.Pix, w, h, 0, 0, w, h, false); err == nil {
			return gozxing.NewBinaryBitmap(gozxing.NewHybridBinarizer(src))
		}
	}
	return gozxing.NewBinaryBitmapFromImage(img)
}

func typeOfFormat(f gozxing.BarcodeFormat) BarcodeType {
	switch f {
	case gozxing.BarcodeFormat_QR_CODE:
		return TypeQR
	case gozxing.BarcodeFormat_DATA_MATRIX:
		return TypeDataMatrix
	case gozxing.BarcodeFormat_AZTEC:
		return TypeAztec
	case gozxing.BarcodeFormat_PDF_417:
		return TypePDF417
	case gozxing.BarcodeFormat_EAN_13:
		return TypeEAN13
	case gozxing.BarcodeFormat_EAN_8:
		return TypeEAN8
	case gozxing.BarcodeFormat_UPC_A:
		return TypeUPCA
	case gozxing.BarcodeFormat_CODE_128:
		return TypeCode128
	case gozxing.BarcodeFormat_CODE_39:
		return TypeCode39
	case gozxing.BarcodeFormat_ITF:
		return TypeITF
	}
	return BarcodeType(f.String())
}

// quietZone is the white border added around the image; result points are
// shifted back by it.
const quietZone = 40

func withQuietZone(img *image.Gray) *image.Gray {
	b := img.Rect
	out := image.NewGray(image.Rect(0, 0, b.Dx()+2*quietZone, b.Dy()+2*quietZone))
	for i := range out.Pix {
		out.Pix[i] = 0xff
	}
	for y := 0; y < b.Dy(); y++ {
		copy(out.Pix[(y+quietZone)*out.Stride+quietZone:], img.Pix[img.PixOffset(b.Min.X, b.Min.Y+y):][:b.Dx()])
	}
	return out
}

// toGray converts img to 8-bit luminance at origin (0,0). JPEG photos
// arrive as YCbCr, whose Y plane already is the luminance.
func toGray(img image.Image) *image.Gray {
	b := img.Bounds()
	out := image.NewGray(image.Rect(0, 0, b.Dx(), b.Dy()))
	switch src := img.(type) {
	case *image.Gray:
		for y := 0; y < b.Dy(); y++ {
			copy(out.Pix[y*out.Stride:][:b.Dx()], src.Pix[src.PixOffset(b.Min.X, b.Min.Y+y):])
		}
	case *image.YCbCr:
		for y := 0; y < b.Dy(); y++ {
			copy(out.Pix[y*out.Stride:][:b.Dx()], src.Y[src.YOffset(b.Min.X, b.Min.Y+y):])
		}
	default:
		// Transparent pixels count as white paper, not black (draw.Over
		// onto white, as a viewer would show them).
		draw.Draw(out, out.Rect, image.NewUniform(color.White), image.Point{}, draw.Src)
		draw.Draw(out, out.Rect, img, b.Min, draw.Over)
	}
	return out
}

// shrink reduces img by the integer factor k, averaging k×k blocks.
func shrink(img *image.Gray, k int) *image.Gray {
	w, h := img.Rect.Dx()/k, img.Rect.Dy()/k
	out := image.NewGray(image.Rect(0, 0, w, h))
	sums := make([]int, w)
	for y := 0; y < h; y++ {
		clear(sums)
		for dy := 0; dy < k; dy++ {
			row := img.Pix[(y*k+dy)*img.Stride:]
			for x := 0; x < w; x++ {
				for dx := 0; dx < k; dx++ {
					sums[x] += int(row[x*k+dx])
				}
			}
		}
		for x := 0; x < w; x++ {
			out.Pix[y*out.Stride+x] = uint8(sums[x] / (k * k))
		}
	}
	return out
}

func scale2x(img *image.Gray) *image.Gray {
	b := img.Rect
	out := image.NewGray(image.Rect(0, 0, 2*b.Dx(), 2*b.Dy()))
	for y := 0; y < out.Rect.Dy(); y++ {
		for x := 0; x < out.Rect.Dx(); x++ {
			out.Pix[y*out.Stride+x] = img.Pix[img.PixOffset(b.Min.X+x/2, b.Min.Y+y/2)]
		}
	}
	return out
}

func invert(img *image.Gray) *image.Gray {
	out := image.NewGray(img.Rect)
	for i, v := range img.Pix {
		out.Pix[i] = 255 - v
	}
	return out
}

// unsharp sharpens img: the difference to a Gaussian blur of radius sigma,
// times amount, is added back.
func unsharp(img *image.Gray, sigma, amount float64) *image.Gray {
	blurred := gaussian(img, sigma)
	out := image.NewGray(img.Rect)
	for i, v := range img.Pix {
		d := float64(v) + amount*(float64(v)-float64(blurred.Pix[i]))
		out.Pix[i] = uint8(min(max(d, 0), 255))
	}
	return out
}

// gaussian blurs img (separable kernel, edges repeated).
func gaussian(img *image.Gray, sigma float64) *image.Gray {
	r := int(math.Ceil(sigma * 3))
	k := make([]float64, 2*r+1)
	sum := 0.0
	for i := range k {
		x := float64(i - r)
		k[i] = math.Exp(-x * x / (2 * sigma * sigma))
		sum += k[i]
	}
	for i := range k {
		k[i] /= sum
	}
	w, h := img.Rect.Dx(), img.Rect.Dy()
	tmp := make([]float64, w*h)
	for y := 0; y < h; y++ {
		row := img.Pix[img.PixOffset(img.Rect.Min.X, img.Rect.Min.Y+y):]
		for x := 0; x < w; x++ {
			v := 0.0
			for i := -r; i <= r; i++ {
				v += k[i+r] * float64(row[min(max(x+i, 0), w-1)])
			}
			tmp[y*w+x] = v
		}
	}
	out := image.NewGray(img.Rect)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			v := 0.0
			for i := -r; i <= r; i++ {
				v += k[i+r] * tmp[min(max(y+i, 0), h-1)*w+x]
			}
			out.Pix[out.PixOffset(img.Rect.Min.X+x, img.Rect.Min.Y+y)] = uint8(min(max(v, 0), 255))
		}
	}
	return out
}
