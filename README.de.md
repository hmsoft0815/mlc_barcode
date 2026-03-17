# MLC Barcode CLI & MCP Server

Ein Werkzeug zur Generierung von Barcodes, das sowohl als Kommandozeilen-Tool (CLI) als auch als Model Context Protocol (MCP) Server genutzt werden kann.

<img src="assets/mlc_barcode_mpc3.png" >

## Version
Aktuelle Version: **1.1.0**

## Funktionen
- Unterstützt mehrere Barcode-Typen: QR, DataMatrix, Code128, Code39, EAN-13, EAN-8, UPC-A, ITF.
- Ausgabeformate: SVG (vektorbasiert) und PNG.
- Anpassbare Größe und optionale Textanzeige für generierte SVG-Bilder.
- MCP-Server Integration für LLMs (bietet das Tool `generate_barcode` an).
- **Optionale Artifact-Anbindung**: Generierte Barcodes können direkt an den `mlcartifact` Dienst gesendet werden.
- Saubere Projektstruktur nach Go Best Practices.

## Installation

Stellen Sie sicher, dass Go installiert ist.

```bash
git clone <repository-url>
cd mlc_barcode
make build
```

Die Binärdateien befinden sich in `bin/`:
- `barcode`: CLI Werkzeug
- `mcp-barcode-server`: MCP Server

## Benutzung als CLI

```bash
# Version anzeigen
./bin/barcode -version

# Einen QR-Code als SVG generieren
./bin/barcode -type qr -data "Hallo Welt" -out test.svg

# Optional: Als Artifact speichern
./bin/barcode -type qr -data "Hallo Welt" -out test.png -artifact -artifact-addr localhost:9590
```

### Parameter

- `-type`: Barcode-Typ (Standard: `qr`).
- `-data`: Der zu kodierende Inhalt (erforderlich, sofern keine strukturierten Flags genutzt werden).
- `-out`: Ausgabedatei mit Endung `.svg` oder `.png` (Standard: `barcode.svg`).
- `-width`: Breite in Pixeln.
- `-height`: Höhe in Pixeln.
- `-text`: Text unter dem Barcode anzeigen (Standard: `false`).
- `-fg` / `-bg`: Vorder- und Hintergrundfarben (z. B. `red`, `#ff0000`, `transparent`).
- `-version`: Version anzeigen und beenden.

#### Strukturierte QR-Flags (Automatische Formatierung)
- **WLAN**: `-wifi-ssid`, `-wifi-pass`, `-wifi-enc` (WPA/WEP/nopass).
- **vCard**: `-vcard-first`, `-vcard-last`, `-vcard-email`, `-vcard-tel`.
- **Termin**: `-event-summary`, `-event-start` (YYYYMMDDTHHMMSS), `-event-end`, `-event-tz`.

## Beispielausgabe
<img src="showcase/assets/qr.png" > <img src="showcase/assets/ean13.svg" >

Detaillierte Beispiele finden Sie im **[Showcase](showcase/SHOWCASE.de.md)**.

## Benutzung als MCP Server

<img src="assets/mlc_barcode_mpc4.png" >

Der Server unterstützt Stdio (Standard) und SSE.

### Integration in Claude Desktop (Stdio)
Ergänzen Sie Ihre `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "mlc-barcode": {
      "command": "/pfad/zu/mlc_barcode/bin/mcp-barcode-server",
      "args": ["-artifact-addr", "localhost:9590"]
    }
  }
}
```

Das MCP-Tool `generate_barcode` hat zusätzliche Parameter:
- `save_artifact` (boolean): Wenn true, wird der Barcode via mlc_artifact gespeichert (benötigit mlc artifact mcp server) gespeichert.
- `filename` (string): Optionaler Dateiname im Artifact-Speicher.

### SSE Modus
```bash
./bin/mcp-barcode-server -addr :8080 -artifact-addr localhost:9590
```

## Entwicklung

- `make build`: Kompiliert alles.
- `make run-server`: Startet den MCP Server über stdio.
- `make clean`: Aufräumen.
- `make test`: Unit-Tests ausführen.

## Lizenz

Copyright (c) 2026 Michael Lechner.
Lizenziert unter der MIT-Lizenz.
