package main

import (
	"context"

	"github.com/mlcmcp/mlc_barcode/internal/barcodes"
	"github.com/mlcmcp/mlc_barcode/internal/qrformats"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerWifiTools(s *mcp.Server) {
	props := getCommonProperties()
	props["ssid"] = map[string]any{"type": "string", "description": "WIFI Network Name (SSID)"}
	props["password"] = map[string]any{"type": "string", "description": "WIFI Password"}
	props["encryption"] = map[string]any{
		"type":        "string",
		"description": "Encryption type",
		"enum":        []string{"WPA", "WEP", "nopass"},
		"default":     "WPA",
	}
	props["hidden"] = map[string]any{"type": "boolean", "description": "Hidden network"}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "generate_wifi_qr",
		Description: "Generates a QR code for WIFI access",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": props,
			"required":   []string{"ssid"},
		},
	}, func(ctx context.Context, request *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
		ssid, _ := args["ssid"].(string)
		pass, _ := args["password"].(string)
		enc, _ := args["encryption"].(string)
		hidden, _ := args["hidden"].(bool)

		data := qrformats.FormatWifi(qrformats.WifiOptions{
			SSID:       ssid,
			Password:   pass,
			Encryption: enc,
			Hidden:     hidden,
		})
		res, err := handleBarcodeGeneration(ctx, barcodes.TypeQR, data, args)
		return res, nil, err
	})
}

func registerVCardTools(s *mcp.Server) {
	props := getCommonProperties()
	props["first_name"] = map[string]any{"type": "string"}
	props["last_name"] = map[string]any{"type": "string"}
	props["org"] = map[string]any{"type": "string", "description": "Organization"}
	props["title"] = map[string]any{"type": "string"}
	props["phone"] = map[string]any{"type": "string"}
	props["email"] = map[string]any{"type": "string"}
	props["address"] = map[string]any{"type": "string"}
	props["city"] = map[string]any{"type": "string"}
	props["zip"] = map[string]any{"type": "string"}
	props["country"] = map[string]any{"type": "string"}
	props["url"] = map[string]any{"type": "string"}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "generate_vcard_qr",
		Description: "Generates a QR code for a vCard 3.0 contact",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": props,
			"required":   []string{"first_name", "last_name"},
		},
	}, func(ctx context.Context, request *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
		opts := qrformats.VCardOptions{}
		opts.FirstName, _ = args["first_name"].(string)
		opts.LastName, _ = args["last_name"].(string)
		opts.Organization, _ = args["org"].(string)
		opts.Title, _ = args["title"].(string)
		opts.Phone, _ = args["phone"].(string)
		opts.Email, _ = args["email"].(string)
		opts.Address, _ = args["address"].(string)
		opts.City, _ = args["city"].(string)
		opts.Zip, _ = args["zip"].(string)
		opts.Country, _ = args["country"].(string)
		opts.URL, _ = args["url"].(string)

		data := qrformats.FormatVCard(opts)
		res, err := handleBarcodeGeneration(ctx, barcodes.TypeQR, data, args)
		return res, nil, err
	})
}

func registerVCalendarTools(s *mcp.Server) {
	props := getCommonProperties()
	props["summary"] = map[string]any{"type": "string", "description": "Event title"}
	props["description"] = map[string]any{"type": "string"}
	props["location"] = map[string]any{"type": "string"}
	props["start_time"] = map[string]any{"type": "string", "description": "YYYYMMDDTHHMMSS(Z)"}
	props["end_time"] = map[string]any{"type": "string", "description": "YYYYMMDDTHHMMSS(Z)"}
	props["timezone"] = map[string]any{"type": "string", "description": "e.g. Europe/Berlin"}
	props["latitude"] = map[string]any{"type": "number"}
	props["longitude"] = map[string]any{"type": "number"}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "generate_event_qr",
		Description: "Generates a QR code for an iCalendar (RFC 5545) event",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": props,
			"required":   []string{"summary", "start_time"},
		},
	}, func(ctx context.Context, request *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
		opts := qrformats.VCalendarOptions{}
		opts.Summary, _ = args["summary"].(string)
		opts.Description, _ = args["description"].(string)
		opts.Location, _ = args["location"].(string)
		opts.StartTime, _ = args["start_time"].(string)
		opts.EndTime, _ = args["end_time"].(string)
		opts.TimeZone, _ = args["timezone"].(string)
		if v, ok := args["latitude"].(float64); ok {
			opts.Latitude = v
		}
		if v, ok := args["longitude"].(float64); ok {
			opts.Longitude = v
		}

		data := qrformats.FormatVCalendar(opts)
		res, err := handleBarcodeGeneration(ctx, barcodes.TypeQR, data, args)
		return res, nil, err
	})
}
