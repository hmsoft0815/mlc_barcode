package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"sort"
	"strings"

	"github.com/mlcmcp/mlc_barcode/internal/barcodes"
	"github.com/mlcmcp/mlc_barcode/internal/qrformats"
	_ "golang.org/x/image/webp"
)

// runDecode prints one line per code found: type, a tab, the content —
// followed by an indented line with the fields of a known payload.
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
		// One line per code: multi-line payloads (vCard, GiroCode) get \n.
		fmt.Printf("%s\t%s\n", d.Type, strings.NewReplacer("\r", `\r`, "\n", `\n`).Replace(d.Text))
		if p := qrformats.Parse(d.Text); p.Kind != "text" {
			keys := make([]string, 0, len(p.Fields))
			for k := range p.Fields {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			fmt.Printf("\t%s:", p.Kind)
			for _, k := range keys {
				fmt.Printf(" %s=%q", k, p.Fields[k])
			}
			fmt.Println()
		}
	}
	return 0
}
