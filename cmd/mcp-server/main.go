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

// listTTLMs is how long clients may cache the tool and prompt lists.
const listTTLMs = 60 * 60 * 1000

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
	s, err := newServer(true, *addr == "")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if *addr != "" {
		fmt.Fprintf(os.Stderr, "Starting Barcode MCP Server on %s: Streamable HTTP at /mcp, legacy SSE at /sse\n", *addr)
		// SSE is the 2024-11-05 transport: its clients speak 2025-11-25 at
		// most, where the Skills extension does not exist — they get a
		// server that does not announce it.
		legacy, err := newServer(false, false)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		if err := newHTTPServer(*addr, s, legacy).ListenAndServe(); err != nil {
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

// newServer builds the MCP server with every tool and prompt. withSkills
// adds the Skills extension (protocol 2026-07-28); allowPath lets
// decode_barcode read files by path (stdio only).
func newServer(withSkills, allowPath bool) (*mcp.Server, error) {
	caps := &mcp.ServerCapabilities{
		Tools:   &mcp.ToolCapabilities{ListChanged: true},
		Prompts: &mcp.PromptCapabilities{ListChanged: true},
	}
	if withSkills {
		declareSkills(caps)
	}
	s := mcp.NewServer(
		&mcp.Implementation{
			Name:    "mlc-barcode-server",
			Version: version.Version,
		},
		&mcp.ServerOptions{
			Instructions: serverInstructions,
			Capabilities: caps,
			// Tools and prompts are fixed when the server is built, so
			// clients may keep the lists (and server/discover) for a while
			// instead of fetching them again for every use (SEP-2549);
			// list_changed would still announce a change.
			SetCacheable: func(_ context.Context, _ mcp.Request, c *mcp.Cacheable) {
				if c.TTLMs == 0 {
					c.TTLMs = listTTLMs
				}
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
	registerDecodeTools(s, allowPath)
	registerPrompts(s)
	if withSkills {
		if err := registerSkills(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}
