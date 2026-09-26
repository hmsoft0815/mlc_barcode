package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"

	"github.com/mlcmcp/mlc_barcode/internal/barcodes"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	_ "golang.org/x/image/webp"
)

// Limits against oversized input: bytes before decoding, pixels from the
// image header before the pixels are allocated.
const (
	maxImageBytes  = 20 << 20
	maxImagePixels = 40_000_000
)

type decodedCode struct {
	BarcodeType string      `json:"barcode_type"`
	Text        string      `json:"text"`
	Points      []pointJSON `json:"points,omitempty"`
}

type pointJSON struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type decodeOutput struct {
	Count int           `json:"count"`
	Codes []decodedCode `json:"codes"`
}

var decodeOutputSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"count": map[string]any{"type": "integer", "description": "Number of codes found"},
		"codes": map[string]any{
			"type": "array",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"barcode_type": map[string]any{"type": "string", "description": "Symbology, e.g. qr, ean13, aztec"},
					"text":         map[string]any{"type": "string", "description": "Content of the code (EAN/UPC including check digit)"},
					"points": map[string]any{"type": "array", "description": "Corner or finder points in image pixels",
						"items": map[string]any{"type": "object", "properties": map[string]any{
							"x": map[string]any{"type": "integer"}, "y": map[string]any{"type": "integer"},
						}}},
				},
				"required": []string{"barcode_type", "text"},
			},
		},
	},
	"required": []string{"count", "codes"},
}

// registerDecodeTools adds decode_barcode. Reading files from disk is only
// offered over stdio: over HTTP any client could read the server's files.
func registerDecodeTools(s *mcp.Server, allowPath bool) {
	props := map[string]any{
		"image_base64": map[string]any{
			"type":        "string",
			"description": "PNG, JPEG, GIF or WebP image, base64 encoded (a data: URI prefix is accepted)",
		},
	}
	desc := "Reads barcodes and QR codes from an image and reports symbology and content. " +
		"Finds several codes per image. Reads qr, datamatrix, aztec, ean13, ean8, upca, code128, code39 and itf; pdf417 not yet."
	if allowPath {
		props["path"] = map[string]any{
			"type":        "string",
			"description": "Absolute path of a PNG, JPEG, GIF or WebP file on this computer",
		}
		desc += " Pass either path or image_base64."
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:         "decode_barcode",
		Title:        "Read barcodes from an image",
		Description:  desc,
		InputSchema:  map[string]any{"type": "object", "properties": props},
		OutputSchema: decodeOutputSchema,
	}, func(_ context.Context, _ *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
		path, _ := args["path"].(string)
		b64, _ := args["image_base64"].(string)
		if !allowPath {
			path = ""
		}

		img, err := loadImage(path, b64)
		if err != nil {
			return toolError(err), nil, nil
		}
		found, err := barcodes.Decode(img)
		if err != nil {
			return toolError(err), nil, nil
		}

		out := decodeOutput{Count: len(found)}
		var summary strings.Builder
		fmt.Fprintf(&summary, "Found %d code(s):", len(found))
		for _, d := range found {
			c := decodedCode{BarcodeType: string(d.Type), Text: d.Text}
			for _, p := range d.Points {
				c.Points = append(c.Points, pointJSON{X: p.X, Y: p.Y})
			}
			out.Codes = append(out.Codes, c)
			fmt.Fprintf(&summary, "\n- %s: %s", d.Type, d.Text)
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: summary.String()}}}, out, nil
	})
}

func toolError(err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}}}
}

func loadImage(path, b64 string) (image.Image, error) {
	var data []byte
	switch {
	case path != "" && b64 != "":
		return nil, errors.New("pass either path or image_base64, not both")
	case path != "":
		st, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("cannot read %s: %v — pass an absolute path to an existing image file", path, err)
		}
		if st.Size() > maxImageBytes {
			return nil, fmt.Errorf("image file is %d MB, at most %d MB are accepted", st.Size()>>20, maxImageBytes>>20)
		}
		if data, err = os.ReadFile(path); err != nil {
			return nil, fmt.Errorf("cannot read %s: %v", path, err)
		}
	case b64 != "":
		if i := strings.Index(b64, ","); strings.HasPrefix(b64, "data:") && i > 0 {
			b64 = b64[i+1:]
		}
		if base64.StdEncoding.DecodedLen(len(b64)) > maxImageBytes {
			return nil, fmt.Errorf("image is larger than %d MB", maxImageBytes>>20)
		}
		var err error
		if data, err = base64.StdEncoding.DecodeString(strings.TrimSpace(b64)); err != nil {
			return nil, fmt.Errorf("image_base64 is not valid base64: %v", err)
		}
	default:
		return nil, errors.New("no image given — pass path (absolute file path) or image_base64")
	}

	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, errors.New("unsupported image format — use PNG, JPEG, GIF or WebP (HEIC from iPhones: export as JPEG first)")
	}
	if cfg.Width*cfg.Height > maxImagePixels {
		return nil, fmt.Errorf("image is %dx%d pixels, at most %d megapixels are accepted — scale it down", cfg.Width, cfg.Height, maxImagePixels/1_000_000)
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("cannot decode %s image: %v", format, err)
	}
	return img, nil
}
