# Telefonnummer-QR-Format

Öffnet die Wählanwendung mit eingetragener Nummer. Wählen muss der Nutzer selbst.

Standard: **RFC 3966** (`tel:`-URI).

## Parameter

| CLI-Flag | MCP-Parameter | Beschreibung |
|----------|---------------|-------------|
| `-tel` | `phone_number` | **Pflicht.** Die zu wählende Nummer. |

## Anwendungsbeispiele

### CLI-Beispiel
```bash
./bin/barcode -tel "+49 30 123456" -out anruf.png
```

### MCP-Tool-Beispiel
```json
{
  "name": "generate_tel_qr",
  "arguments": {
    "phone_number": "+493012345678"
  }
}
```

## Hinweise

Leerzeichen und Formatierungszeichen werden entfernt, `+49 30 123456` und `+493012345678` ergeben also dieselben Nutzdaten.

Die internationale Schreibweise mit führendem `+` ist die richtige. Eine nationale Nummer funktioniert nur für Anrufende aus demselben Land.
