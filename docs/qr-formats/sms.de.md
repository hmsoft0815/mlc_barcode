# SMS-QR-Format

Öffnet eine neue SMS mit vorausgefülltem Empfänger und Text.

Format: `smsto:<Nummer>:<Text>` — die Schreibweise, die Android und iOS beide verstehen.

## Parameter

| CLI-Flag | MCP-Parameter | Beschreibung |
|----------|---------------|-------------|
| `-sms` | `phone_number` | **Pflicht.** Empfängernummer. |
| `-sms-text` | `message` | Nachrichtentext. |

## Anwendungsbeispiele

### CLI-Beispiel
```bash
./bin/barcode -sms "+493012345678" -sms-text "Bin da" -out sms.png
```

### MCP-Tool-Beispiel
```json
{
  "name": "generate_sms_qr",
  "arguments": {
    "phone_number": "+493012345678",
    "message": "Bin da"
  }
}
```

## Hinweise

Der Text wird unverändert eingesetzt, es wird nichts maskiert. Ein Doppelpunkt im Text erzeugt daher einen dritten Doppelpunkt in den Nutzdaten; Lesegeräte trennen an den ersten beiden und stören sich nicht daran. In `formats_test.go` als Verhalten festgeschrieben — es zu ändern hiesse, eine Maskierung zu erfinden, auf die sich kein Lesegerät geeinigt hat.

Durch das Scannen wird nichts versendet. Die Nachricht wird verfasst und wartet auf den Nutzer.
