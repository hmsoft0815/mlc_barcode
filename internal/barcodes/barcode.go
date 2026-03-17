/*
 * PDF Generation Service
 * Copyright (c) 2026 Michael Lechner
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */
package barcodes

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/code39"
	"github.com/boombuler/barcode/datamatrix"
	"github.com/boombuler/barcode/ean"
	"github.com/boombuler/barcode/qr"
	"github.com/boombuler/barcode/twooffive"
)

type BarcodeType string

const (
	TypeQR         BarcodeType = "qr"
	TypeDataMatrix BarcodeType = "datamatrix"
	TypeCode128    BarcodeType = "code128"
	TypeCode39     BarcodeType = "code39"
	TypeEAN13      BarcodeType = "ean13"
	TypeEAN8       BarcodeType = "ean8"
	TypeUPCA       BarcodeType = "upca"
	TypeITF        BarcodeType = "itf"
)

// BarcodeOptions holds configuration for barcode generation
type BarcodeOptions struct {
	Width           int
	Height          int
	ShowText        bool
	ForegroundColor string // SVG color (e.g. "black", "#000")
	BackgroundColor string // SVG color (e.g. "white", "#fff")
	// For QR codes
	QRLevel qr.ErrorCorrectionLevel
	QRMode  qr.Encoding
	// For Code39
	IncludeChecksum bool
	FullASCIICode39 bool
}

// DefaultOptions returns recommended default options for a barcode type
func DefaultOptions(btype BarcodeType) BarcodeOptions {
	opts := BarcodeOptions{
		QRLevel:         qr.M,
		QRMode:          qr.Auto,
		ShowText:        false,
		ForegroundColor: "black",
		BackgroundColor: "white",
	}

	switch btype {
	case TypeQR, TypeDataMatrix:
		opts.Width = 256
		opts.Height = 256
	case TypeEAN13, TypeUPCA:
		opts.Width = 450
		opts.Height = 150
	case TypeEAN8:
		opts.Width = 320
		opts.Height = 150
	default:
		opts.Width = 600
		opts.Height = 150
	}

	return opts
}

// Generate generates a barcode object
func Generate(btype BarcodeType, data string, opts BarcodeOptions) (barcode.Barcode, error) {
	var bc barcode.Barcode
	var err error

	data = strings.TrimSpace(data)

	switch btype {
	case TypeQR:
		bc, err = qr.Encode(data, opts.QRLevel, opts.QRMode)
	case TypeDataMatrix:
		bc, err = datamatrix.Encode(data)
	case TypeCode128:
		bc, err = code128.Encode(data)
	case TypeCode39:
		bc, err = code39.Encode(data, opts.IncludeChecksum, opts.FullASCIICode39)
	case TypeEAN13, TypeEAN8, TypeUPCA:
		bc, err = ean.Encode(data)
	case TypeITF:
		bc, err = twooffive.Encode(data, true)
	default:
		return nil, fmt.Errorf("unsupported barcode type: %s", btype)
	}

	if err != nil {
		return nil, err
	}

	// Scale to requested size
	if opts.Width > 0 && opts.Height > 0 {
		bc, err = barcode.Scale(bc, opts.Width, opts.Height)
		if err != nil {
			return nil, err
		}
	}

	return bc, nil
}

// GeneratePNG returns a PNG byte slice for the given barcode
func GeneratePNG(btype BarcodeType, data string, opts BarcodeOptions) ([]byte, error) {
	bc, err := Generate(btype, data, opts)
	if err != nil {
		return nil, err
	}

	// barcode library doesn't easily support custom colors in Encode,
	// so we might need to manually remap colors if they aren't black/white.
	// For now, we'll implement a simple color remapper if colors are provided.

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

// GenerateSVG returns an SVG string for the given barcode
func GenerateSVG(btype BarcodeType, data string, opts BarcodeOptions) (string, error) {
	bc, err := Generate(btype, data, opts)
	if err != nil {
		return "", err
	}

	return barcodeToSVG(bc, opts)
}

// barcodeToSVG converts a barcode.Barcode to an SVG string
func barcodeToSVG(bc barcode.Barcode, opts BarcodeOptions) (string, error) {
	if bc == nil {
		return "", errors.New("barcode is nil")
	}

	bounds := bc.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	// If we show text, we need extra height in the viewBox
	viewBoxHeight := height
	textHeight := 0
	if opts.ShowText {
		textHeight = height / 5 // Reserve 20% of height for text
		if textHeight < 20 {
			textHeight = 20
		}
		viewBoxHeight += textHeight
	}

	var pathData strings.Builder

	// Optimize by grouping consecutive black pixels in each row
	for y := 0; y < height; y++ {
		inBar := false
		startX := 0
		for x := 0; x < width; x++ {
			r, g, b, a := bc.At(x, y).RGBA()
			isBlack := (a > 0x8000) && (r < 0x8000 || g < 0x8000 || b < 0x8000)

			if isBlack {
				if !inBar {
					startX = x
					inBar = true
				}
			} else if inBar {
				w := x - startX
				fmt.Fprintf(&pathData, "M%d %d h%d v1 h-%d z ", startX, y, w, w)
				inBar = false
			}
		}
		if inBar {
			w := width - startX
			fmt.Fprintf(&pathData, "M%d %d h%d v1 h-%d z ", startX, y, w, w)
		}
	}

	textElement := ""
	if opts.ShowText {
		content := bc.Content()
		fontSize := textHeight * 8 / 10
		textY := height + (textHeight * 7 / 10)
		textElement = fmt.Sprintf(
			`<text x="%d" y="%d" font-family="monospace" font-size="%d" text-anchor="middle" fill="%s">%s</text>`,
			width/2, textY, fontSize, opts.ForegroundColor, content,
		)
	}

	bgStyle := opts.BackgroundColor
	if strings.ToLower(bgStyle) == "transparent" {
		bgStyle = "none"
	}

	svg := fmt.Sprintf(
		`<svg viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg" shape-rendering="crispEdges">
  <rect width="%d" height="%d" fill="%s"/>
  <path d="%s" fill="%s"/>
  %s
</svg>`,
		width, viewBoxHeight,
		width, viewBoxHeight, bgStyle,
		pathData.String(), opts.ForegroundColor,
		textElement,
	)

	return svg, nil
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

// Helper functions for backward compatibility or ease of use

func CreateBarcodeToSVG(barcodetype string, data string, width, height int) (string, error) {
	btype := BarcodeType(strings.ToLower(barcodetype))
	opts := DefaultOptions(btype)
	if width > 0 {
		opts.Width = width
	}
	if height > 0 {
		opts.Height = height
	}
	return GenerateSVG(btype, data, opts)
}
