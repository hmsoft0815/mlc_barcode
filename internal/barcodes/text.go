package barcodes

import (
	"image"
	"image/color"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// Font sizes are in output pixels (= SVG viewBox units).
const (
	MinFontSize = 6
	MaxFontSize = 200
)

var (
	captionFontOnce sync.Once
	captionFont     *opentype.Font
)

// loadCaptionFont parses the embedded Go Regular font once. It is used to
// draw the PNG caption and to measure text for both renderers, so SVG and
// PNG shrink an over-long caption the same way.
func loadCaptionFont() *opentype.Font {
	captionFontOnce.Do(func() {
		f, err := opentype.Parse(goregular.TTF)
		if err == nil {
			captionFont = f
		}
	})
	return captionFont
}

func captionFace(size int) font.Face {
	f := loadCaptionFont()
	if f == nil {
		return nil
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{Size: float64(size), DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		return nil
	}
	return face
}

func textWidth(s string, size int) int {
	face := captionFace(size)
	if face == nil {
		return len([]rune(s)) * size * 6 / 10
	}
	defer face.Close()
	return font.MeasureString(face, s).Ceil()
}

// captionText returns the text drawn under the barcode.
func captionText(opts BarcodeOptions, content string) string {
	if opts.CustomText != "" {
		return opts.CustomText
	}
	return content
}

// captionLayout returns the font size and the height of the caption band
// below a barcode of the given size. opts.FontSize 0 keeps the old
// automatic size (70 % of a band that is 20 % of the barcode height).
// The size is reduced until the text fits into the barcode width.
func captionLayout(opts BarcodeOptions, text string, width, height int) (fontSize, bandHeight int) {
	if opts.FontSize > 0 {
		fontSize = min(max(opts.FontSize, MinFontSize), MaxFontSize)
	} else {
		bandHeight = max(height/5, 20)
		fontSize = bandHeight * 7 / 10
	}

	maxWidth := width * 95 / 100
	for fontSize > MinFontSize && textWidth(text, fontSize) > maxWidth {
		fontSize--
	}

	if opts.FontSize > 0 || bandHeight == 0 {
		bandHeight = fontSize * 14 / 10
	}
	return fontSize, bandHeight
}

// captionBaseline places the text baseline inside the band so that the
// capitals keep a gap to the bars above and descenders stay inside.
func captionBaseline(top, bandHeight int) int {
	return top + bandHeight*3/4
}

// drawCaption draws text centred into the band starting at y = top.
func drawCaption(img *image.RGBA, text string, fontSize, top, bandHeight int, fg color.Color) {
	face := captionFace(fontSize)
	if face == nil {
		return
	}
	defer face.Close()

	w := font.MeasureString(face, text).Ceil()
	x := (img.Bounds().Dx() - w) / 2
	baseline := captionBaseline(top, bandHeight)
	d := &font.Drawer{Dst: img, Src: image.NewUniform(fg), Face: face, Dot: fixed.P(x, baseline)}
	d.DrawString(text)
}
