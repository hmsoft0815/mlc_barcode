package barcodes

import (
	"image/color"
	"strings"
	"testing"
)

func TestGenerateSVG(t *testing.T) {
	tests := []struct {
		name     string
		btype    BarcodeType
		data     string
		showText bool
		wantErr  bool
	}{
		{"QR Code", TypeQR, "Hello QR", false, false},
		{"DataMatrix", TypeDataMatrix, "Hello DM", false, false},
		{"Code128", TypeCode128, "12345", true, false},
		{"Code39", TypeCode39, "HELLO39", true, false},
		{"EAN13", TypeEAN13, "4006381333931", false, false},
		{"EAN8", TypeEAN8, "96385074", false, false},
		{"UPCA", TypeUPCA, "012345678905", false, false},
		{"ITF", TypeITF, "123456", false, false},
		{"Invalid Type", BarcodeType("invalid"), "data", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := DefaultOptions(tt.btype)
			opts.ShowText = tt.showText
			opts.BackgroundColor = "transparent"
			opts.ForegroundColor = "#ff0000"
			svg, err := GenerateSVG(tt.btype, tt.data, opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateSVG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !strings.Contains(svg, "<svg") || !strings.Contains(svg, "</svg>") {
					t.Errorf("GenerateSVG() produced invalid SVG: %s", svg)
				}
				if tt.showText && !strings.Contains(svg, "<text") {
					t.Errorf("GenerateSVG() with showText did not include <text tag")
				}
			}
		})
	}
}

func TestGeneratePNG(t *testing.T) {
	tests := []struct {
		name    string
		btype   BarcodeType
		data    string
		wantErr bool
	}{
		{"QR Code PNG", TypeQR, "Hello PNG", false},
		{"Code39 PNG", TypeCode39, "ABC-123", false},
		{"DataMatrix PNG", TypeDataMatrix, "DM PNG", false},
		{"EAN13 PNG", TypeEAN13, "4006381333931", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := DefaultOptions(tt.btype)
			opts.ForegroundColor = "#00f"
			opts.BackgroundColor = "#fff"
			pngData, err := GeneratePNG(tt.btype, tt.data, opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePNG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(pngData) < 10 {
				t.Errorf("GeneratePNG() produced too small output: %d bytes", len(pngData))
			}
		})
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions(TypeQR)
	if opts.Width != 256 || opts.Height != 256 {
		t.Errorf("DefaultOptions(TypeQR) = %dx%d, want 256x256", opts.Width, opts.Height)
	}

	optsEAN := DefaultOptions(TypeEAN13)
	if optsEAN.Width != 450 || optsEAN.Height != 150 {
		t.Errorf("DefaultOptions(TypeEAN13) = %dx%d, want 450x150", optsEAN.Width, optsEAN.Height)
	}

	optsEAN8 := DefaultOptions(TypeEAN8)
	if optsEAN8.Width != 320 || optsEAN8.Height != 150 {
		t.Errorf("DefaultOptions(TypeEAN8) = %dx%d, want 320x150", optsEAN8.Width, optsEAN8.Height)
	}

	optsDefault := DefaultOptions(TypeCode128)
	if optsDefault.Width != 600 || optsDefault.Height != 150 {
		t.Errorf("DefaultOptions(TypeCode128) = %dx%d, want 600x150", optsDefault.Width, optsDefault.Height)
	}
}

func TestParseColor(t *testing.T) {
	tests := []struct {
		input string
		want  color.Color
	}{
		{"", color.Transparent},
		{"transparent", color.Transparent},
		{"none", color.Transparent},
		{"white", color.White},
		{"black", color.Black},
		{"red", color.RGBA{255, 0, 0, 255}},
		{"green", color.RGBA{0, 255, 0, 255}},
		{"blue", color.RGBA{0, 0, 255, 255}},
		{"#f00", color.RGBA{255, 0, 0, 255}},
		{"#00ff00", color.RGBA{0, 255, 0, 255}},
		{"#invalid", color.Black},
	}

	for _, tt := range tests {
		got := parseColor(tt.input, color.Black)
		if got != tt.want {
			t.Errorf("parseColor(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestCreateBarcodeToSVG(t *testing.T) {
	svg, err := CreateBarcodeToSVG("qr", "test-svg", 200, 200)
	if err != nil {
		t.Fatalf("CreateBarcodeToSVG failed: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Errorf("expected valid svg, got: %s", svg)
	}
}
