package barcodes

import (
	"strings"
	"testing"
)

func TestGenerateSVG(t *testing.T) {
	tests := []struct {
		name    string
		btype   BarcodeType
		data    string
		wantErr bool
	}{
		{"QR Code", TypeQR, "Hello QR", false},
		{"DataMatrix", TypeDataMatrix, "Hello DM", false},
		{"Code128", TypeCode128, "12345", false},
		{"Invalid Type", BarcodeType("invalid"), "data", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := DefaultOptions(tt.btype)
			svg, err := GenerateSVG(tt.btype, tt.data, opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateSVG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !strings.Contains(svg, "<svg") {
					t.Errorf("GenerateSVG() produced invalid SVG: %s", svg)
				}
				if !strings.Contains(svg, "</svg>") {
					t.Errorf("GenerateSVG() produced invalid SVG: %s", svg)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := DefaultOptions(tt.btype)
			png, err := GeneratePNG(tt.btype, tt.data, opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePNG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(png) < 10 {
				t.Errorf("GeneratePNG() produced too small output: %d bytes", len(png))
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
}
