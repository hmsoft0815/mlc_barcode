package gui

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mlcmcp/mlc_barcode/internal/photokit"
)

// LibraryPhoto is one photo taken from the system photo library, ready for
// the same path as a picked file: a data: URL the decoder and <img> read.
// An empty DataURL means the user closed the picker.
type LibraryPhoto struct {
	Name    string `json:"name"`
	DataURL string `json:"dataUrl"`
}

// PhotoLibraryAvailable reports whether "From Photos" works here (macOS
// 13+). Elsewhere the frontend hides the button and keeps the file picker.
func (a *BarcodeApp) PhotoLibraryAvailable() bool { return photokit.Supported() }

// PickFromPhotos shows the system photo picker for one photo. It runs on
// the service goroutine, as photokit requires; the picker needs no
// permission and hands over the original file.
func (a *BarcodeApp) PickFromPhotos() (LibraryPhoto, error) {
	dir, err := os.MkdirTemp("", "mlc-barcode-photos-")
	if err != nil {
		return LibraryPhoto{}, err
	}
	defer os.RemoveAll(dir)

	photos, err := photokit.Pick(photokit.PickOptions{Limit: 1, Dir: dir})
	if errors.Is(err, photokit.ErrCancelled) || (err == nil && len(photos) == 0) {
		return LibraryPhoto{}, nil
	}
	if len(photos) == 0 {
		return LibraryPhoto{}, err
	}
	p := photos[0]

	path, mime := p.Path, ""
	switch p.Type {
	case "public.jpeg":
		mime = "image/jpeg"
	case "public.png":
		mime = "image/png"
	default:
		// HEIC, RAW, TIFF …: neither the decoder nor every web view reads
		// them. sips ships with macOS and writes a JPEG the decoder takes.
		path = filepath.Join(dir, "converted.jpg")
		if out, cerr := exec.Command("sips", "-s", "format", "jpeg", p.Path, "--out", path).CombinedOutput(); cerr != nil {
			return LibraryPhoto{}, fmt.Errorf("convert %s (%s): %v: %s", p.Name, p.Type, cerr, out)
		}
		mime = "image/jpeg"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return LibraryPhoto{}, err
	}
	return LibraryPhoto{
		Name:    p.Name,
		DataURL: "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(data),
	}, nil
}
