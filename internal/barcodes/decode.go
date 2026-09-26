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
)

// Decoded is one barcode found in an image.
type Decoded struct {
	Type   BarcodeType   // our name, e.g. qr or ean13; the reader's name for formats we do not generate
	Text   string        // content as read (EAN/UPC including the check digit)
	Points []image.Point // corner or finder points in image coordinates
}

// Decode finds every barcode in img. PDF417 cannot be read yet (see
// PLAN.md, decoder phase 4). When nothing is found it returns an
// InputError with code ErrNothingFound.
func Decode(img image.Image) ([]Decoded, error) {
	// Attempts from cheap to expensive; the first that finds anything wins.
	// A quiet zone helps codes cropped to their edge, doubling helps dense
	// codes on few pixels, inverting helps light codes on dark ground.
	padded := withQuietZone(img)
	attempts := []func() image.Image{
		func() image.Image { return padded },
		func() image.Image { return scale2x(padded) },
		func() image.Image { return invert(padded) },
	}
	for _, next := range attempts {
		if found := decodeOnce(next()); len(found) > 0 {
			return found, nil
		}
	}
	return nil, inputError(ErrNothingFound,
		"no barcode found in the image — check that the code is sharp, fully visible and not too small; PDF417 cannot be read yet",
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

// decodeOnce runs every reader on one image and collects distinct results.
// QR codes go through the multi reader, so several QR codes in one image
// are all reported.
func decodeOnce(img image.Image) []Decoded {
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
			d.Points = append(d.Points, image.Pt(int(p.GetX())-quietZone, int(p.GetY())-quietZone))
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
