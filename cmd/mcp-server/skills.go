package main

import (
	"embed"
	"fmt"
	"io/fs"

	"github.com/hmsoft0815/mlc_mcptester/pkg/mcpskills"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Agent Skills the server publishes (Skills extension,
// io.modelcontextprotocol/skills): how-tos a client loads when a task needs
// one — an event invitation, Wi-Fi for guests, a contact card, a GiroCode
// from an invoice, checking a medicine pack, a label sheet. They are built
// into the binary; each file is a skill://<name>/<file> resource.
//
// The official go-sdk does not support the extension yet; mcpskills from
// mlc mcp-tester serves it (manifests with SHA-256 digests, skills/list,
// skills/get, resources/directory/read).
//
//go:embed skills
var skillFiles embed.FS

// declareSkills adds the extension to the capabilities (before the server
// is created).
func declareSkills(caps *mcp.ServerCapabilities) {
	mcpskills.Declare(caps, true)
}

// registerSkills serves the embedded skills.
func registerSkills(s *mcp.Server) error {
	sub, err := fs.Sub(skillFiles, "skills")
	if err != nil {
		return err
	}
	if _, err := mcpskills.Serve(s, sub, nil); err != nil {
		return fmt.Errorf("skills: %w", err)
	}
	return nil
}
