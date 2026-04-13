# MLC Barcode

Ein leistungsstarkes Werkzeug zur Generierung von Barcodes, das sowohl als eigenständiges Kommandozeilen-Tool (CLI) als auch als Model Context Protocol (MCP) Server genutzt werden kann.

## Was es macht

MLC Barcode wurde entwickelt, um die Barcode-Generierung nahtlos in Ihre täglichen Workflows und automatisierten Pipelines zu integrieren. Ob Sie einen einfachen QR-Code über das Terminal generieren, Etiketten automatisieren oder Ihre bevorzugten Large Language Models mit der Fähigkeit ausstatten möchten, Barcodes direkt im Chat zu erstellen – MLC Barcode ist die Lösung.

## Wichtigste Funktionen

- **Vielseitige Barcode-Unterstützung**:  Generieren Sie QR, DataMatrix, Code128, Code39, EAN-13, EAN-8, UPC-A und ITF.
- **LLM-Integration über MCP**: Der integrierte MCP-Server stellt das `generate_barcode`-Tool für KI-Clients wie Claude Desktop bereit, wodurch Barcodes mittels natürlicher Sprache generiert werden können.
- **Artifact-Anbindung**: Generierte Bilder können direkt an den `mlcartifact`-Dienst gesendet werden, um den lokalen Arbeitsbereich übersichtlich zu halten.
- **Smarte QR-Code Vorlagen**: Eingebaute Unterstützung für die Erstellung strukturierter QR-Codes, wie z.B. automatisch formatierte WLAN-Zugangsdaten, vCards und Kalendereinträge.
- **Hochwertige Ausgabe**: Erstellen Sie pixelgenaue Vektorgrafiken (`SVG`) oder Standard-Rasterbilder (`PNG`) mit anpassbarer Größe und transparentem Hintergrund.

## Schnellstart

### Claude Desktop
Fügen Sie Folgendes zu Ihrer `claude_desktop_config.json` hinzu:

```json
{
  "mcpServers": {
    "mlc-barcode": {
      "command": "mlc_barcode"
    }
  }
}
```

### Gemini-CLI
Fügen Sie den Server zu Ihrer `~/.gemini/settings.json` hinzu:

```json
{
  "mcpServers": {
    "mlc-barcode": {
      "command": "mlc_barcode"
    }
  }
}
```

### MCP-Tester
Ein neues Profil hinzufügen:

```bash
mcp-tester profile add barcode -c "mlc_barcode"
```
