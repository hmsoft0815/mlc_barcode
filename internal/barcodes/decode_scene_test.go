package barcodes

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"testing"
	"time"
)

// scene is a photo-like image: a tinted background with printed lines,
// like the side of a pack, with codes pasted at given positions.
func scene(t testing.TB, w, h int, codes []image.Image, at []image.Point) *image.RGBA {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.RGBA{230, 225, 215, 255}), image.Point{}, draw.Src)
	for y := 0; y < h; y += 40 {
		draw.Draw(img, image.Rect(0, y, w, y+6), image.NewUniform(color.RGBA{70, 70, 70, 255}), image.Point{}, draw.Src)
	}
	for i, c := range codes {
		draw.Draw(img, c.Bounds().Add(at[i]), c, image.Point{}, draw.Src)
	}
	return img
}

func pharmaCode(i int) string {
	return fmt.Sprintf("[)>\x1e06\x1d9N111234567842\x1d1TCH%02d\x1dD280331\x1dSSN%04d\x1e\x04", i, i)
}

func codeImage(t testing.TB, btype BarcodeType, data string, size int) image.Image {
	t.Helper()
	opts := DefaultOptions(btype)
	if btype == TypePDF417 || btype == TypeEAN13 {
		opts.Width, opts.Height = size*2, size*2/3
	} else {
		opts.Width, opts.Height = size, size
	}
	return renderPNG(t.(*testing.T), btype, data, opts)
}

// A single code must be found wherever it sits in the photo (B-20260926-05:
// DataMatrix and Aztec used to be found only near the centre).
func TestDecodeAnywhereInImage(t *testing.T) {
	const w, h = 1400, 1000
	samples := map[BarcodeType]string{
		TypeDataMatrix: pharmaCode(1), TypeAztec: "TICKET-ICE-599-FRA-MUC", TypeQR: "https://mlcgo.eu",
		TypePDF417: "BOARDING PASS LH123", TypeEAN13: "4006381333931",
	}
	for btype, data := range samples {
		code := codeImage(t, btype, data, 180)
		cb := code.Bounds()
		for _, pos := range []image.Point{
			{40, 60}, {(w - cb.Dx()) / 2, 60}, {w - cb.Dx() - 40, 60},
			{40, (h - cb.Dy()) / 2}, {(w - cb.Dx()) / 2, (h - cb.Dy()) / 2}, {w - cb.Dx() - 40, (h - cb.Dy()) / 2},
			{40, h - cb.Dy() - 60}, {(w - cb.Dx()) / 2, h - cb.Dy() - 60}, {w - cb.Dx() - 40, h - cb.Dy() - 60},
		} {
			t.Run(fmt.Sprintf("%s/%d,%d", btype, pos.X, pos.Y), func(t *testing.T) {
				t.Parallel()
				found, _ := Decode(scene(t, w, h, []image.Image{code}, []image.Point{pos}))
				expectOne(t, found, btype, data)
			})
		}
	}
}

// Several packs in one photo.
func TestDecodeSeveralPacks(t *testing.T) {
	for _, n := range []int{3, 4, 6} {
		t.Run(fmt.Sprintf("%d packs", n), func(t *testing.T) {
			t.Parallel()
			var codes []image.Image
			var at []image.Point
			for i := 0; i < n; i++ {
				codes = append(codes, codeImage(t, TypeDataMatrix, pharmaCode(i), 170))
				at = append(at, image.Pt(60+(i%3)*450, 80+(i/3)*450))
			}
			found, _ := Decode(scene(t, 1400, 1000, codes, at))
			for i := 0; i < n; i++ {
				expectOne(t, found, TypeDataMatrix, pharmaCode(i))
			}
		})
	}
}

