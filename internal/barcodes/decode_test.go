package barcodes

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strings"
	"testing"
)

// roundTrip generates a code as PNG, reads it back and returns what the
// decoder found.
func roundTrip(t *testing.T, btype BarcodeType, data string, opts BarcodeOptions) []Decoded {
	t.Helper()
	img := renderPNG(t, btype, data, opts)
	found, err := Decode(img)
	if err != nil {
		t.Fatalf("%s %q: %v", btype, data, err)
	}
	return found
}

func renderPNG(t *testing.T, btype BarcodeType, data string, opts BarcodeOptions) image.Image {
	t.Helper()
	b, err := GeneratePNG(btype, data, opts)
	if err != nil {
		t.Fatalf("generate %s %q: %v", btype, data, err)
	}
	img, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return img
}

func expectOne(t *testing.T, found []Decoded, btype BarcodeType, text string) {
	t.Helper()
	for _, d := range found {
		if d.Type == btype && d.Text == text {
			return
		}
	}
	t.Errorf("want %s %q, decoded %+v", btype, text, found)
}

// Every symbology we generate must read
// back to exactly the input — EAN/UPC with the completed check digit.
func TestRoundTripAllSymbologies(t *testing.T) {
	cases := []struct {
		btype      BarcodeType
		data, want string
	}{
		{TypeQR, "https://mlcgo.eu/produkte?x=Größe", ""},
		{TypeDataMatrix, "MLC-DM-12345 äöü", ""},
		{TypeAztec, "TICKET-ICE-599-FRA-MUC", ""},
		{TypePDF417, "BOARDING PASS LH123 FRA-JFK", ""},
		{TypePDF417, "Größe äöü – 12 €", ""},
		{TypeCode128, "MLC-128-abc", ""},
		{TypeCode128, "Hello World 123", ""},
		{TypeCode39, "CODE39-TEST", ""},
		{TypeEAN13, "4006381333931", ""},
		{TypeEAN13, "400638133393", "4006381333931"},
		{TypeEAN8, "9638507", "96385074"},
		{TypeUPCA, "036000291452", ""},
		{TypeUPCA, "03600029145", "036000291452"},
		{TypeITF, "12345678", ""},
		{TypeITF, "0012345678", ""},
	}
	for _, c := range cases {
		want := c.want
		if want == "" {
			want = c.data
		}
		expectOne(t, roundTrip(t, c.btype, c.data, DefaultOptions(c.btype)), c.btype, want)
	}
}

// The 2D sweep that found the PDF417 encoder bug, now against our own reader.
func TestRoundTrip2DSweep(t *testing.T) {
	sources := map[string]string{
		"ascii":  strings.Repeat("abcdefghij", 6),
		"umlaut": strings.Repeat("äöüßÄÖÜ", 9),
		"digits": strings.Repeat("0123456789", 6),
		"mixed":  strings.Repeat("Größe 12 € x;y,z", 4),
	}
	for _, btype := range []BarcodeType{TypeQR, TypeDataMatrix, TypeAztec, TypePDF417} {
		for name, src := range sources {
			runes := []rune(src)
			for n := 1; n <= 60; n++ {
				data := strings.TrimSpace(string(runes[:n]))
				if data == "" {
					continue
				}
				t.Run(fmt.Sprintf("%s/%s/%d", btype, name, n), func(t *testing.T) {
					t.Parallel()
					if btype == TypeAztec && strings.ContainsRune(data, '€') {
						// outside ISO-8859-1: rejected with a helpful message
						_, err := Generate(btype, data, DefaultOptions(btype))
						if ie, ok := AsInputError(err); !ok || ie.Code != ErrAztecCharset {
							t.Errorf("want ErrAztecCharset, got %v", err)
						}
						return
					}
					expectOne(t, roundTrip(t, btype, data, DefaultOptions(btype)), btype, data)
				})
			}
		}
	}
}

