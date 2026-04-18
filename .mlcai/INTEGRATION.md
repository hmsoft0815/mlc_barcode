# 🤖 AI \u0026 Integration Context: MLC Barcode — CLI \u0026 MCP Server

## 1. Identität \u0026 Zweck
- **Kernaufgabe:** Generierung verschiedenster Barcode-Typen (QR, DataMatrix, Code128, Code39, EAN-13, EAN-8, UPC-A, ITF) als Standalone CLI-Tool oder als Model Context Protocol (MCP) Server.
- **Technischer Stack:** Go 1.25.5, MCP-SDK (Stdio/SSE), `boombuler/barcode` f\u00fcr das Rendering.
- **Hoster/Infrastruktur:** CLI f\u00fcr lokale Ausf\u00fchrung, MCP-Server f\u00fcr Integration in LLM-Clients (z.B. Claude Desktop).

## 2. Die \"Nachbarschaft\" (System-Kontext)
- **Upstream (Wovon h\u00e4nge ich ab?):**
  - **boombuler/barcode** -\u003e [GitHub](https://github.com/boombuler/barcode) -\u003e Kern-Bibliothek f\u00fcr das Barcode-Rendering.
  - **mlcartifact** -\u003e [Internal/Optional] -\u003e Optionaler Service zum Speichern generierter Barcodes als persistente Artefakte.
- **Downstream (Wer nutzt mich?):**
  - **Claude Desktop / LLM-Clients** -\u003e Nutzt das MCP-Tool `generate_barcode`.
  - **DevOps/Scripts** -\u003e Nutzen das CLI-Tool `barcode` f\u00fcr automatisierte Asset-Generierung.
- **Shared Resources:**
  - Kommuniziert optional mit der `mlcartifact` API (default: localhost:9590).

## 3. Schnittstellen-Vertrag
- **Prim\u00e4re API:** Model Context Protocol (MCP) via Stdio oder SSE.
- **Auth-Mechanismus:** Keiner (lokaler Stdio-Prozess oder offenes SSE).
- **MCP-Tool:**
  - `generate_barcode`: { type: string, data: string, filename?: string, save_artifact?: boolean, ... }
- **CLI-Interface:** Umfangreiche Flags f\u00fcr Farben, Gr\u00f6\u00dfe, Text-Einblendung und strukturierte QR-Daten (WIFI, vCard, Event).
- **API-Doku-Link:** Siehe [README.md](../README.md) und [docs/qr-formats/](../docs/qr-formats/) f\u00fcr spezielle Formate.

## 4. Leitplanken \u0026 Regeln
- **Naming:** Go-Standard (CamelCase f\u00fcr exportierte Typen/Funktionen).
- **Testing:** Unit-Tests in `internal/` (z.B. `barcode_test.go`, `formats_test.go`).
- **Sicherheit:** Keine Speicherung sensibler Daten (WIFI-Passw\u00f6rter etc. landen nur im generierten Barcode).

## 5. Aktueller Fokus (Status)
- **Bekannte Probleme:** Keine kritischen Probleme; Fokus liegt auf Stabilit\u00e4t f\u00fcr LLM-Workflows.
- **Status:** Stabil (Version 1.2.0), aktuell keine Erweiterungen geplant.

---

## 📋 Meta

- **Zuletzt aktualisiert:** 2026-04-18
- **Aktualisiert von:** gemini-2.0-flash-exp
- **Status:** Aktuell
