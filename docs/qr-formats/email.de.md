# E-Mail-QR-Format

Öffnet das Standard-Mailprogramm mit vorausgefülltem Empfänger, Betreff und Text.

Standard: **RFC 6068** (`mailto:`-URI).

## Parameter

| CLI-Flag | MCP-Parameter | Beschreibung |
|----------|---------------|-------------|
| `-mail-to` | `to` | **Pflicht.** Empfängeradresse. |
| `-mail-subject` | `subject` | Betreffzeile. |
| `-mail-body` | `body` | Nachrichtentext. |

## Anwendungsbeispiele

### CLI-Beispiel
```bash
./bin/barcode -mail-to "info@example.com" -mail-subject "Anfrage" \
  -mail-body "Guten Tag," -out mail.png
```

### MCP-Tool-Beispiel
```json
{
  "name": "generate_email_qr",
  "arguments": {
    "to": "info@example.com",
    "subject": "Anfrage",
    "body": "Guten Tag,"
  }
}
```

## Hinweise

Betreff und Text werden prozentkodiert, damit Umlaute, Zeilenumbrüche und kaufmännische Und unbeschadet ankommen. Ein Leerzeichen wird zu `%20`, nicht zu `+`: RFC 6068 gibt dem `+` keine Sonderbedeutung, ein regelkonformes Programm zeigte sonst ein Pluszeichen.

Ein langer Text macht den QR-Code dichter und aus der Entfernung schwerer scanbar. Ab ein paar Zeilen ist ein Link auf eine Seite der bessere Träger.
