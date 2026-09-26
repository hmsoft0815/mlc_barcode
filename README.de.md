# MLC Barcode — Desktop GUI, CLI & MCP Server

> **[mlcgo.eu](https://mlcgo.eu)** — Werkzeuge, Bibliotheken und Handbücher · [Produktseite](https://mlcgo.eu/products/mlc-barcode/)


Ein Werkzeug zur Generierung von Barcodes und QR-Codes — als Desktop-GUI, als Kommandozeilen-Tool (CLI) und als Model Context Protocol (MCP) Server.

**MCP-Konformität geprüft mit [mlc mcp-tester](https://github.com/hmsoft0815/mlc_mcptester)** — Spezifikation 2026-07-28, Qualitäts-Score 100/100, siehe [MCP-Konformität](#mcp-konformität).

<img src="assets/qrcode-product-teaser-960.jpg" alt="MLC Barcode — Desktop GUI, CLI and MCP Server" width="960">

<img src="assets/mlc_barcode_mpc3.png" >

## Warum es das gibt

Am Anfang stand ein ganz praktischer Bedarf: ein MCP-Server, mit dem ein
KI-Assistent Barcodes und QR-Codes erzeugen kann — EAN-Etiketten, WLAN-Zugänge,
Zahlungscodes, Visitenkarten —, ohne dass jemand die zugrunde liegenden Formate
von Hand zusammenbaut. Dieser Server ist bis heute der Kern des Projekts.

Die Desktop-App kam obendrauf. Sie hat sich als praktisch erwiesen für schnelle
Tests und für alle, die an der Tastatur nicht so sattelfest sind oder gerade kein
LLM zur Hand haben: eintippen, Code sehen, Etikettenbogen drucken.

### Warum Open Source — und warum mit Namensnennung

Ich glaube an Open Source. Code, den man lesen kann, ist Code, dem man vertrauen
kann: Man sieht, was er mit den eigenen Daten macht, kann ihn selbst bauen und
weiter betreiben, auch wenn ich es irgendwann nicht mehr tue.

Offen heißt aber nicht herrenlos. Die Arbeit zu nehmen und als die eigene
auszugeben, ist nicht in Ordnung. Deshalb steht das Projekt unter der
**MIT-Lizenz mit Namensnennungsklausel**: frei nutzen, ändern, weitergeben —
auch kommerziell —, solange das eigene Produkt den Autor sichtbar nennt. Passt
diese Nennung nicht ins Produkt, gibt es auf Anfrage eine kommerzielle Lizenz
ohne sie. Siehe [LICENSE](LICENSE).

## Version
Aktuelle Version: **1.5.0**

## Funktionen
- Unterstützt mehrere Barcode-Typen: QR, DataMatrix, Aztec, PDF417, Code128, Code39, EAN-13, EAN-8, UPC-A, ITF.
- Ausgabeformate: SVG (vektorbasiert) und PNG.
- Anpassbare Größe und optionale Beschriftung (Freitext, einstellbare Schriftgröße) in SVG und PNG.
- Desktop-GUI mit Batch-Generator und A4-Etikettenbögen.
- **Liest** Barcodes aus Bildern (CLI `-decode`, MCP `decode_barcode`): Format und Inhalt, auch mehrere Codes pro Bild — jeder erzeugte Code wird gegen diesen Reader getestet.

### Datei-Import (GUI: Batch-Generator und Etiketten-Druck)

- **TXT:** Eine Zeile = ein Barcode-Inhalt; es wird nichts aufgeteilt.
- **CSV:** Spalte 1 = Inhalt (z. B. für den QR-Code), Spalte 2 = *optionaler* Etikett-Text. Der Batch-Generator nutzt nur Spalte 1.
- Trennzeichen `;` oder Tab werden erkannt; ein Komma gilt nur, wenn jede Zeile gleich viele hat. Ohne erkennbares Trennzeichen ist jede Zeile ein Inhalt.
- Enthält der Inhalt selbst das Trennzeichen, das Feld in Anführungszeichen setzen, z. B. `"WIFI:T:WPA;S:Gast;P:geheim;;";Gäste-WLAN`. Leere Zeilen werden übersprungen, ein UTF-8-BOM (Excel) wird entfernt.

```csv
ART-1001;Schraube M4
ART-1002
"WIFI:T:WPA;S:Gast;P:geheim;;";Gäste-WLAN
```
- MCP-Server Integration für LLMs (bietet das Tool `generate_barcode` an).
- **Optionale Artifact-Anbindung**: Generierte Barcodes können direkt an den `mlcartifact` Dienst gesendet werden.
- Saubere Projektstruktur nach Go Best Practices.

## Installation

Stellen Sie sicher, dass Go installiert ist.

```bash
git clone <repository-url>
cd mlc_barcode
task build
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

- `-decode <bild>`: Codes aus einer PNG/JPEG/GIF/WebP-Datei lesen statt erzeugen; gibt je Code `Typ<TAB>Inhalt` aus.
- `-type`: Barcode-Typ (Standard: `qr`).
- `-no-quiet-zone`: Ruhezone (Rand um den Code) weglassen — standardmäßig an: QR 4 Module, DataMatrix 1, 1D-Codes ~10 links/rechts. Nur abschalten, wenn Sie selbst Rand lassen; ohne sie scheitern DataMatrix und ITF auf farbigem oder unruhigem Grund.
- `-data`: Der zu kodierende Inhalt (erforderlich, sofern keine strukturierten Flags genutzt werden).
- `-out`: Ausgabedatei mit Endung `.svg` oder `.png` (Standard: `barcode.svg`).
- `-width`: Breite in Pixeln.
- `-height`: Höhe in Pixeln.
- `-text`: Text unter dem Barcode anzeigen (Standard: `false`).
- `-fg` / `-bg`: Vorder- und Hintergrundfarben (z. B. `red`, `#ff0000`, `transparent`).
- `-version`: Version anzeigen und beenden.

#### Strukturierte QR-Flags (Automatische Formatierung)
- **WLAN**: `-wifi-ssid`, `-wifi-pass`, `-wifi-enc` (WPA/WEP/nopass). [Dokumentation](docs/qr-formats/wifi.de.md)
- **vCard**: `-vcard-first`, `-vcard-last`, `-vcard-email`, `-vcard-tel`. [Dokumentation](docs/qr-formats/vcard.de.md)
- **Termin**: `-event-summary`, `-event-start` (YYYYMMDDTHHMMSS), `-event-end`, `-event-tz`. [Dokumentation](docs/qr-formats/event.de.md)
- **GiroCode (EPC)**: `-epc-name`, `-epc-iban`, `-epc-bic`, `-epc-amount`, `-epc-ref`. [Dokumentation](docs/qr-formats/epc.de.md)
- **Krypto**: `-crypto-coin`, `-crypto-addr`, `-crypto-amount`, `-crypto-msg`. [Dokumentation](docs/qr-formats/crypto.de.md)
- **Geo**: `-geo-lat`, `-geo-lon`, `-geo-query`. [Dokumentation](docs/qr-formats/geo.de.md)
- **Telefon**: `-tel`. [Dokumentation](docs/qr-formats/tel.de.md)
- **SMS**: `-sms`, `-sms-text`. [Dokumentation](docs/qr-formats/sms.de.md)
- **E-Mail**: `-mail-to`, `-mail-subject`, `-mail-body`. [Dokumentation](docs/qr-formats/email.de.md)

## Beispielausgabe
<img src="showcase/assets/qr.png" > <img src="showcase/assets/ean13.svg" >

Detaillierte Beispiele finden Sie im **[Showcase](showcase/SHOWCASE.de.md)**.

## Benutzung als MCP Server

<img src="assets/mlc_barcode_mpc4.png" >

Der Server unterstützt stdio (Standard) und mit `-addr` Streamable HTTP unter `/mcp` sowie das alte SSE unter `/sse`.

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

### HTTP-Modus (Streamable HTTP, SSE)
```bash
./bin/mcp-barcode-server -addr :8080 -artifact-addr localhost:9590
```

Clients verbinden sich mit `http://<host>:8080/mcp` (Streamable HTTP, zustandslos wie von Spezifikation 2026-07-28 verlangt), ältere Clients, die nur SSE können, mit `http://<host>:8080/sse`. Browser-Anfragen fremder Herkunft werden abgewiesen; eine eigene Anmeldung gibt es nicht — nur im vertrauenswürdigen Netz oder hinter einem Proxy mit Authentifizierung betreiben.

Schritt-für-Schritt-Anleitung für Claude Desktop, Claude Code, Gemini CLI und Cursor: **[MCP-Anleitung auf mlcgo.eu](https://mlcgo.eu/products/mlc-barcode/de/mcp/)**.

### MCP-Konformität

Der Server wird gegen die MCP-Spezifikation **2026-07-28** mit dem **[mlc mcp-tester](https://github.com/hmsoft0815/mlc_mcptester)** geprüft. Der Tester sieht den Server so, wie ein Client ihn sieht — über den echten Transport, nicht durch direkten Aufruf der Handler:

- `inspect`: Protokollversion, Capabilities, Metadaten von Werkzeugen und Prompts, Instructions — **Qualitäts-Score 100/100**.
- Protokoll-Testskript [`tests/mcp/barcode.mcp`](tests/mcp/barcode.mcp): jedes Werkzeug und jede Symbologie, PNG/SVG, Beschriftungen, Ergänzen und Ablehnen von Prüfziffern, leere und schemawidrige Eingaben (als Tool-Fehler, die das Modell korrigieren kann), unbekannte Werkzeuge (`-32602`) sowie `structuredContent`, geprüft gegen das Output-Schema jedes Werkzeugs — so, wie strikte Clients (etwa das offizielle TypeScript-SDK) es verlangen.
- Beide Prompts (`payment_qr_from_invoice`, `product_labels`) über `prompts/get`, einschließlich des Fehlers `-32602` bei fehlendem Argument.
- `http-check` am Streamable-HTTP-Endpunkt: Pflicht-Header samt Base64-Form, Fehlercodes und HTTP-Status, Origin-Prüfung, kein GET/DELETE, keine Sessions — dazu `inspect` über Streamable HTTP (100/100) und das alte SSE (90/100: SSE handelt höchstens 2025-11-25 aus).

Selbst ausführen (braucht `mcp-tester` im `PATH` oder einen Checkout neben diesem Repo):

```bash
task test:mcp
```

## Entwicklung

- `task build`: Kompiliert alles.
- `task dev:server`: Startet den MCP Server über stdio.
- `task clean`: Aufräumen.
- `task test`: Unit-Tests ausführen.
- `task test:mcp`: MCP-Konformitätsprüfung mit dem mlc mcp-tester.

## Referenz

Das **[MCP-Handbuch](https://mlcgo.eu/books/mcp-handbuch/)** erklärt das Model Context Protocol von Grund auf —
Tools, Resources, Prompts, Transporte, Sicherheit und das Artifact-Pattern.
Auf Deutsch und Englisch.

## Barcode- & 2D-Code-Spezifikationen

Diese Übersicht enthält die offiziellen ISO-Normen, Spezifikationen und Dokumentationen für die bereits implementierten sowie für geplante Formate.

### Integrierte Formate

| Kurz | Langform | Offizielle Spezifikation & Dokumentation |
| :--- | :--- | :--- |
| **QR** | Quick Response Code | https://iso.org (ISO/IEC 18004) |
| **DataMatrix** | DataMatrix ECC 200 | https://iso.org (ISO/IEC 16022) |
| **Code128** | Code 128 | https://iso.org (ISO/IEC 15417) |
| **Code39** | Code 39 (3 of 9) | https://iso.org (ISO/IEC 16388) |
| **EAN-13** | International Article Number (13-stellig) | https://iso.org (ISO/IEC 15420) |
| **EAN-8** | International Article Number (8-stellig) | https://iso.org (ISO/IEC 15420) |
| **UPC-A** | Universal Product Code | https://iso.org (ISO/IEC 15420) |
| **ITF** | Interleaved 2 of 5 | https://iso.org (ISO/IEC 16390) |
| **Aztec** | Aztec Code | https://iso.org (ISO/IEC 24778) |
| **PDF417** | Portable Data File 417 | https://iso.org (ISO/IEC 15438) |

### Geplante Formate

| Kurz | Langform | Offizielle Spezifikation & Dokumentation |
| :--- | :--- | :--- |
| **GS1 DataMatrix** | GS1 DataMatrix (mit FNC1-Steuerung) | https://gs1.org (GS1 Guideline) |

### Hilfreiche offene Ressourcen (Implementierungshilfen)

Da ISO-Normen kostenpflichtig sind, bieten die folgenden freien Dokumentationen exzellente Algorithmen für Prüfziffern-Berechnungen und Zeichensatz-Tabellen:

* **Prüfziffern & Zeichensätze:** https://grandzebu.net (umfangreiche Berechnungen für Code 128 und EAN)
* **GS1-Datenstrukturen:** https://gs1.org (Definitionen aller Application Identifier)

---

## Lizenz

Copyright (c) 2026 Michael Lechner.
Lizenziert unter der MIT-Lizenz mit Namensnennungsklausel — frei nutzbar, veränderbar und
weitergebbar, auch kommerziell, solange Produkte, die es verwenden, "Michael Lechner" sichtbar
nennen (Dokumentation, Info-Dialog o. ä.). Eine kommerzielle Lizenz ohne Namensnennung gibt es
auf Anfrage. Siehe [LICENSE](LICENSE).

---
**Hinweis:** Aktuell sind keine weiteren Erweiterungen oder größeren Änderungen geplant, da sich die Werkzeuge im Alltag – insbesondere auch in Verbindung mit Large Language Models (LLMs) – bestens bewährt haben.

<!-- mlcai-private -->
## Projektdokumentation (`.mlcai/`)

`.mlcai/` ist ein **privates Git-Submodul**: interne Planung, Backlog und Arbeitsnotizen, gepflegt mit dem MLC Doc Hub. Es ist nicht öffentlich zugänglich — **ohne** `--recurse-submodules` klonen; für den Build wird es nicht gebraucht. Links nach `.mlcai/` funktionieren nur mit Zugriff (`git submodule update --init .mlcai`).
