package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hmsoft0815/mlcartifact/client"
	"github.com/mlcmcp/mlc_barcode/internal/version"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	addr := flag.String("addr", "", "Listen address for HTTP (e.g. \":8080\"): Streamable HTTP at /mcp, legacy SSE at /sse. If empty, uses stdio.")
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
			Instructions: serverInstructions,
			Capabilities: &mcp.ServerCapabilities{
				Tools:   &mcp.ToolCapabilities{ListChanged: true},
				Prompts: &mcp.PromptCapabilities{ListChanged: true},
			},
		},
	)

	registerBarcodeTools(s)
	registerWifiTools(s)
	registerVCardTools(s)
	registerVCalendarTools(s)
	registerEPCTools(s)
	registerCryptoTools(s)
	registerGeoTools(s)
	registerCommunicationTools(s)
	registerPrompts(s)

	if *addr != "" {
		fmt.Fprintf(os.Stderr, "Starting Barcode MCP Server on %s: Streamable HTTP at /mcp, legacy SSE at /sse\n", *addr)
		if err := newHTTPServer(*addr, s).ListenAndServe(); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
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
