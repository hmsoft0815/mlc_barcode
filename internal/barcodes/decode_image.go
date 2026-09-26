package barcodes

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"  // decoders for DecodeImageBytes
	_ "image/jpeg" //
	_ "image/png"  //

	_ "golang.org/x/image/webp" //
)

// Limits against oversized input: bytes before decoding, pixels from the
// image header before the pixels are allocated.
const (
	MaxImageBytes  = 20 << 20
	MaxImagePixels = 40_000_000
)

// DecodeImageBytes reads a PNG, JPEG, GIF or WebP file and decodes the
// barcodes in it. It also returns the image size, so a UI can place the
// reported points over the picture.
func DecodeImageBytes(data []byte) ([]Decoded, image.Point, error) {
	if len(data) > MaxImageBytes {
		return nil, image.Point{}, inputError(ErrImageTooLarge,
			fmt.Sprintf("image is %d MB, at most %d MB are accepted", len(data)>>20, MaxImageBytes>>20),
			map[string]string{"size_mb": itoa(len(data) >> 20), "max_mb": itoa(MaxImageBytes >> 20)})
	}
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, image.Point{}, inputError(ErrImageFormat,
			"unsupported image format — use PNG, JPEG, GIF or WebP (HEIC from iPhones: export as JPEG first)", nil)
	}
	if cfg.Width*cfg.Height > MaxImagePixels {
		return nil, image.Point{}, inputError(ErrImageTooLarge,
			fmt.Sprintf("image is %dx%d pixels, at most %d megapixels are accepted — scale it down", cfg.Width, cfg.Height, MaxImagePixels/1_000_000),
			map[string]string{"width": itoa(cfg.Width), "height": itoa(cfg.Height), "max_mp": itoa(MaxImagePixels / 1_000_000)})
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, image.Point{}, inputError(ErrImageFormat, fmt.Sprintf("cannot decode %s image: %v", format, err), nil)
	}
	size := image.Pt(img.Bounds().Dx(), img.Bounds().Dy())
	found, err := Decode(img)
	return found, size, err
}
