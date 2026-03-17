# MCP Server Integration

This document describes how to integrate and configure the MLC Barcode MCP Server.

## Overview

The `mcp-barcode-server` provides a Model Context Protocol (MCP) interface for Large Language Models (LLMs) to generate various types of barcodes and specialized QR codes.

## Installation

1. Build the project:
   ```bash
   make build
   ```
2. The binary will be located at `bin/mcp-barcode-server`.

## Configuration Parameters

- `-addr`: (Optional) Listen address for SSE (e.g., `:8080`). If empty, the server uses stdio.
- `-artifact-addr`: (Optional) The gRPC address of the [mlcartifact server](https://github.com/hmsoft0815/mlcartifact) (e.g., `localhost:9590`).
- `-version`: Show version and exit.

## Available Tools

All tools share common optional parameters for formatting:
- `format`: `svg` (default) or `png`.
- `width` / `height`: Optional dimensions.
- `fg_color` / `bg_color`: Colors (e.g., `black`, `#ff0000`, `transparent`).
- `save_artifact` / `filename`: Only available if `-artifact-addr` is specified.

### 1. `generate_barcode`
Generates a standard barcode.
- **Required**: `type`, `data`.

### 2. `generate_wifi_qr`
Generates a QR code for WIFI access.
- **Required**: `ssid`.
- **Optional**: `password`, `encryption` (WPA/WEP/nopass), `hidden`.

### 3. `generate_vcard_qr`
Generates a QR code for a vCard 3.0 contact (referencing [RFC 6350](https://www.rfc-editor.org/rfc/rfc6350.html)).
- **Required**: `first_name`, `last_name`.
- **Optional**: `org`, `title`, `phone`, `email`, `address`, `city`, `zip`, `country`, `url`.

### 4. `generate_event_qr`
Generates a QR code for an iCalendar (RFC 5545) event.
- **Required**: `summary`, `start_time` (YYYYMMDDTHHMMSS).
- **Optional**: `end_time`, `description`, `location`, `timezone` (e.g. Europe/Berlin), `latitude`, `longitude`.

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

## Dependency: MLC Artifact Server

The artifact integration is **optional**. If no `-artifact-addr` is provided at startup, the MCP server will **not** offer the artifact-related parameters in its tool schema to keep the interface clean.
