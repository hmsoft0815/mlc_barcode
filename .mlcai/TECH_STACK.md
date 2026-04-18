# 🛠 Tech Stack \u0026 Constraints: MLC Barcode — CLI \u0026 MCP Server

## Kern-Versionen
- **Sprache:** Go 1.25.5.
- **Frameworks:** Keines (Standardbibliothek + MCP-SDK).
- **Architektur:** Go-idiomatic, flache Package-Struktur (`internal/`).

## Bibliotheken (Erlaubt/Fixiert)
- **Barcode-Rendering:** `github.com/boombuler/barcode` (v1.1.0) -\u003e Generiert SVG (vektorbasiert) und PNG.
- **MCP Integration:** `github.com/modelcontextprotocol/go-sdk` (v1.4.1).
- **Speicher-Anbindung:** `github.com/hmsoft0815/mlcartifact` (v0.4.1).

## Projektstruktur
- `cmd/barcode`: CLI-Einstiegspunkt (ca. 150 Zeilen).
- `cmd/mcp-server`: MCP-Server-Einstiegspunkt (ca. 340 Zeilen).
- `internal/barcodes`: Kern-Generierungslogik (ca. 300 Zeilen).
- `internal/qrformats`: Spezial-Formate f\u00fcr WIFI, vCard und vCalendar.
- `internal/version`: Verwaltung der Version (aktuell 1.2.0).

## Einschränkungen (Constraints)
- **Kein Gin/HTTP:** F\u00fcr HTTP-Anforderungen wird ausschlie\u00dflich das MCP-SDK (SSE) genutzt.
- **Plattform-unabh\u00e4ngig:** Pfade werden ausschlie\u00dflich via `path/filepath` gehandhabt.
- **Performance:** Barcode-Generierung erfolgt rein lokal ohne externe APIs (au\u00dfer optionaler Upload zu `mlcartifact`).

## Styling-Regeln
- **Early Returns:** Zur Reduzierung der zyklomatischen Komplexit\u00e4t bevorzugen (Guard Clauses).
- **Dokumentation:** Jede exportierte Funktion MUSS einen GoDoc-Kommentar enthalten.
- **Fehler-Handling:** Errors werden explizit behandelt und mit Kontext via `fmt.Errorf(\"context: %w\", err)` gewrappt.
\n## Build \u0026 Tooling\n- **Makefile:** Zentrales Build-Tool (Targets: `build`, `clean`, `test`, `build-linux`, `build-windows`, `build-macos`).\n- **Tests:** F\u00fcr alle kritischen Generierungsfunktionen in `internal/` vorhanden.\n"}

---

## 📋 Meta

- **Zuletzt aktualisiert:** 2026-04-18
- **Aktualisiert von:** gemini-2.0-flash-exp
- **Status:** Aktuell
