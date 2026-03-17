# Projektplan: QR-Code Format-Erweiterungen

Dieser Plan beschreibt die Integration von spezialisierten QR-Code-Formaten (vCard, WIFI, vCalendar) basierend auf der Entwurfsdatei `qrcode-enhancements.md`.

## Zielsetzung
Das CLI-Tool und der MCP-Server sollen Vorlagen erhalten, mit denen komplexe QR-Code-Inhalte einfach generiert werden können, ohne dass der Benutzer die genaue Syntax (z.B. vCard-Header) kennen muss.

## Meilensteine

### 1. Internes Format-Paket (`internal/qrformats`)
Entwicklung einer Bibliothek zur Validierung und Formatierung der spezifischen Inhalts-Strings.
- [ ] **WIFI**: Generator für `WIFI:T:WPA;S:SSID;P:PASSWORD;;`
- [ ] **vCard**: Generator für Kontaktinformationen (vCard 3.0 / [RFC 6350](https://www.rfc-editor.org/rfc/rfc6350.html)).
- [ ] **vCalendar**: Generator für Event-Einladungen (iCalendar 2.0 / [RFC 5545](https://datatracker.ietf.org/doc/html/rfc5545)).

### 2. CLI-Erweiterung (`cmd/barcode`)
Hinzufügen von Unterbefehlen oder speziellen Flags zur Nutzung der Formate.
- [ ] Implementierung eines interaktiven Modus oder spezifischer Flags (z.B. `-wifi-ssid`, `-vcard-name`).
- [ ] Integration in die bestehende Logik.

### 3. MCP-Server Erweiterung (`cmd/mcp-server`)
Erweiterung der MCP-Schnittstelle, um LLMs die gezielte Erstellung dieser Formate zu ermöglichen.
- [ ] Neues Tool: `generate_wifi_qr`
- [ ] Neues Tool: `generate_vcard_qr`
- [ ] Neues Tool: `generate_event_qr`

### 4. Dokumentation & Showcase
- [ ] Aktualisierung der READMEs.
- [ ] Aufnahme von Beispielen in den `showcase/`.

## Zeitplan
1. **Phase 1**: Kern-Logik in `internal/qrformats`.
2. **Phase 2**: Integration in den MCP-Server (da dies den höchsten Nutzwert für LLMs hat).
3. **Phase 3**: CLI-Anpassung.
4. **Phase 4**: Tests und Dokumentation.
