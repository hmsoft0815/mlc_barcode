// Package version carries the build identity of the CLI and the MCP server.
package version

// Version is the release this binary was built from.
//
// A var, not a const, and deliberately so: the value comes from the VERSION
// file at build time via `-ldflags -X`, which is the single source of truth
// (see Taskfile.yml). The literal here is only what an unstamped `go build`
// or a `go test` gets, and "dev" says exactly that rather than claiming a
// release number it may not be.
var Version = "dev"

// Author is fixed and has no reason to come from the build.
const Author = "Michael Lechner"
