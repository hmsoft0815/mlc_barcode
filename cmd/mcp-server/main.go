package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mlcmcp/mlc_barcode/internal/barcodes"
	"github.com/mlcmcp/mlc_barcode/internal/qrformats"
	"github.com/mlcmcp/mlc_barcode/internal/version"
	"github.com/hmsoft0815/mlcartifact/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var artifactClient *client.Client

func main() {
	addr := flag.String("addr", "", "Listen address for SSE (e.g. \":8080\"). If empty, uses stdio.")
	artifactAddr := flag.String("artifact-addr", os.Getenv("ARTIFACT_GRPC_ADDR"), "Address of the mlcartifact gRPC server")
	showVersion := flag.Bool("version", false, "Show version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("MLC Barcode MCP Server v%s\nAuthor: %s\n", version.Version, version.Author)
		return
	}

	if *artifactAddr != "" {
		var err error
		artifactClient, err = client.NewClientWithAddr(*artifactAddr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not connect to artifact server at %s: %v\n", *artifactAddr, err)
		} else {
			fmt.Fprintf(os.Stderr, "Connected to artifact server at %s\n", *artifactAddr)
		}
	}

	ctx := context.Background()
	s := mcp.NewServer(
		&mcp.Implementation{
			Name:    "mlc-barcode-server",
			Version: version.Version,
		},
		&mcp.ServerOptions{
			Capabilities: &mcp.ServerCapabilities{
				Tools: &mcp.ToolCapabilities{ListChanged: true},
			},
		},
	)

	registerBarcodeTools(s)
	registerWifiTools(s)
	registerVCardTools(s)
	registerVCalendarTools(s)

	if *addr != "" {
		fmt.Fprintf(os.Stderr, "Starting Barcode MCP Server on SSE (%s)...\n", *addr)
		handler := mcp.NewSSEHandler(func(*http.Request) *mcp.Server { return s }, nil)
		if err := http.ListenAndServe(*addr, handler); err != nil {
			log.Fatalf("SSE server failed: %v", err)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Starting Barcode MCP Server on stdio...\n")
		transport := &mcp.StdioTransport{}
		session, err := s.Connect(ctx, transport, nil)
		if err != nil {
			log.Fatal(err)
		}
		session.Wait()
	}
}

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
