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
		"description": "Barcode type: qr, datamatrix, aztec (tickets, compact, good on screens), pdf417 (boarding passes, shipping labels), code128 (ASCII), code39 (A-Z 0-9 -.$/+%), ean13, ean8, upca, itf (even number of digits)",
		"enum":        barcodeTypeNames(),
	}
	props["data"] = map[string]any{
		"type":        "string",
		"description": "The data to encode. ean13/ean8/upca: digits only, 13/8/12 digits with check digit or 12/7/11 without (it is then computed); a wrong check digit is rejected with the expected one",
	}

	addBarcodeTool(s, &mcp.Tool{
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
		res, out, err := handleBarcodeGeneration(ctx, barcodes.BarcodeType(strings.ToLower(btypeStr)), data, args)
		return res, out, err
	})
}
