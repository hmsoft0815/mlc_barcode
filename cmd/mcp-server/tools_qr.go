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

func registerEPCTools(s *mcp.Server) {
	props := getCommonProperties()
	props["name"] = map[string]any{"type": "string", "description": "Beneficiary name (Empfänger)"}
	props["iban"] = map[string]any{"type": "string", "description": "IBAN (e.g. DE89370400440532013000)"}
	props["bic"] = map[string]any{"type": "string", "description": "BIC (optional)"}
	props["amount"] = map[string]any{"type": "number", "description": "Amount in EUR (e.g. 12.50)"}
	props["reference"] = map[string]any{"type": "string", "description": "Remittance information / Verwendungszweck"}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "generate_epc_qr",
		Description: "Generates an EPC-QR-Code (GiroCode / SEPA-Überweisung) for banking apps and invoice payments",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": props,
			"required":   []string{"name", "iban"},
		},
	}, func(ctx context.Context, request *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
		name, _ := args["name"].(string)
		iban, _ := args["iban"].(string)
		bic, _ := args["bic"].(string)
		ref, _ := args["reference"].(string)
		var amount float64
		if v, ok := args["amount"].(float64); ok {
			amount = v
		}

		data := qrformats.FormatEPC(qrformats.EPCOptions{
			Name:      name,
			IBAN:      iban,
			BIC:       bic,
			Amount:    amount,
			Reference: ref,
		})
		res, err := handleBarcodeGeneration(ctx, barcodes.TypeQR, data, args)
		return res, nil, err
	})
}

func registerCryptoTools(s *mcp.Server) {
	props := getCommonProperties()
	props["coin"] = map[string]any{"type": "string", "description": "Cryptocurrency (bitcoin, ethereum, solana, etc.)", "default": "bitcoin"}
	props["address"] = map[string]any{"type": "string", "description": "Crypto wallet address"}
	props["amount"] = map[string]any{"type": "number", "description": "Optional payment amount"}
	props["label"] = map[string]any{"type": "string", "description": "Optional recipient label"}
	props["message"] = map[string]any{"type": "string", "description": "Optional payment message/memo"}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "generate_crypto_qr",
		Description: "Generates a cryptocurrency payment QR code (Bitcoin BIP 21, Ethereum EIP-681, etc.)",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": props,
			"required":   []string{"address"},
		},
	}, func(ctx context.Context, request *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
		coin, _ := args["coin"].(string)
		addr, _ := args["address"].(string)
		label, _ := args["label"].(string)
		msg, _ := args["message"].(string)
		var amount float64
		if v, ok := args["amount"].(float64); ok {
			amount = v
		}

		data := qrformats.FormatCrypto(qrformats.CryptoOptions{
			Coin:    coin,
			Address: addr,
			Amount:  amount,
			Label:   label,
			Message: msg,
		})
		res, err := handleBarcodeGeneration(ctx, barcodes.TypeQR, data, args)
		return res, nil, err
	})
}

func registerGeoTools(s *mcp.Server) {
	props := getCommonProperties()
	props["latitude"] = map[string]any{"type": "number", "description": "Latitude coordinate (e.g. 52.5200)"}
	props["longitude"] = map[string]any{"type": "number", "description": "Longitude coordinate (e.g. 13.4050)"}
	props["query"] = map[string]any{"type": "string", "description": "Location name or search query"}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "generate_geo_qr",
		Description: "Generates a Geo Location (RFC 5870) QR code that opens Google Maps or Apple Maps",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": props,
			"required":   []string{"latitude", "longitude"},
		},
	}, func(ctx context.Context, request *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
		var lat, lon float64
		if v, ok := args["latitude"].(float64); ok {
			lat = v
		}
		if v, ok := args["longitude"].(float64); ok {
			lon = v
		}
		query, _ := args["query"].(string)

		data := qrformats.FormatGeo(qrformats.GeoOptions{
			Latitude:  lat,
			Longitude: lon,
			Query:     query,
		})
		res, err := handleBarcodeGeneration(ctx, barcodes.TypeQR, data, args)
		return res, nil, err
	})
}

func registerCommunicationTools(s *mcp.Server) {
	// Telephone Tool
	telProps := getCommonProperties()
	telProps["phone_number"] = map[string]any{"type": "string", "description": "Phone number to dial (e.g. +49123456789)"}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "generate_tel_qr",
		Description: "Generates a telephone QR code (tel:) that opens the phone dialer",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": telProps,
			"required":   []string{"phone_number"},
		},
	}, func(ctx context.Context, request *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
		num, _ := args["phone_number"].(string)
		data := qrformats.FormatTel(qrformats.TelOptions{PhoneNumber: num})
		res, err := handleBarcodeGeneration(ctx, barcodes.TypeQR, data, args)
		return res, nil, err
	})

	// SMS Tool
	smsProps := getCommonProperties()
	smsProps["phone_number"] = map[string]any{"type": "string", "description": "Recipient phone number"}
	smsProps["message"] = map[string]any{"type": "string", "description": "SMS message text"}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "generate_sms_qr",
		Description: "Generates an SMS QR code (smsto:) with pre-filled recipient and message",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": smsProps,
			"required":   []string{"phone_number"},
		},
	}, func(ctx context.Context, request *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
		num, _ := args["phone_number"].(string)
		msg, _ := args["message"].(string)
		data := qrformats.FormatSMS(qrformats.SMSOptions{PhoneNumber: num, Message: msg})
		res, err := handleBarcodeGeneration(ctx, barcodes.TypeQR, data, args)
		return res, nil, err
	})

	// Email Tool
	mailProps := getCommonProperties()
	mailProps["to"] = map[string]any{"type": "string", "description": "Recipient email address"}
	mailProps["subject"] = map[string]any{"type": "string", "description": "Email subject"}
	mailProps["body"] = map[string]any{"type": "string", "description": "Email body content"}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "generate_email_qr",
		Description: "Generates an Email QR code (mailto:) with recipient, subject, and body",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": mailProps,
			"required":   []string{"to"},
		},
	}, func(ctx context.Context, request *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
		to, _ := args["to"].(string)
		subject, _ := args["subject"].(string)
		body, _ := args["body"].(string)
		data := qrformats.FormatEmail(qrformats.EmailOptions{To: to, Subject: subject, Body: body})
		res, err := handleBarcodeGeneration(ctx, barcodes.TypeQR, data, args)
		return res, nil, err
	})
}
