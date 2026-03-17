# MCP Server Integration

Dieses Dokument beschreibt die Integration und Konfiguration des MLC Barcode MCP Servers.

## Überblick

Der `mcp-barcode-server` stellt eine Model Context Protocol (MCP) Schnittstelle für Large Language Models (LLMs) bereit, um verschiedene Barcode-Typen zu generieren. Er kann im **stdio** Modus (Standard) oder im **SSE** Modus betrieben werden.

## Installation

1. Projekt bauen:
   ```bash
   make build
   ```
2. Die Binärdatei befindet sich unter `bin/mcp-barcode-server`.

## Konfiguration

### Parameter

- `-addr`: (Optional) Adresse für SSE (z. B. `:8080`). Wenn leer, wird stdio verwendet.
- `-artifact-addr`: (Optional) Die gRPC-Adresse des [mlcartifact Servers](https://github.com/hmsoft0815/mlcartifact) (z. B. `localhost:9590`).
- `-version`: Version anzeigen und beenden.

### Tool: `generate_barcode`

Der Server stellt ein einziges Tool namens `generate_barcode` zur Verfügung.

**Basis-Parameter:**
- `type`: Barcode-Typ (`qr`, `datamatrix`, `code128`, `code39`, `ean13`, `ean8`, `upca`, `itf`).
- `data`: Der zu kodierende Text.
- `format`: `svg` (Standard) oder `png`.
- `width` / `height`: Optionale Abmessungen.
- `fg_color` / `bg_color`: Farben (z. B. `black`, `#ff0000`, `transparent`).

**Artifact-Parameter (Nur verfügbar, wenn `-artifact-addr` beim Start angegeben wurde):**
- `save_artifact`: (Boolean) Wenn true, wird der generierte Barcode an den Artifact-Server gesendet.
- `filename`: (String) Optionaler Dateiname für den Artifact-Speicher.

## Integrationsbeispiele

### Claude Desktop (Stdio)

Ergänzen Sie Ihre `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "mlc-barcode": {
      "command": "/absoluter/pfad/zu/bin/mcp-barcode-server",
      "args": ["-artifact-addr", "localhost:9590"]
    }
  }
}
```

### SSE Modus

Server als Hintergrundprozess starten:
```bash
./bin/mcp-barcode-server -addr :8080 -artifact-addr localhost:9590
```

## Abhängigkeit: MLC Artifact Server

Die Artifact-Integration ist **optional**. Wenn Sie die Funktion `save_artifact` nutzen möchten, muss eine Instanz des [MLC Artifact Servers](https://github.com/hmsoft0815/mlcartifact) laufen und via gRPC erreichbar sein.

Wird beim Start keine `-artifact-addr` angegeben, bietet der MCP-Server die Artifact-bezogenen Parameter im Tool-Schema **nicht** an, um die Schnittstelle sauber zu halten.
