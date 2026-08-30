package barcodes

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strings"
)

// GeneratePNG returns a PNG byte slice for the given barcode
func GeneratePNG(btype BarcodeType, data string, opts BarcodeOptions) ([]byte, error) {
	bc, err := Generate(btype, data, opts)
	if err != nil {
		return nil, err
	}

	fg := parseColor(opts.ForegroundColor, color.Black)
	bg := parseColor(opts.BackgroundColor, color.White)

	bounds := bc.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, &image.Uniform{bg}, image.Point{}, draw.Src)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := bc.At(x, y).RGBA()
			// If it's "black" in the source barcode
			if a > 0x8000 && (r < 0x8000 || g < 0x8000 || b < 0x8000) {
				img.Set(x, y, fg)
			}
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func parseColor(s string, def color.Color) color.Color {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" || s == "transparent" || s == "none" {
		return color.Transparent
	}

	switch s {
	case "white":
		return color.White
	case "black":
		return color.Black
	case "red":
		return color.RGBA{255, 0, 0, 255}
	case "green":
		return color.RGBA{0, 255, 0, 255}
	case "blue":
		return color.RGBA{0, 0, 255, 255}
	}

	// Hex parser
	if strings.HasPrefix(s, "#") {
		var r, g, b uint8
		format := "#%02x%02x%02x"
		if len(s) == 4 {
			format = "#%1x%1x%1x"
			var r1, g1, b1 uint8
			fmt.Sscanf(s, format, &r1, &g1, &b1)
			r = r1 * 17
			g = g1 * 17
			b = b1 * 17
		} else if len(s) == 7 {
			fmt.Sscanf(s, format, &r, &g, &b)
		} else {
			return def
		}
		return color.RGBA{r, g, b, 255}
	}

	return def
}
