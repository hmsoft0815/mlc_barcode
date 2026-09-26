package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/mlcmcp/mlc_barcode/internal/barcodes"
	_ "golang.org/x/image/webp"
)

// runDecode prints one line per code found: type, a tab, the content.
// Exit code 1 when the file cannot be read or holds no code.
func runDecode(path string) int {
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s is not a PNG, JPEG, GIF or WebP image: %v\n", path, err)
		return 1
	}
	found, err := barcodes.Decode(img)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}
	for _, d := range found {
		fmt.Printf("%s\t%s\n", d.Type, d.Text)
	}
	return 0
}
