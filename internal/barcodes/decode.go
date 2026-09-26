package barcodes

import (
	"image"
	"image/color"
	"image/draw"
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

// Decode finds every barcode in img, all ten symbologies we generate
// (PDF417 through the ported ZXing reader in internal/pdf417decode). When
// nothing is found it returns an InputError with code ErrNothingFound.
func Decode(img image.Image) ([]Decoded, error) {
	// Attempts from cheap to expensive; the first that finds anything wins.
	// A quiet zone helps codes cropped to their edge, doubling helps dense
	// codes on few pixels, inverting helps light codes on dark ground.
	padded := withQuietZone(img)
	origin := image.Pt(-quietZone, -quietZone) // padded → caller's coordinates
	if found := decodeRegions(padded, origin); len(found) > 0 {
		return found, nil
	}
	if found := decodeOnce(scale2x(padded), origin, 2); len(found) > 0 {
		return found, nil
	}
	if found := decodeOnce(invert(padded), origin, 1); len(found) > 0 {
		return found, nil
	}
	return nil, inputError(ErrNothingFound,
		"no barcode found in the image — check that the code is sharp, fully visible and not too small",
		nil)
}

type namedReader struct {
	format gozxing.BarcodeFormat
	reader gozxing.Reader
}

func readers() []namedReader {
	return []namedReader{
		{gozxing.BarcodeFormat_DATA_MATRIX, datamatrix.NewDataMatrixReader()},
		{gozxing.BarcodeFormat_AZTEC, aztec.NewAztecReader()},
		{gozxing.BarcodeFormat_EAN_13, oned.NewMultiFormatUPCEANReader(nil)},
		{gozxing.BarcodeFormat_CODE_128, oned.NewCode128Reader()},
		{gozxing.BarcodeFormat_CODE_39, oned.NewCode39Reader()},
		{gozxing.BarcodeFormat_ITF, oned.NewITFReader()},
	}
}

// regionScanPixels is the image size from which Decode also searches
// overlapping parts of the image. Detectors such as Aztec's start at the
// image centre and miss a code that sits beside another one.
const regionScanPixels = 400_000

// decodeRegions decodes the whole image and, for large images, five
// overlapping regions (four quadrants and the centre, 60 % each).
func decodeRegions(img image.Image, origin image.Point) []Decoded {
	found := decodeOnce(img, origin, 1)
	b := img.Bounds()
	if b.Dx()*b.Dy() < regionScanPixels {
		return found
	}
	w, h := b.Dx()*3/5, b.Dy()*3/5
	for _, at := range []image.Point{
		{0, 0}, {b.Dx() - w, 0}, {0, b.Dy() - h}, {b.Dx() - w, b.Dy() - h}, {(b.Dx() - w) / 2, (b.Dy() - h) / 2},
	} {
		part := crop(img, image.Rect(at.X, at.Y, at.X+w, at.Y+h))
		found = merge(found, decodeOnce(withQuietZone(part), origin.Add(at).Sub(image.Pt(quietZone, quietZone)), 1))
	}
	return found
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

func crop(img image.Image, r image.Rectangle) image.Image {
	out := image.NewGray(image.Rect(0, 0, r.Dx(), r.Dy()))
	draw.Draw(out, out.Bounds(), img, r.Min, draw.Src)
	return out
}

// decodeOnce runs every reader on one image and collects distinct results.
// QR codes go through the multi reader, so several QR codes in one image
// are all reported. A result point p maps to p/scale + origin in the
// caller's image.
func decodeOnce(img image.Image, origin image.Point, scale int) []Decoded {
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return nil
	}
	hints := map[gozxing.DecodeHintType]interface{}{gozxing.DecodeHintType_TRY_HARDER: true}

	var results []*gozxing.Result
	if qrs, err := multiqr.NewQRCodeMultiReader().DecodeMultiple(bmp, hints); err == nil {
		results = append(results, qrs...)
	} else if r, err := qrcode.NewQRCodeReader().Decode(bmp, hints); err == nil {
		results = append(results, r)
	}
	for _, nr := range readers() {
		if r, err := nr.reader.Decode(bmp, hints); err == nil {
			results = append(results, r)
		}
	}
	if rs, err := pdf417decode.NewReader().DecodeMultiple(bmp, hints); err == nil {
		results = append(results, rs...)
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
			d.Points = append(d.Points, image.Pt(int(p.GetX())/scale, int(p.GetY())/scale).Add(origin))
		}
		out = append(out, d)
	}
	return out
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

func withQuietZone(img image.Image) image.Image {
	b := img.Bounds()
	out := image.NewGray(image.Rect(0, 0, b.Dx()+2*quietZone, b.Dy()+2*quietZone))
	draw.Draw(out, out.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	draw.Draw(out, image.Rect(quietZone, quietZone, quietZone+b.Dx(), quietZone+b.Dy()), img, b.Min, draw.Over)
	return out
}

func scale2x(img image.Image) image.Image {
	b := img.Bounds()
	out := image.NewGray(image.Rect(0, 0, 2*b.Dx(), 2*b.Dy()))
	for y := 0; y < out.Bounds().Dy(); y++ {
		for x := 0; x < out.Bounds().Dx(); x++ {
			out.Set(x, y, img.At(b.Min.X+x/2, b.Min.Y+y/2))
		}
	}
	return out
}

func invert(img image.Image) image.Image {
	b := img.Bounds()
	out := image.NewGray(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			g := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			out.SetGray(x, y, color.Gray{Y: 255 - g.Y})
		}
	}
	return out
}
