package barcodes

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Real photos from ZXing's blackbox test sets (Apache-2.0, see
// testdata/zxing/NOTICE): each image comes with a .txt holding the content.
// minFound is what the decoder must at least read; raise it when the
// decoder gets better, never lower it.
var blackboxSets = []struct {
	dir      string
	btype    BarcodeType
	minFound int
}{
	{"datamatrix-1", TypeDataMatrix, 23},
	{"datamatrix-2", TypeDataMatrix, 18},
	{"aztec-2", TypeAztec, 7},
	// Codes on bottles, made with mlc OpticScript (task testdata:cylinder,
	// tests/testdata/cylinder.js): s<degrees of circumference>-t<bottle
	// turned>-v<camera looking down>.
	// The 2D detectors sample on a perspective grid; turned bottles (t35)
	// and a camera looking down (v20, rows become arcs) are read through
	// decodeCurved (cylinder.go). Before it: DataMatrix 3, QR 5.
	{"../cylinder/datamatrix", TypeDataMatrix, 11},
	{"../cylinder/qr", TypeQR, 14},
	{"../cylinder/ean13", TypeEAN13, 8},
}

func TestDecodeBlackbox(t *testing.T) {
	for _, set := range blackboxSets {
		t.Run(set.dir, func(t *testing.T) {
			images, _ := filepath.Glob(filepath.Join("testdata", "zxing", set.dir, "*.*"))
			total, found := 0, 0
			var missed []string
			start := time.Now()
			for _, path := range images {
				ext := filepath.Ext(path)
				if ext == ".txt" || ext == ".bin" {
					continue
				}
				want, err := os.ReadFile(strings.TrimSuffix(path, ext) + ".txt")
				if err != nil {
					continue
				}
				f, err := os.Open(path)
				if err != nil {
					t.Fatal(err)
				}
				img, _, err := image.Decode(f)
				f.Close()
				if err != nil {
					t.Fatalf("%s: %v", path, err)
				}
				total++
				codes, _ := Decode(img)
				ok := false
				for _, c := range codes {
					if c.Type == set.btype && c.Text == string(want) {
						ok = true
					}
				}
				if ok {
					found++
				} else {
					missed = append(missed, filepath.Base(path))
				}
			}
			t.Logf("%s: %d/%d read in %v; missed %v", set.dir, found, total, time.Since(start).Round(time.Millisecond), missed)
			if found < set.minFound {
				t.Errorf("%s: only %d/%d read, expected at least %d", set.dir, found, total, set.minFound)
			}
		})
	}
}
