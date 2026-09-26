# MLC Barcode — Desktop GUI, CLI & MCP Server

> **[mlcgo.eu](https://mlcgo.eu)** — tools, libraries and manuals · [Product page](https://mlcgo.eu/products/mlc-barcode/)


A tool for generating barcodes and QR codes — as a desktop GUI, a Command Line Interface (CLI) and a Model Context Protocol (MCP) server.

**MCP compliance checked with [mlc mcp-tester](https://github.com/hmsoft0815/mlc_mcptester)** — spec 2026-07-28, quality score 100/100, see [MCP compliance](#mcp-compliance).

<img src="assets/qrcode-product-teaser-960.jpg" alt="MLC Barcode — Desktop GUI, CLI and MCP Server" width="960">

<img src="assets/mlc_barcode_mpc3.png" >

## Why this exists

It started with a plain need: an MCP server that lets an AI assistant generate
barcodes and QR codes — EAN labels, Wi-Fi access, payment codes, business
cards — without anyone hand-assembling the underlying formats. That server is
still the core of the project.

The desktop app came on top. It turned out to be handy for quick tests, and for
everyone who is not fluent at the keyboard or has no LLM at hand: type, see the
code, print the label sheet.

### Why open source — and why with attribution

I believe in open source. Code you can read is code you can trust: you can see
what it does with your data, build it yourself and keep it running when I no
longer do.

Open does not mean ownerless, though. Taking the work and passing it off as your
own is not okay. That is why the licence is **MIT with an attribution clause**:
use it freely, change it, ship it — also commercially — as long as your product
visibly credits the author. If that credit does not fit your product, a
commercial licence without it is available on request. See [LICENSE](LICENSE).

## Version
Current Version: **1.5.0**

## Features

- Supports multiple barcode types: QR, DataMatrix, Aztec, PDF417, Code128, Code39, EAN-13, EAN-8, UPC-A, ITF.
- Output formats: SVG (vector-based) and PNG.
- Adjustable size and optional caption (custom text, adjustable font size) in SVG and PNG.
- Desktop GUI with batch generator and A4 label sheets.
- **Reads** barcodes from images (CLI `-decode`, MCP `decode_barcode`): symbology and content, several codes per image — every generated code is tested against this reader.

### File import (GUI: batch generator and label printer)

- **TXT:** one line = one barcode content; nothing is split.
- **CSV:** column 1 = content (e.g. for the QR code), column 2 = *optional* label text. The batch generator only uses column 1.
- Separator `;` or tab is detected; a comma only counts when every line has the same number of them. Without a detectable separator every line is one content.
- Quote a field that contains the separator itself, e.g. `"WIFI:T:WPA;S:Guest;P:secret;;";Guest Wi-Fi`. Empty lines are skipped, a UTF-8 BOM (Excel) is removed.

```csv
ART-1001;M4 screw
ART-1002
"WIFI:T:WPA;S:Guest;P:secret;;";Guest Wi-Fi
```
- MCP Server integration for LLMs (provides the `generate_barcode` tool).
- **Optional Artifact connection**: Generated barcodes can be sent directly to the `mlcartifact` service.

## Installation

Make sure you have Go installed.

```bash
git clone <repository-url>
cd mlc_barcode
task build
```

The binaries are located in `bin/`:

- `barcode`: CLI tool
- `mcp-barcode-server`: MCP server

## CLI Usage

```bash
# Show version
./bin/barcode -version

# Generate a QR code as SVG
./bin/barcode -type qr -data "Hello World" -out test.svg

# Optional: Save to artifact service
./bin/barcode -type qr -data "Hello World" -out test.png -artifact -artifact-addr localhost:9590
```

### Parameters

- `-decode <image>`: Read codes from a PNG/JPEG/GIF/WebP file instead of generating; prints `type<TAB>content` per code.
- `-type`: Barcode type (default: `qr`).
- `-no-quiet-zone`: Omit the blank margin around the code (on by default: QR 4 modules, DataMatrix 1, 1D codes ~10 left/right). Only if you add your own margin — without it DataMatrix and ITF fail on coloured or busy backgrounds.
- `-data`: Data to encode (required if no structured flags are used).
- `-out`: Output file with `.svg` or `.png` extension (default: `barcode.svg`).
- `-width`: Width in pixels.
- `-height`: Height in pixels.
- `-text`: Show text below barcode (default: `false`).
- `-fg` / `-bg`: Foreground and Background colors (e.g., `red`, `#ff0000`, `transparent`).
- `-version`: Show version and exit.

#### Structured QR Flags (Automatic formatting)
- **WIFI**: `-wifi-ssid`, `-wifi-pass`, `-wifi-enc` (WPA/WEP/nopass). [Documentation](docs/qr-formats/wifi.md)
- **vCard**: `-vcard-first`, `-vcard-last`, `-vcard-email`, `-vcard-tel`. [Documentation](docs/qr-formats/vcard.md)
- **Event**: `-event-summary`, `-event-start` (YYYYMMDDTHHMMSS), `-event-end`, `-event-tz`. [Documentation](docs/qr-formats/event.md)
- **GiroCode (EPC)**: `-epc-name`, `-epc-iban`, `-epc-bic`, `-epc-amount`, `-epc-ref`. [Documentation](docs/qr-formats/epc.md)
- **Crypto**: `-crypto-coin`, `-crypto-addr`, `-crypto-amount`, `-crypto-msg`. [Documentation](docs/qr-formats/crypto.md)
- **Geo**: `-geo-lat`, `-geo-lon`, `-geo-query`. [Documentation](docs/qr-formats/geo.md)
- **Telephone**: `-tel`. [Documentation](docs/qr-formats/tel.md)
- **SMS**: `-sms`, `-sms-text`. [Documentation](docs/qr-formats/sms.md)
- **Email**: `-mail-to`, `-mail-subject`, `-mail-body`. [Documentation](docs/qr-formats/email.md)

## Example output
<img src="showcase/assets/qr.png" > <img src="showcase/assets/ean13.svg" >

Detailed examples can be found in the **[Showcase](showcase/SHOWCASE.md)**.

## MCP Server Usage
<img src="assets/mlc_barcode_mpc4.png" >


The server supports stdio (default) and, with `-addr`, Streamable HTTP at `/mcp` plus legacy SSE at `/sse`.

### Claude Desktop Integration (Stdio)

Add this to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "mlc-barcode": {
      "command": "/path/to/mlc_barcode/bin/mcp-barcode-server",
      "args": ["-artifact-addr", "localhost:9590"]
    }
  }
}
```

The MCP tool `generate_barcode` has additional parameters:

- `save_artifact` (boolean): If true, saves the barcode to the artifact service.
- `filename` (string): Optional filename in the artifact store.

### HTTP mode (Streamable HTTP, SSE)

```bash
./bin/mcp-barcode-server -addr :8080 -artifact-addr localhost:9590
```

Clients connect to `http://<host>:8080/mcp` (Streamable HTTP, stateless as spec 2026-07-28 requires); older SSE-only clients to `http://<host>:8080/sse`. Cross-origin browser requests are rejected; there is no built-in authentication, so run it on a trusted network or behind an authenticating proxy.

