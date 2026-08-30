package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hmsoft0815/mlcartifact/client"
	"github.com/mlcmcp/mlc_barcode/internal/version"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

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
