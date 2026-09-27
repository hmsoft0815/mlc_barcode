# MCP Server Integration

This document describes how to integrate and configure the MLC Barcode MCP Server.

## Overview

The `mcp-barcode-server` provides a Model Context Protocol (MCP) interface for Large Language Models (LLMs) to generate various types of barcodes and specialized QR codes.

## Installation

1. Build the project:
   ```bash
   task build
   ```
2. The binary will be located at `bin/mcp-barcode-server`.

## Configuration Parameters

- `-addr`: (Optional) Listen address for HTTP (e.g., `:8080`): **Streamable HTTP** at `/mcp` (stateless, spec 2026-07-28), legacy SSE at `/sse` for older clients. Cross-origin browser requests are rejected. If empty, the server uses stdio.
- `-artifact-addr`: (Optional) The gRPC address of the [mlcartifact server](https://github.com/hmsoft0815/mlcartifact) (e.g., `localhost:9590`).
- `-version`: Show version and exit.

## Available Tools

All tools share common optional parameters for formatting:
- `format`: `svg` (default) or `png`.
- `width` / `height`: Optional dimensions.
- `fg_color` / `bg_color`: Colors (e.g., `black`, `#ff0000`, `transparent`).
- `text`: Show a caption below the barcode (the encoded content).
- `caption`: Your own caption text instead of the encoded content (implies `text`).
- `quiet_zone`: Blank margin scanners need (default `true`); the requested size includes it.
- `font_size`: Caption size in pixels (6–200); omitted = automatic. Over-long captions shrink to the barcode width.
- `save_artifact` / `filename`: Only available if `-artifact-addr` is specified.

### 1. `generate_barcode`
Generates a standard barcode.
- **Required**: `type`, `data`.
- `type`: `qr`, `datamatrix`, `aztec`, `pdf417`, `code128`, `code39`, `ean13`, `ean8`, `upca`, `itf`.
- Errors name the offending character or the amount of data and say what to do instead (e.g. `code128` for lower case, `aztec` for more text than a QR code holds).
- `ean13` / `ean8` / `upca`: digits only; with check digit (13/8/12 digits) or without (12/7/11 — it is computed). A wrong check digit is rejected, the error names the expected one.

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
- **Required**: `summary`, `start_time` (YYYYMMDDTHHMMSS, or YYYYMMDD for an all-day event; `end_time` is then the day after the last day).
- **Optional**: `end_time`, `description`, `location`, `timezone` (e.g. Europe/Berlin), `latitude`, `longitude`, `organizer` (as shown, e.g. "Lechner, Michael"), `organizer_email`.

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

## Reading codes: `decode_barcode`

Reads barcodes and QR codes from an image and reports symbology, content and corner points — several codes per image.
- **Input:** `path` (absolute path of a PNG, JPEG, GIF or WebP file — **only over stdio**; over HTTP the server would otherwise hand its own files to any client) or `image_base64` (a `data:` URI prefix is accepted). At most 20 MB and 40 megapixels. HEIC (iPhone photos) is not supported — export as JPEG.
- **Reads:** `qr`, `datamatrix`, `aztec`, `pdf417`, `ean13`, `ean8`, `upca`, `code128`, `code39`, `itf`.
- **Output:** `structuredContent` `{count, codes: [{barcode_type, text, content, points}]}`; nothing found is a tool error with hints.
- **`content`:** known payloads split into fields — `epc` (GiroCode: name, iban, `iban_valid` from the mod-97 check, bic, amount, currency, reference), `wifi`, `vcard`, `event`, `geo`, `tel`, `sms`, `email`, `crypto`, `url`, `pharma` (securPharm pack code, GS1 or IFA format: pzn, `pzn_valid`, gtin/ppn with check, batch, expiry YYMMDD, serial). Omitted for plain text.
- CLI: `barcode -decode image.png` prints `type<TAB>content` per code.

## Structured output

Every successful call returns, next to the image, `structuredContent` matching the tool's output schema:
`barcode_type`, `format`, `mime_type`, `encoded_data` (exactly what the code contains — the EAN/UPC with its check digit, the generated vCard, iCalendar or EPC text) and `artifact_id` when saved to mlcartifact.

## Prompts

- `payment_qr_from_invoice` (`invoice`, required): extract beneficiary, IBAN, amount and reference from invoice text and call `generate_epc_qr`; asks instead of guessing when the IBAN or amount is unclear.
- `product_labels` (`items`, required; `type`, default `ean13`): one labelled barcode per line `code;label text`.

## Skills

The server publishes **Agent Skills** (Skills extension `io.modelcontextprotocol/skills`, protocol 2026-07-28) — how-tos a client loads when a task needs one. Each file is a `skill://<name>/SKILL.md` resource; `skills/list` and `skills/get` return manifests with SHA-256 digests.

| Skill | For |
|---|---|
| `event-invite` | a calendar QR code for an invitation or appointment: time zone, place, organizer, all-day events |
| `wifi-guest-access` | Wi-Fi for guests: security type, special characters, guest networks |
| `contact-card` | a vCard from a signature or business card |
| `girocode-from-invoice` | a GiroCode from an invoice — never guessing IBAN or amount, checking the IBAN by reading the code back |
| `medicine-pack-check` | reading a medicine pack from a photo: PZN, batch, expiry (end of month for day 00), expired or not |
| `label-sheet` | barcodes for a list: which symbology for what, captions, sizes for printing |

Only on the current transports (stdio, Streamable HTTP `/mcp`); the legacy SSE endpoint speaks 2025-11-25, where the extension does not exist, and does not announce it. Checked with `mcp-tester skills --verify` in `task test:mcp`.

## Compliance

Checked with [mlc mcp-tester](https://github.com/hmsoft0815/mlc_mcptester) against spec 2026-07-28: `task test:mcp`.

## Dependency: MLC Artifact Server

The artifact integration is **optional**. If no `-artifact-addr` is provided at startup, the MCP server will **not** offer the artifact-related parameters in its tool schema to keep the interface clean.
