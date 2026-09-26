package main

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// newHTTPServer serves the MCP server over the network:
//
//   - /mcp  Streamable HTTP, stateless as spec 2026-07-28 requires (no
//     sessions, no GET/DELETE). The tools never call back into the client,
//     so nothing is lost.
//   - /sse  the SSE transport of spec 2024-11-05, kept for older clients.
//
// Cross-origin requests from browsers are rejected (the spec requires
// servers to validate Origin, against DNS rebinding). Base64-encoded Mcp-*
// headers are decoded first, see decodeMCPHeaders.
func newHTTPServer(addr string, s *mcp.Server) *http.Server {
	getServer := func(*http.Request) *mcp.Server { return s }

	mux := http.NewServeMux()
	mux.Handle("/mcp", mcp.NewStreamableHTTPHandler(getServer, &mcp.StreamableHTTPOptions{Stateless: true}))
	mux.Handle("/sse", mcp.NewSSEHandler(getServer, nil))

	return &http.Server{
		Addr:              addr,
		Handler:           http.NewCrossOriginProtection().Handler(decodeMCPHeaders(mux)),
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}
}

// decodeMCPHeaders decodes Mcp-* header values sent as =?base64?…?=, which
// spec 2026-07-28 allows for values that are not plain ASCII. go-sdk v1.8.0
// compares the raw header with the request body and answers such a request
// with -32020 (found with mcp-tester http-check). Drop this once the SDK
// decodes the values itself.
func decodeMCPHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for name, values := range r.Header {
			if !strings.HasPrefix(name, "Mcp-") {
				continue
			}
			for i, v := range values {
				if decoded, ok := decodeBase64Header(v); ok {
					values[i] = decoded
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

func decodeBase64Header(v string) (string, bool) {
	inner, ok := strings.CutPrefix(v, "=?base64?")
	if !ok {
		return "", false
	}
	inner, ok = strings.CutSuffix(inner, "?=")
	if !ok {
		return "", false
	}
	b, err := base64.StdEncoding.DecodeString(inner)
	if err != nil {
		if b, err = base64.RawStdEncoding.DecodeString(inner); err != nil {
			return "", false
		}
	}
	if !utf8.Valid(b) {
		return "", false
	}
	return string(b), true
}
