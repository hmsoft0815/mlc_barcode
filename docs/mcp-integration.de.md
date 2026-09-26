# MCP Server Integration

Dieses Dokument beschreibt die Integration und Konfiguration des MLC Barcode MCP Servers.

## Überblick

Der `mcp-barcode-server` stellt eine Model Context Protocol (MCP) Schnittstelle für Large Language Models (LLMs) bereit, um verschiedene Barcode-Typen und spezialisierte QR-Codes zu generieren.

## Installation

1. Projekt bauen:
   ```bash
   task build
   ```
2. Die Binärdatei befindet sich unter `bin/mcp-barcode-server`.

## Konfigurationsparameter

- `-addr`: (Optional) Adresse für HTTP (z. B. `:8080`): **Streamable HTTP** unter `/mcp` (zustandslos, Spezifikation 2026-07-28), altes SSE unter `/sse` für ältere Clients. Browser-Anfragen fremder Herkunft werden abgewiesen. Wenn leer, wird stdio verwendet.
- `-artifact-addr`: (Optional) Die gRPC-Adresse des [mlcartifact Servers](https://github.com/hmsoft0815/mlcartifact) (z. B. `localhost:9590`).
- `-version`: Version anzeigen und beenden.

## Verfügbare Tools

Alle Tools teilen sich gemeinsame optionale Parameter für die Formatierung:
- `format`: `svg` (Standard) oder `png`.
- `width` / `height`: Optionale Abmessungen.
- `fg_color` / `bg_color`: Farben (z. B. `black`, `#ff0000`, `transparent`).
- `text`: Beschriftung unter dem Barcode anzeigen (der codierte Inhalt).
- `caption`: Eigener Beschriftungstext statt des codierten Inhalts (schaltet `text` ein).
- `quiet_zone`: Rand, den Scanner brauchen (Standard `true`); die angeforderte Größe schließt ihn ein.
- `font_size`: Schriftgröße der Beschriftung in Pixeln (6–200); weggelassen = automatisch. Zu lange Beschriftungen werden auf Barcodebreite verkleinert.
- `save_artifact` / `filename`: Nur verfügbar, wenn `-artifact-addr` beim Start angegeben wurde.

### 1. `generate_barcode`
Generiert einen Standard-Barcode.
- **Erforderlich**: `type`, `data`.
- `type`: `qr`, `datamatrix`, `aztec`, `pdf417`, `code128`, `code39`, `ean13`, `ean8`, `upca`, `itf`.
- Fehlermeldungen nennen das störende Zeichen oder die Datenmenge und sagen, was stattdessen zu tun ist (z. B. `code128` für Kleinbuchstaben, `aztec` für mehr Text, als ein QR-Code fasst).
- `ean13` / `ean8` / `upca`: nur Ziffern; mit Prüfziffer (13/8/12 Stellen) oder ohne (12/7/11 — sie wird berechnet). Eine falsche Prüfziffer wird abgelehnt, die Fehlermeldung nennt die richtige.

### 2. `generate_wifi_qr`
Generiert einen QR-Code für den WLAN-Zugriff.
- **Erforderlich**: `ssid`.
- **Optional**: `password`, `encryption` (WPA/WEP/nopass), `hidden`.

### 3. `generate_vcard_qr`
Generiert einen QR-Code für einen vCard 3.0 Kontakt (unter Berücksichtigung von [RFC 6350](https://www.rfc-editor.org/rfc/rfc6350.html)).
- **Erforderlich**: `first_name`, `last_name`.
- **Optional**: `org`, `title`, `phone`, `email`, `address`, `city`, `zip`, `country`, `url`.

### 4. `generate_event_qr`
Generiert einen QR-Code für einen iCalendar (RFC 5545) Termin.
- **Erforderlich**: `summary`, `start_time` (YYYYMMDDTHHMMSS, oder YYYYMMDD für einen ganztägigen Termin; `end_time` ist dann der Tag nach dem letzten Tag).
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

## Codes lesen: `decode_barcode`

Liest Barcodes und QR-Codes aus einem Bild und meldet Format, Inhalt und Eckpunkte — auch mehrere Codes pro Bild.
- **Eingabe:** `path` (absoluter Pfad einer PNG-, JPEG-, GIF- oder WebP-Datei — **nur über stdio**; über HTTP würde der Server sonst seine eigenen Dateien an jeden Client geben) oder `image_base64` (ein `data:`-URI-Präfix ist erlaubt). Höchstens 20 MB und 40 Megapixel. HEIC (iPhone-Fotos) geht nicht — als JPEG exportieren.
- **Liest:** `qr`, `datamatrix`, `aztec`, `pdf417`, `ean13`, `ean8`, `upca`, `code128`, `code39`, `itf`.
- **Ausgabe:** `structuredContent` `{count, codes: [{barcode_type, text, content, points}]}`; wird nichts gefunden, gibt es einen Tool-Fehler mit Hinweisen.
- **`content`:** bekannte Inhalte in Feldern — `epc` (GiroCode: name, iban, `iban_valid` aus der Mod-97-Prüfung, bic, amount, currency, reference), `wifi`, `vcard`, `event`, `geo`, `tel`, `sms`, `email`, `crypto`, `url`, `pharma` (Arzneimittel-Code nach securPharm, GS1- oder IFA-Format: pzn, `pzn_valid`, gtin/ppn mit Prüfung, batch, expiry JJMMTT, serial). Bei reinem Text weggelassen.
- CLI: `barcode -decode bild.png` gibt je Code `Typ<TAB>Inhalt` aus.

## Strukturierte Ausgabe

Jeder erfolgreiche Aufruf liefert neben dem Bild `structuredContent` passend zum Output-Schema des Werkzeugs:
`barcode_type`, `format`, `mime_type`, `encoded_data` (genau das, was im Code steht — EAN/UPC mit Prüfziffer, die erzeugte vCard, der iCalendar- oder EPC-Text) und `artifact_id`, wenn in mlcartifact gespeichert wurde.

## Prompts

- `payment_qr_from_invoice` (`invoice`, Pflicht): Empfänger, IBAN, Betrag und Verwendungszweck aus einem Rechnungstext lesen und `generate_epc_qr` aufrufen; fragt nach, statt zu raten, wenn IBAN oder Betrag unklar sind.
- `product_labels` (`items`, Pflicht; `type`, Standard `ean13`): ein beschrifteter Barcode je Zeile `Code;Etikett-Text`.

## Konformität

Geprüft mit dem [mlc mcp-tester](https://github.com/hmsoft0815/mlc_mcptester) gegen die Spezifikation 2026-07-28: `task test:mcp`.

## Abhängigkeit: MLC Artifact Server

Die Artifact-Integration ist **optional**. Wird beim Start keine `-artifact-addr` angegeben, bietet der MCP-Server die Artifact-bezogenen Parameter im Tool-Schema **nicht** an, um die Schnittstelle sauber zu halten.
