package main

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func decodeSession(t *testing.T, allowPath bool) *mcp.ClientSession {
	t.Helper()
	ctx := context.Background()
	s := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0"}, nil)
	registerDecodeTools(s, allowPath)
	st, ct := mcp.NewInMemoryTransports()
	if _, err := s.Connect(ctx, st, nil); err != nil {
		t.Fatal(err)
	}
	cs, err := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "0"}, nil).Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { cs.Close() })
	return cs
}

// Over HTTP any client could read the server's files, so decode_barcode
// must neither offer nor honour "path" there.
func TestDecodeToolHidesPathOverHTTP(t *testing.T) {
	for _, allow := range []bool{true, false} {
		cs := decodeSession(t, allow)
		tools, err := cs.ListTools(context.Background(), nil)
		if err != nil {
			t.Fatal(err)
		}
		schema, _ := json.Marshal(tools.Tools[0].InputSchema)
		if got := strings.Contains(string(schema), `"path"`); got != allow {
			t.Errorf("allowPath=%v: schema offers path = %v: %s", allow, got, schema)
		}
	}

	cs := decodeSession(t, false)
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "decode_barcode", Arguments: map[string]any{"path": "/etc/hostname"},
	})
	if err == nil && (res == nil || !res.IsError) {
		t.Error("path was honoured although file access is off")
	}
}
