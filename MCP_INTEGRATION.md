# MCP Server Integration

This document describes how to integrate and configure the MLC Barcode MCP Server.

## Overview

The `mcp-barcode-server` provides a Model Context Protocol (MCP) interface for Large Language Models (LLMs) to generate various types of barcodes. It can be run in **stdio** mode (standard) or **SSE** mode.

## Installation

1. Build the project:
   ```bash
   make build
   ```
2. The binary will be located at `bin/mcp-barcode-server`.

## Configuration

### Parameters

- `-addr`: (Optional) Listen address for SSE (e.g., `:8080`). If empty, the server uses stdio.
- `-artifact-addr`: (Optional) The gRPC address of the [mlcartifact server](https://github.com/hmsoft0815/mlcartifact) (e.g., `localhost:9590`).
- `-version`: Show version and exit.

### Tool: `generate_barcode`

The server provides a single tool called `generate_barcode`.

**Base Parameters:**
- `type`: Barcode type (`qr`, `datamatrix`, `code128`, `code39`, `ean13`, `ean8`, `upca`, `itf`).
- `data`: The string to encode.
- `format`: `svg` (default) or `png`.
- `width` / `height`: Optional dimensions.
- `fg_color` / `bg_color`: Colors (e.g., `black`, `#ff0000`, `transparent`).

**Artifact Parameters (Only available if `-artifact-addr` is specified):**
- `save_artifact`: (Boolean) If true, the generated barcode will be sent to the artifact server.
- `filename`: (String) Optional filename for the artifact storage.

## Integration Examples

### Claude Desktop (Stdio)

Add the following to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "mlc-barcode": {
      "command": "/absolute/path/to/bin/mcp-barcode-server",
      "args": ["-artifact-addr", "localhost:9590"]
    }
  }
}
```

### SSE Mode

Run the server as a background process:
```bash
./bin/mcp-barcode-server -addr :8080 -artifact-addr localhost:9590
```

## Dependency: MLC Artifact Server

The artifact integration is **optional**. If you want to use the `save_artifact` feature, you must have an instance of the [MLC Artifact Server](https://github.com/hmsoft0815/mlcartifact) running and accessible via gRPC.

If no `-artifact-addr` is provided at startup, the MCP server will **not** offer the artifact-related parameters in its tool schema to keep the interface clean.