// A shipping label: two Code 128 and a DataMatrix. 1D readers report one
// code per pass; the second Code 128 used to go unnoticed.
func TestDecodeShippingLabel(t *testing.T) {
	codes := []struct {
		btype BarcodeType
		data  string
	}{{TypeCode128, "TBA123456789000"}, {TypeCode128, "DE-LABEL-42"}, {TypeDataMatrix, "AMZN-XYZ-0001"}}
	var imgs []image.Image
	for _, c := range codes {
		imgs = append(imgs, codeImage(t, c.btype, c.data, 200))
	}
	found, _ := Decode(scene(t, 1400, 1000, imgs, []image.Point{{100, 100}, {100, 600}, {900, 300}}))
	for _, c := range codes {
		expectOne(t, found, c.btype, c.data)
	}
}

// A phone photo is 12 MP; decoding must stay interactive.
func TestDecodePhonePhotoTime(t *testing.T) {
	if testing.Short() {
		t.Skip("timing test")
	}
	code := codeImage(t, TypeDataMatrix, pharmaCode(7), 500)
	img := scene(t, 4000, 3000, []image.Image{code}, []image.Point{{600, 700}})
	start := time.Now()
	found, _ := Decode(img)
	elapsed := time.Since(start)
	t.Logf("12 MP photo: %d found in %v", len(found), elapsed)
	expectOne(t, found, TypeDataMatrix, pharmaCode(7))
	if elapsed > 2*time.Second {
		t.Errorf("too slow: %v", elapsed)
	}
}

// A camera frame as the GUI sends it (1600 px long side) is decoded about
// three times a second; it must not keep the loop waiting.
func BenchmarkDecodeCameraFrame(b *testing.B) {
	t := &testing.T{}
	code := codeImage(t, TypeDataMatrix, pharmaCode(3), 160)
	img := scene(b, 1600, 900, []image.Image{code}, []image.Point{{300, 200}})
	for b.Loop() {
		if found, _ := Decode(img); len(found) != 1 {
			b.Fatalf("found %d", len(found))
		}
	}
}

// The scan view's code switch runs only one kind of readers: square codes
// (QR, DataMatrix, Aztec) or wide ones (1D, PDF417). Each finds its own
// kind and ignores the other.
func TestDecodeFamily(t *testing.T) {
	dm := codeImage(t, TypeDataMatrix, pharmaCode(4), 200)
	ean := codeImage(t, TypeEAN13, "4006381333931", 200)
	img := scene(t, 1400, 900, []image.Image{dm, ean}, []image.Point{{150, 200}, {800, 300}})
	for _, c := range []struct {
		family      Family
		want, never BarcodeType
	}{{FamilySquare, TypeDataMatrix, TypeEAN13}, {FamilyWide, TypeEAN13, TypeDataMatrix}} {
		found, _ := DecodeWith(img, DecodeOptions{Live: true, Family: c.family})
		has := map[BarcodeType]bool{}
		for _, d := range found {
			has[d.Type] = true
		}
		if !has[c.want] || has[c.never] {
			t.Errorf("family %d: want %s and no %s, got %+v", c.family, c.want, c.never, found)
		}
	}
	found, _ := DecodeWith(img, DecodeOptions{Live: true})
	if len(found) != 2 {
		t.Errorf("all readers: want both codes, got %+v", found)
	}
}

// The front camera's fixed focus blurs a code held close. Sharpening reads
// one blur step more — measured: a 120 px DataMatrix at blur sigma 2.0 is
// read only sharpened.
func TestDecodeSharpenFrontCamera(t *testing.T) {
	code := codeImage(t, TypeDataMatrix, pharmaCode(5), 120)
	blurred := gaussian(toGray(scene(t, 900, 600, []image.Image{code}, []image.Point{{300, 150}})), 2)
	plain, _ := DecodeWith(blurred, DecodeOptions{Live: true, Family: FamilySquare})
	sharp, _ := DecodeWith(blurred, DecodeOptions{Live: true, Family: FamilySquare, Sharpen: true})
	if len(sharp) != 1 {
		t.Errorf("sharpened: want the DataMatrix, got %+v", sharp)
	}
	if len(plain) != 0 {
		t.Logf("note: read without sharpening too (%d) — the case is no longer on the edge", len(plain))
	}
}
