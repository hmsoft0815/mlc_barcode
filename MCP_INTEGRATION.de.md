# MCP Server Integration

Dieses Dokument beschreibt die Integration und Konfiguration des MLC Barcode MCP Servers.

## Überblick

Der `mcp-barcode-server` stellt eine Model Context Protocol (MCP) Schnittstelle für Large Language Models (LLMs) bereit, um verschiedene Barcode-Typen und spezialisierte QR-Codes zu generieren.

## Installation

1. Projekt bauen:
   ```bash
   make build
   ```
2. Die Binärdatei befindet sich unter `bin/mcp-barcode-server`.

## Konfigurationsparameter

- `-addr`: (Optional) Adresse für SSE (z. B. `:8080`). Wenn leer, wird stdio verwendet.
- `-artifact-addr`: (Optional) Die gRPC-Adresse des [mlcartifact Servers](https://github.com/hmsoft0815/mlcartifact) (z. B. `localhost:9590`).
- `-version`: Version anzeigen und beenden.

## Verfügbare Tools

Alle Tools teilen sich gemeinsame optionale Parameter für die Formatierung:
- `format`: `svg` (Standard) oder `png`.
- `width` / `height`: Optionale Abmessungen.
- `fg_color` / `bg_color`: Farben (z. B. `black`, `#ff0000`, `transparent`).
- `save_artifact` / `filename`: Nur verfügbar, wenn `-artifact-addr` beim Start angegeben wurde.

### 1. `generate_barcode`
Generiert einen Standard-Barcode.
- **Erforderlich**: `type`, `data`.

### 2. `generate_wifi_qr`
Generiert einen QR-Code für den WLAN-Zugriff.
- **Erforderlich**: `ssid`.
- **Optional**: `password`, `encryption` (WPA/WEP/nopass), `hidden`.

### 3. `generate_vcard_qr`
Generiert einen QR-Code für einen vCard 3.0 Kontakt.
- **Erforderlich**: `first_name`, `last_name`.
- **Optional**: `org`, `title`, `phone`, `email`, `address`, `city`, `zip`, `country`, `url`.

### 4. `generate_event_qr`
Generiert einen QR-Code für einen iCalendar (RFC 5545) Termin.
- **Erforderlich**: `summary`, `start_time` (YYYYMMDDTHHMMSS).
- **Optional**: `end_time`, `description`, `location`, `timezone` (z. B. Europe/Berlin), `latitude`, `longitude`.

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

## Abhängigkeit: MLC Artifact Server

Die Artifact-Integration ist **optional**. Wird beim Start keine `-artifact-addr` angegeben, bietet der MCP-Server die Artifact-bezogenen Parameter im Tool-Schema **nicht** an, um die Schnittstelle sauber zu halten.
