package main

import (
	"context"
	"fmt"

	"github.com/hmsoft0815/mlcartifact/client"
	"github.com/mlcmcp/mlc_barcode/internal/barcodes"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var artifactClient *client.Client

func getCommonProperties() map[string]any {
	props := map[string]any{
		"format": map[string]any{
			"type":        "string",
			"description": "Output format (svg or png)",
			"enum":        []string{"svg", "png"},
			"default":     "svg",
		},
		"width": map[string]any{
			"type":        "integer",
			"description": "Width of the image (optional)",
		},
		"height": map[string]any{
			"type":        "integer",
			"description": "Height of the image (optional)",
		},
		"text": map[string]any{
			"type":        "boolean",
			"description": "Show text below barcode (if supported)",
		},
		"fg_color": map[string]any{
			"type":        "string",
			"description": "Foreground color (e.g. black, #ff0000)",
			"default":     "black",
		},
		"bg_color": map[string]any{
			"type":        "string",
			"description": "Background color (e.g. white, transparent, #ffffff)",
			"default":     "white",
		},
	}

	if artifactClient != nil {
		props["save_artifact"] = map[string]any{
			"type":        "boolean",
			"description": "If true, saves the barcode to mlcartifact service",
		}
		props["filename"] = map[string]any{
			"type":        "string",
			"description": "Optional filename for the artifact (e.g. 'mybarcode.png')",
		}
	}
	return props
}

func handleBarcodeGeneration(ctx context.Context, btype barcodes.BarcodeType, data string, args map[string]any) (*mcp.CallToolResult, error) {
	format, _ := args["format"].(string)
	if format == "" {
		format = "svg"
	}

	opts := barcodes.DefaultOptions(btype)
	if w, ok := args["width"].(float64); ok {
		opts.Width = int(w)
	}
	if h, ok := args["height"].(float64); ok {
		opts.Height = int(h)
	}
	if t, ok := args["text"].(bool); ok {
		opts.ShowText = t
	}
	if fg, ok := args["fg_color"].(string); ok && fg != "" {
		opts.ForegroundColor = fg
	}
	if bg, ok := args["bg_color"].(string); ok && bg != "" {
		opts.BackgroundColor = bg
	}

	var content []byte
	var err error
	var mimeType string

	switch format {
	case "svg":
		var svg string
		svg, err = barcodes.GenerateSVG(btype, data, opts)
		content = []byte(svg)
		mimeType = "image/svg+xml"
	case "png":
		content, err = barcodes.GeneratePNG(btype, data, opts)
		mimeType = "image/png"
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error generating barcode: %v", err)}},
		}, nil
	}

	results := []mcp.Content{}

	// Optional artifact saving
	saveArtifact, _ := args["save_artifact"].(bool)
	if saveArtifact && artifactClient != nil {
		fname, _ := args["filename"].(string)
		if fname == "" {
			fname = fmt.Sprintf("barcode_%s.%s", btype, format)
		}
		resp, err := artifactClient.Write(ctx, fname, content, client.WithMimeType(mimeType), client.WithDescription(fmt.Sprintf("Generated %s barcode", btype)))
		if err != nil {
			results = append(results, &mcp.TextContent{Text: fmt.Sprintf("Error saving artifact: %v", err)})
		} else {
			results = append(results, &mcp.TextContent{Text: fmt.Sprintf("Artifact saved as '%s' (ID: %s)", fname, resp.Id)})
		}
	}

	if format == "svg" {
		results = append(results, &mcp.TextContent{Text: string(content)})
	} else {
		results = append(results, &mcp.TextContent{Text: fmt.Sprintf("Generated %s barcode in %s format.", btype, format)})
		results = append(results, &mcp.ImageContent{
			Data:     content,
			MIMEType: mimeType,
		})
	}

	return &mcp.CallToolResult{Content: results}, nil
}