func TestRoundTripStyling(t *testing.T) {
	caption := DefaultOptions(TypeEAN13)
	caption.ShowText, caption.CustomText, caption.FontSize = true, "Schraube M4", 28
	expectOne(t, roundTrip(t, TypeEAN13, "4006381333931", caption), TypeEAN13, "4006381333931")

	colored := DefaultOptions(TypeQR)
	colored.ForegroundColor, colored.BackgroundColor = "#1b5e20", "#fff59d"
	expectOne(t, roundTrip(t, TypeQR, "farbig", colored), TypeQR, "farbig")

	transparent := DefaultOptions(TypeDataMatrix)
	transparent.BackgroundColor = "transparent"
	expectOne(t, roundTrip(t, TypeDataMatrix, "transparent", transparent), TypeDataMatrix, "transparent")
}

func TestDecodeLightOnDark(t *testing.T) {
	opts := DefaultOptions(TypeQR)
	opts.ForegroundColor, opts.BackgroundColor = "#ffffff", "#000000"
	expectOne(t, roundTrip(t, TypeQR, "hell auf dunkel", opts), TypeQR, "hell auf dunkel")
}

func TestDecodeSeveralCodesInOneImage(t *testing.T) {
	a := renderPNG(t, TypeQR, "erster", DefaultOptions(TypeQR))
	b := renderPNG(t, TypeQR, "zweiter", DefaultOptions(TypeQR))
	sheet := image.NewRGBA(image.Rect(0, 0, a.Bounds().Dx()+b.Bounds().Dx()+60, a.Bounds().Dy()))
	draw.Draw(sheet, sheet.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	draw.Draw(sheet, a.Bounds(), a, image.Point{}, draw.Src)
	draw.Draw(sheet, b.Bounds().Add(image.Pt(a.Bounds().Dx()+60, 0)), b, image.Point{}, draw.Src)

	found, err := Decode(sheet)
	if err != nil {
		t.Fatal(err)
	}
	expectOne(t, found, TypeQR, "erster")
	expectOne(t, found, TypeQR, "zweiter")
}

func TestDecodeNothingFound(t *testing.T) {
	blank := image.NewGray(image.Rect(0, 0, 200, 200))
	draw.Draw(blank, blank.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	_, err := Decode(blank)
	if ie, ok := AsInputError(err); !ok || ie.Code != ErrNothingFound {
		t.Errorf("want ErrNothingFound, got %v", err)
	}
}

// A sheet with different symbologies side by side: Aztec is only found by
// the region scan, because its detector starts at the image centre.
func TestDecodeMixedSheet(t *testing.T) {
	parts := []image.Image{
		renderPNG(t, TypeQR, "QR auf dem Blatt", DefaultOptions(TypeQR)),
		renderPNG(t, TypeAztec, "AZTEC AUF DEM BLATT", DefaultOptions(TypeAztec)),
		renderPNG(t, TypeEAN13, "4006381333931", DefaultOptions(TypeEAN13)),
	}
	sheet := image.NewRGBA(image.Rect(0, 0, 1400, 900))
	draw.Draw(sheet, sheet.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	at := []image.Point{{60, 60}, {800, 60}, {400, 600}}
	for i, p := range parts {
		draw.Draw(sheet, p.Bounds().Add(at[i]), p, image.Point{}, draw.Src)
	}

	found, err := Decode(sheet)
	if err != nil {
		t.Fatal(err)
	}
	expectOne(t, found, TypeQR, "QR auf dem Blatt")
	expectOne(t, found, TypeAztec, "AZTEC AUF DEM BLATT")
	expectOne(t, found, TypeEAN13, "4006381333931")

	// Points are in sheet coordinates: the QR's lie inside its square.
	qrBox := parts[0].Bounds().Add(at[0])
	for _, d := range found {
		if d.Type != TypeQR {
			continue
		}
		for _, p := range d.Points {
			if !p.In(qrBox) {
				t.Errorf("QR point %v outside %v", p, qrBox)
			}
		}
	}
}
