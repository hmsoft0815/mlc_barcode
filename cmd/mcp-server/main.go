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

func registerBarcodeTools(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "generate_barcode",
		Description: "Generates a barcode image (SVG or PNG) from data",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"type": map[string]any{
					"type":        "string",
					"description": "Barcode type (qr, datamatrix, code128, code39, ean13, ean8, upca, itf)",
					"enum":        []string{"qr", "datamatrix", "code128", "code39", "ean13", "ean8", "upca", "itf"},
				},
				"data": map[string]any{
					"type":        "string",
					"description": "The data to encode in the barcode",
				},
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
				"save_artifact": map[string]any{
					"type":        "boolean",
					"description": "If true, saves the barcode to mlcartifact service (requires server connection)",
				},
				"filename": map[string]any{
					"type":        "string",
					"description": "Optional filename for the artifact (e.g. 'mybarcode.png')",
				},
			},
			"required": []string{"type", "data"},
		},
	}, func(ctx context.Context, request *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
		btypeStr, _ := args["type"].(string)
		data, _ := args["data"].(string)
		format, _ := args["format"].(string)
		if format == "" {
			format = "svg"
		}

		btype := barcodes.BarcodeType(strings.ToLower(btypeStr))
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
			return nil, nil, fmt.Errorf("unsupported format: %s", format)
		}

		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error generating barcode: %v", err)}},
			}, nil, nil
		}

		results := []mcp.Content{}

		// Optional artifact saving
		saveArtifact, _ := args["save_artifact"].(bool)
		if saveArtifact {
			if artifactClient == nil {
				results = append(results, &mcp.TextContent{Text: "Warning: Artifact saving requested but no artifact server connection available."})
			} else {
				fname, _ := args["filename"].(string)
				if fname == "" {
					fname = fmt.Sprintf("barcode_%s.%s", btype, format)
				}
				resp, err := artifactClient.Write(ctx, fname, content, client.WithMimeType(mimeType), client.WithDescription(fmt.Sprintf("Generated %s barcode for: %s", btype, data)))
				if err != nil {
					results = append(results, &mcp.TextContent{Text: fmt.Sprintf("Error saving artifact: %v", err)})
				} else {
					results = append(results, &mcp.TextContent{Text: fmt.Sprintf("Artifact saved as '%s' (ID: %s)", fname, resp.Id)})
				}
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

		return &mcp.CallToolResult{
			Content: results,
		}, nil, nil
	})
}