Step-by-step setup for Claude Desktop, Claude Code, Gemini CLI and Cursor: **[MCP server guide on mlcgo.eu](https://mlcgo.eu/products/mlc-barcode/en/mcp/)**.

### MCP compliance

The server is checked against the MCP specification **2026-07-28** with **[mlc mcp-tester](https://github.com/hmsoft0815/mlc_mcptester)**, which looks at a server the way a client does — over the real transport, not by calling handlers directly:

- `inspect`: protocol version, capabilities, tool and prompt metadata, instructions — **quality score 100/100**.
- Protocol test script [`tests/mcp/barcode.mcp`](tests/mcp/barcode.mcp): every tool and symbology, PNG/SVG, captions, check-digit completion and rejection, empty and schema-violating input (reported as tool errors the model can correct), unknown tools (`-32602`), and `structuredContent` validated against each tool's output schema — the check strict clients such as the official TypeScript SDK apply.
- Both prompts (`payment_qr_from_invoice`, `product_labels`) via `prompts/get`, including the `-32602` error for a missing argument.
- `http-check` on the Streamable HTTP endpoint: required headers and their Base64 form, error codes and HTTP status, Origin validation, no GET/DELETE, no sessions — plus `inspect` over Streamable HTTP (100/100) and legacy SSE (90/100: SSE negotiates 2025-11-25 at most).

Run it yourself (needs `mcp-tester` in `PATH` or a checkout next to this repo):

```bash
task test:mcp
```

## Development

- `task build`: Compiles everything.
- `task dev:server`: Starts the MCP server via stdio.
- `task clean`: Clean up.
- `task test`: Run unit tests.
- `task test:mcp`: MCP compliance check with mlc mcp-tester.

## Reference

The **[MCP Handbook](https://mlcgo.eu/books/mcp-handbuch/)** explains the Model Context Protocol from the ground
up — tools, resources, prompts, transports, security and the artifact pattern.
Available in English and German.

## Barcode & 2D code specifications

The official ISO standards, specifications and documentation for the formats already implemented and for those planned.

### Integrated formats

| Short | Long form | Official specification & documentation |
| :--- | :--- | :--- |
| **QR** | Quick Response Code | https://iso.org (ISO/IEC 18004) |
| **DataMatrix** | DataMatrix ECC 200 | https://iso.org (ISO/IEC 16022) |
| **Code128** | Code 128 | https://iso.org (ISO/IEC 15417) |
| **Code39** | Code 39 (3 of 9) | https://iso.org (ISO/IEC 16388) |
| **EAN-13** | International Article Number (13 digits) | https://iso.org (ISO/IEC 15420) |
| **EAN-8** | International Article Number (8 digits) | https://iso.org (ISO/IEC 15420) |
| **UPC-A** | Universal Product Code | https://iso.org (ISO/IEC 15420) |
| **ITF** | Interleaved 2 of 5 | https://iso.org (ISO/IEC 16390) |
| **Aztec** | Aztec Code | https://iso.org (ISO/IEC 24778) |
| **PDF417** | Portable Data File 417 | https://iso.org (ISO/IEC 15438) |

### Planned formats

| Short | Long form | Official specification & documentation |
| :--- | :--- | :--- |
| **GS1 DataMatrix** | GS1 DataMatrix (with FNC1) | https://gs1.org (GS1 Guideline) |

### Helpful open resources (implementation aids)

ISO standards are not free; these open sources document the algorithms for check digits and the character-set tables well:

* **Check digits & character sets:** https://grandzebu.net (extensive calculations for Code 128 and EAN)
* **GS1 data structures:** https://gs1.org (definitions of all Application Identifiers)

---

## License

Copyright (c) 2026 Michael Lechner.
Licensed under the MIT License with an attribution clause — free to use, modify and
redistribute, also commercially, as long as products using it credit "Michael Lechner"
visibly (docs, About screen or similar). A commercial licence without attribution is
available on request. See [LICENSE](LICENSE).

---
**Note:** Currently, no further expansions or major changes are planned, as the tools have proven themselves effective in everyday use, particularly in combination with Large Language Models (LLMs).

<!-- mlcai-private -->
## Project documentation (`.mlcai/`)

`.mlcai/` is a **private git submodule**: internal planning, backlog and work notes, maintained with the MLC Doc Hub. It is not publicly accessible — clone **without** `--recurse-submodules`; the build does not need it. Links into `.mlcai/` only work with access (`git submodule update --init .mlcai`).

## Who is "Claude" in the commits?

Some commits in this repository are co-authored by Claude, Anthropic's AI
model. It helps write code, keeps our documentation and backlog up to date and
digs through failing builds — every change is reviewed before it is merged.
We don't hide it: [how we work with Claude](https://mlcgo.eu/ai/en.html).
