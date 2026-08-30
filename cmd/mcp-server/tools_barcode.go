package main

import (
	"context"
	"strings"

	"github.com/mlcmcp/mlc_barcode/internal/barcodes"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerBarcodeTools(s *mcp.Server) {
	props := getCommonProperties()
	props["type"] = map[string]any{
		"type":        "string",
		"description": "Barcode type (qr, datamatrix, code128, code39, ean13, ean8, upca, itf)",
		"enum":        []string{"qr", "datamatrix", "code128", "code39", "ean13", "ean8", "upca", "itf"},
	}
	props["data"] = map[string]any{
		"type":        "string",
		"description": "The data to encode in the barcode",
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "generate_barcode",
		Description: "Generates a barcode image (SVG or PNG) from data",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": props,
			"required":   []string{"type", "data"},
		},
	}, func(ctx context.Context, request *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
		btypeStr, _ := args["type"].(string)
		data, _ := args["data"].(string)
		res, err := handleBarcodeGeneration(ctx, barcodes.BarcodeType(strings.ToLower(btypeStr)), data, args)
		return res, nil, err
	})
}
