package pdf417decode

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"strings"
	"testing"

	"github.com/boombuler/barcode"
	"github.com/makiuchi-d/gozxing"
	encoder "github.com/mlcmcp/mlc_barcode/internal/pdf417"
)

// render encodes text with our (fixed) PDF417 encoder, scales it by 3 and
// adds a white margin.
func render(t *testing.T, text string, level byte) image.Image {
	t.Helper()
	bc, err := encoder.Encode(text, level)
	if err != nil {
		t.Fatalf("encode %q: %v", text, err)
	}
	b := bc.Bounds()
	scaled, err := barcode.Scale(bc, b.Dx()*3, b.Dy()*3)
	if err != nil {
		t.Fatal(err)
	}
	img := image.NewGray(image.Rect(0, 0, scaled.Bounds().Dx()+40, scaled.Bounds().Dy()+40))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	draw.Draw(img, scaled.Bounds().Add(image.Pt(20, 20)), scaled, image.Point{}, draw.Src)
	return img
}

func read(t *testing.T, img image.Image) (string, error) {
	t.Helper()
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		t.Fatal(err)
	}
	res, err := NewReader().Decode(bmp, nil)
	if err != nil {
		return "", err
	}
	return res.GetText(), nil
}

func TestRoundTripASCIIAndDigits(t *testing.T) {
	sources := map[string]string{
		"ascii":  strings.Repeat("abcdefghij", 6),
		"upper":  strings.Repeat("BOARDING PASS LH123 ", 3),
		"digits": strings.Repeat("0123456789", 6),
		"mixed":  strings.Repeat("Ab1;x,y:Z-9 ", 5),
	}
	for name, src := range sources {
		for n := 1; n <= 60; n++ {
			text := strings.TrimSpace(src[:n])
			if text == "" {
				continue
			}
			for _, level := range []byte{2, 5} {
				t.Run(fmt.Sprintf("%s/%d/L%d", name, n, level), func(t *testing.T) {
					got, err := read(t, render(t, text, level))
					if err != nil || got != text {
						t.Errorf("read %q, %v; want %q", got, err, text)
					}
				})
			}
		}
	}
}

// Error correction per ISO/IEC 15438: n EC codewords repair n/2 errors.
func TestErrorCorrectionRepairsDamage(t *testing.T) {
	text := "BOARDING PASS LH123 FRA-JFK SEAT 12A GATE B44 PAX MUSTERMANN/MAX"
	for _, tc := range []struct {
		level   byte
		damaged int
	}{{2, 4}, {5, 24}} {
		bc, err := encoder.Encode(text, tc.level)
		if err != nil {
			t.Fatal(err)
		}
		cols, rows := (bc.Bounds().Dx()-1)/17-4, bc.Bounds().Dy()/2
		img := render(t, text, tc.level).(*image.Gray)
		r := rand.New(rand.NewSource(7))
		for i := 0; i < tc.damaged; i++ {
			col, row := 2+r.Intn(cols), r.Intn(rows) // data columns only
			x0, y0 := 20+col*17*3+6, 20+row*2*3
			for y := y0; y < y0+6; y++ {
				for x := x0; x < x0+17*3-12; x++ {
					img.SetGray(x, y, color.Gray{Y: 255})
				}
			}
		}
		if got, err := read(t, img); err != nil || got != text {
			t.Errorf("level %d with %d damaged codewords: %q, %v", tc.level, tc.damaged, got, err)
		}
	}
}

// Text outside ASCII must come back unchanged (the encoder marks UTF-8 with
// ECI 26; without it readers use ISO-8859-1 and show mojibake).
func TestRoundTripUnicode(t *testing.T) {
	for _, text := range []string{"Größe äöü ß", "Preis: 12 €", "Ünïcödé 😀 ok", "ÄÖÜ" + strings.Repeat("x", 30) + "äöü"} {
		got, err := read(t, render(t, text, 2))
		if err != nil || got != text {
			t.Errorf("read %q, %v; want %q", got, err, text)
		}
	}
}

func rotate90(src image.Image) *image.Gray {
	b := src.Bounds()
	out := image.NewGray(image.Rect(0, 0, b.Dy(), b.Dx()))
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			out.Set(b.Dy()-1-y, x, src.At(b.Min.X+x, b.Min.Y+y))
		}
	}
	return out
}

func TestRotated(t *testing.T) {
	text := "ROTATED LABEL 12345"
	img := render(t, text, 2)
	for turns := 1; turns <= 3; turns++ {
		img = rotate90(img)
		if got, err := read(t, img); err != nil || got != text {
			t.Errorf("%d°: %q, %v", turns*90, got, err)
		}
	}
}

func TestMultipleSymbols(t *testing.T) {
	a, b := render(t, "FIRST LABEL", 2), render(t, "SECOND LABEL", 2)
	sheet := image.NewGray(image.Rect(0, 0, max(a.Bounds().Dx(), b.Bounds().Dx()), a.Bounds().Dy()+b.Bounds().Dy()+40))
	draw.Draw(sheet, sheet.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	draw.Draw(sheet, a.Bounds(), a, image.Point{}, draw.Src)
	draw.Draw(sheet, b.Bounds().Add(image.Pt(0, a.Bounds().Dy()+40)), b, image.Point{}, draw.Src)
	bmp, _ := gozxing.NewBinaryBitmapFromImage(sheet)
	results, err := NewReader().DecodeMultiple(bmp, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]bool{}
	for _, r := range results {
		got[r.GetText()] = true
	}
	if !got["FIRST LABEL"] || !got["SECOND LABEL"] {
		t.Errorf("found %v", got)
	}
}
