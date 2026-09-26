package barcodes

import (
	"image"
	"testing"
)

// The decoder pads images with its own quiet zone, which hid that our
// output had none. Here every code is read exactly as exported.
func TestExportedCodesReadWithoutExtraMargin(t *testing.T) {
	for _, c := range []struct {
		btype BarcodeType
		data  string
	}{
		{TypeQR, "https://mlcgo.eu"}, {TypeDataMatrix, "MLC-DM-12345"}, {TypeAztec, "TICKET-ICE-599"}, {TypePDF417, "BOARDING PASS"},
		{TypeCode128, "MLC-128-abc"}, {TypeCode39, "CODE39"}, {TypeEAN13, "4006381333931"},
		{TypeEAN8, "96385074"}, {TypeUPCA, "036000291452"}, {TypeITF, "12345678"},
	} {
		img := renderPNG(t, c.btype, c.data, DefaultOptions(c.btype))
		found := decodeOnce(img, image.Point{}, 1) // no withQuietZone
		ok := false
		for _, d := range found {
			ok = ok || d.Type == c.btype
		}
		if !ok {
			t.Errorf("%s: not readable as exported (found %+v)", c.btype, found)
		}
	}
}

// The requested size is the total size; the quiet zone is taken from it.
func TestQuietZoneKeepsRequestedSize(t *testing.T) {
	for _, btype := range []BarcodeType{TypeQR, TypeEAN13, TypePDF417} {
		opts := DefaultOptions(btype)
		img := renderPNG(t, btype, map[BarcodeType]string{TypeQR: "x", TypeEAN13: "4006381333931", TypePDF417: "BOARDING"}[btype], opts)
		if got := img.Bounds().Size(); got != image.Pt(opts.Width, opts.Height) {
			t.Errorf("%s: size %v, want %dx%d", btype, got, opts.Width, opts.Height)
		}
	}
}

func TestQuietZoneMargin(t *testing.T) {
	opts := DefaultOptions(TypeQR)
	bc, err := Generate(TypeQR, "x", opts)
	if err != nil {
		t.Fatal(err)
	}
	// Version-1 QR: 21 modules + 2×4 quiet = 29; 256 px → 35 px margin.
	for _, p := range []image.Point{{0, 0}, {30, 30}, {255, 255}, {225, 225}} {
		if r, g, b, _ := bc.At(p.X, p.Y).RGBA(); r < 0x8000 || g < 0x8000 || b < 0x8000 {
			t.Errorf("pixel %v inside the quiet zone is dark", p)
		}
	}
	opts.NoQuietZone = true
	bare, _ := Generate(TypeQR, "x", opts)
	if withQZ, without := firstDarkColumn(bc), firstDarkColumn(bare); without >= withQZ-20 {
		t.Errorf("code starts at x=%d without and x=%d with quiet zone", without, withQZ)
	}
}

func firstDarkColumn(img image.Image) int {
	b := img.Bounds()
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			if r, _, _, _ := img.At(x, y).RGBA(); r < 0x8000 {
				return x
			}
		}
	}
	return b.Max.X
}
