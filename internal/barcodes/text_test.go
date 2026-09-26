package barcodes

import (
	"bytes"
	"image/png"
	"strings"
	"testing"
)

func pngHeight(t *testing.T, opts BarcodeOptions) int {
	t.Helper()
	b, err := GeneratePNG(TypeCode128, "ABC-123", opts)
	if err != nil {
		t.Fatalf("GeneratePNG: %v", err)
	}
	img, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	return img.Bounds().Dy()
}

func TestPNGCaptionAddsBand(t *testing.T) {
	opts := DefaultOptions(TypeCode128)
	plain := pngHeight(t, opts)

	opts.ShowText = true
	auto := pngHeight(t, opts)
	if auto <= plain {
		t.Errorf("caption band missing: %d <= %d", auto, plain)
	}

	opts.FontSize = 60
	big := pngHeight(t, opts)
	if big <= auto {
		t.Errorf("font size 60 should need a taller band than auto: %d <= %d", big, auto)
	}
}

func TestCaptionFontSizeInSVG(t *testing.T) {
	opts := DefaultOptions(TypeCode128)
	opts.ShowText = true
	opts.FontSize = 24
	svg, err := GenerateSVG(TypeCode128, "ABC-123", opts)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, `font-size="24"`) {
		t.Errorf("font-size 24 not applied: %s", svg)
	}
}

func TestCaptionShrinksToFit(t *testing.T) {
	opts := DefaultOptions(TypeCode128)
	opts.FontSize = 100
	text := strings.Repeat("sehr langer Text ", 10)
	size, _ := captionLayout(opts, text, 600, 150)
	if size >= 100 {
		t.Errorf("long caption not shrunk: size %d", size)
	}
	if w := textWidth(text, size); w > 600 && size > MinFontSize {
		t.Errorf("caption %dpx wide still exceeds 600px at size %d", w, size)
	}
}

func TestCaptionFontSizeClamped(t *testing.T) {
	opts := BarcodeOptions{FontSize: 1}
	if size, _ := captionLayout(opts, "x", 600, 150); size != MinFontSize {
		t.Errorf("size %d, want %d", size, MinFontSize)
	}
}
