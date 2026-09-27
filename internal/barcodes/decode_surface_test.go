package barcodes

import (
	"image"
	_ "image/jpeg"
	"os"
	"path/filepath"
	"testing"
)

// With the scan view's shape switch on "bottle" (the code about half as
// wide as the bottle, 60°), a camera frame gets every turn and tilt of
// that curvature and reads what the full search of a still photo reads.
func TestDecodeSurfaceHint(t *testing.T) {
	for _, set := range []struct {
		dir   string
		btype BarcodeType
	}{{"datamatrix", TypeDataMatrix}, {"qr", TypeQR}} {
		files, _ := filepath.Glob(filepath.Join("testdata", "cylinder", set.dir, "s60-*.jpg"))
		for _, f := range files {
			want, err := os.ReadFile(f[:len(f)-len(".jpg")] + ".txt")
			if err != nil {
				t.Fatal(err)
			}
			fh, err := os.Open(f)
			if err != nil {
				t.Fatal(err)
			}
			img, _, err := image.Decode(fh)
			fh.Close()
			if err != nil {
				t.Fatal(err)
			}
			for _, opts := range []DecodeOptions{{Live: true, Surface: SurfaceBottle}, {Live: true}} {
				found, _ := DecodeWith(img, opts)
				ok := false
				for _, d := range found {
					ok = ok || (d.Type == set.btype && d.Text == string(want))
				}
				t.Logf("%s %s surface=%d: %v", set.dir, filepath.Base(f), opts.Surface, ok)
				if opts.Surface == SurfaceBottle && !ok {
					t.Errorf("%s: not read with the bottle hint", f)
				}
			}
		}
	}
	// A flat pack skips the bottle search: nothing on a curved code.
	fh, _ := os.Open(filepath.Join("testdata", "cylinder", "datamatrix", "s60-t35-v20.jpg"))
	img, _, _ := image.Decode(fh)
	fh.Close()
	if found, _ := DecodeWith(img, DecodeOptions{Live: true, Surface: SurfaceFlat}); len(found) != 0 {
		t.Errorf("flat hint still unwrapped: %v", found)
	}
}
