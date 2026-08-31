package barcodes

import (
	"errors"
	"fmt"
	"strings"

	"github.com/boombuler/barcode"
)

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
		`<svg width="100%%" height="100%%" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg" shape-rendering="crispEdges">
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
