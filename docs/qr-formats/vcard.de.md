# vCard QR-Code Format

Dieses Dokument beschreibt, wie QR-Codes für digitale Visitenkarten (vCard 3.0) generiert werden.

## Standards
Diese Implementierung folgt dem **vCard 3.0** Standard und referenziert [RFC 6350](https://www.rfc-editor.org/rfc/rfc6350.html) für allgemeine vCard-Strukturen.

## Parameter

| CLI Flag | MCP Parameter | Beschreibung |
|----------|---------------|-------------|
| `-vcard-first` | `first_name` | **Erforderlich.** Vorname des Kontakts. |
| `-vcard-last` | `last_name` | **Erforderlich.** Nachname des Kontakts. |
| (N/A) | `org` | Organisation oder Firmenname. |
| (N/A) | `title` | Berufsbezeichnung. |
| `-vcard-tel` | `phone` | Telefonnummer. |
| `-vcard-email` | `email` | E-Mail-Adresse. |
| (N/A) | `address` | Straße und Hausnummer. |
| (N/A) | `city` | Stadt. |
| (N/A) | `zip` | Postleitzahl. |
| (N/A) | `country` | Land. |
| (N/A) | `url` | Webseite. |

## Nutzungsbeispiele

### CLI Beispiel
```bash
./bin/barcode -vcard-first "Max" -vcard-last "Mustermann" -vcard-email "max.mustermann@example.com" -vcard-tel "+49123456789" -out kontakt.svg
```

### MCP Tool Beispiel
```json
{
  "name": "generate_vcard_qr",
  "arguments": {
    "first_name": "Erika",
    "last_name": "Musterfrau",
    "org": "Tech GmbH",
    "email": "erika@tech.de",
    "url": "https://tech.de"
  }
}
```

## Technische Details
Die Ausgabe ist ein strukturierter String, der mit `BEGIN:VCARD` beginnt und mit `END:VCARD` endet. Er enthält korrekt formatierte `N`, `FN`, `ORG`, `TEL` und `EMAIL` Felder.
