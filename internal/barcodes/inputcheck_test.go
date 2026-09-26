package barcodes

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// Error messages are read by language models; they must name the problem,
// what is allowed and what to do — and never echo the whole input.
func TestHelpfulErrors(t *testing.T) {
	long := strings.Repeat("x", 4000)
	tests := []struct {
		name  string
		btype BarcodeType
		data  string
		want  []string
	}{
		{"code39 lower case", TypeCode39, "hello", []string{"'h' at position 1", "upper case"}},
		{"code39 umlaut", TypeCode39, "GRÖSSE", []string{"'Ö' at position 3", "code128"}},
		{"code128 non-ASCII", TypeCode128, "Größe", []string{"'ö' at position 3", "ASCII only", "qr"}},
		{"itf letter", TypeITF, "12a4", []string{"'a' at position 3", "digits only"}},
		{"itf odd", TypeITF, "12345", []string{"even number", "got 5", "012345"}},
		{"datamatrix text capacity", TypeDataMatrix, long, []string{"4000 characters", "at most 1558", "qr"}},
		{"datamatrix digit capacity", TypeDataMatrix, strings.Repeat("7", 4000), []string{"at most 3116"}},
		{"qr capacity", TypeQR, strings.Repeat("x", 8000), []string{"at most 2331", "aztec"}},
		{"aztec capacity", TypeAztec, strings.Repeat("x", 5000), []string{"too much data for aztec"}},
		{"pdf417 capacity", TypePDF417, strings.Repeat("x", 3000), []string{"too much data for pdf417"}},
		{"unknown type", BarcodeType("gif"), "x", []string{"unsupported", "aztec", "pdf417"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Generate(tt.btype, tt.data, DefaultOptions(tt.btype))
			if err == nil {
				t.Fatal("expected an error")
			}
			if ie, ok := AsInputError(err); !ok || ie.Code == "" {
				t.Errorf("not an InputError with code: %v", err)
			}
			msg := err.Error()
			for _, w := range tt.want {
				if !strings.Contains(msg, w) {
					t.Errorf("message %q lacks %q", msg, w)
				}
			}
			if len(tt.data) > 100 && strings.Contains(msg, tt.data[:100]) {
				t.Errorf("message echoes the input: %.120s…", msg)
			}
		})
	}
}

func TestNewSymbologies(t *testing.T) {
	for _, tt := range []struct {
		btype BarcodeType
		data  string
	}{
		{TypeAztec, "https://mlcgo.eu"},
		{TypeAztec, "Größe"},
		{TypePDF417, "BOARDING PASS LH123 FRA-JFK"},
	} {
		svg, err := GenerateSVG(tt.btype, tt.data, DefaultOptions(tt.btype))
		if err != nil || !strings.Contains(svg, "<svg") {
			t.Errorf("%s %q: %v", tt.btype, tt.data, err)
		}
		if _, err := GeneratePNG(tt.btype, tt.data, DefaultOptions(tt.btype)); err != nil {
			t.Errorf("%s png: %v", tt.btype, err)
		}
	}
}

// A code wider than the requested size is not an error any more: it is
// rendered at its natural size instead of failing in barcode.Scale.
func TestNoScaleErrorForWideCodes(t *testing.T) {
	bc, err := Generate(TypeCode128, strings.Repeat("A", 80), DefaultOptions(TypeCode128))
	if err != nil {
		t.Fatal(err)
	}
	if bc.Bounds().Dx() < 600 {
		t.Errorf("width %d, want at least 600", bc.Bounds().Dx())
	}
}

func TestPDF417SecurityByLength(t *testing.T) {
	for _, tt := range []struct {
		n    int
		want byte
	}{{10, 2}, {100, 3}, {300, 4}, {1000, 5}} {
		if got := pdf417Security(-1, strings.Repeat("a", tt.n)); got != tt.want {
			t.Errorf("%d bytes: level %d, want %d", tt.n, got, tt.want)
		}
	}
	if got := pdf417Security(7, "a"); got != 7 {
		t.Errorf("explicit level ignored: %d", got)
	}
}

// Every error code must have a German text in the GUI catalog, or the GUI
// silently falls back to English.
func TestErrorCodesTranslated(t *testing.T) {
	src, err := os.ReadFile("../../frontend/src/lib/i18n/errors.ts")
	if err != nil {
		t.Fatal(err)
	}
	start := strings.Index(string(src), "const de: Catalog = {")
	if start < 0 {
		t.Fatal("German catalog not found")
	}
	keys := map[string]bool{}
	for _, m := range regexp.MustCompile(`(?m)^  ([a-z0-9_]+): `).FindAllStringSubmatch(string(src[start:]), -1) {
		keys[m[1]] = true
	}
	for _, c := range ErrorCodes {
		if !keys[c] {
			t.Errorf("error code %q has no German text in frontend/src/lib/i18n/errors.ts", c)
		}
	}
}
