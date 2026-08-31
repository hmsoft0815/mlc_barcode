# EPC-QR-Code (GiroCode) Format

Dieses Dokument beschreibt den GiroCode — den QR-Code, den europäische Banking-Apps scannen, um eine SEPA-Überweisung vorauszufüllen.

Standard: European Payments Council **EPC069-12**, Nutzdaten-Version `002`.

## Parameter

| CLI-Flag | MCP-Parameter | Beschreibung |
|----------|---------------|-------------|
| `-epc-name` | `name` | **Pflicht.** Empfängername. |
| `-epc-iban` | `iban` | **Pflicht.** IBAN des Empfängers. Leerzeichen werden automatisch entfernt. |
| `-epc-bic` | `bic` | BIC. In Version 002 optional. |
| `-epc-amount` | `amount` | Betrag in EUR, z. B. `12.50`. Weglassen für einen offenen Betrag. |
| `-epc-ref` | `reference` | Verwendungszweck, maximal 140 Zeichen. |

## Anwendungsbeispiele

### CLI-Beispiel
```bash
./bin/barcode -epc-name "Verein e.V." -epc-iban "DE89 3704 0044 0532 0130 00" \
  -epc-amount 25.00 -epc-ref "Mitgliedsbeitrag 2026" -out spende.png
```

### MCP-Tool-Beispiel
```json
{
  "name": "generate_epc_qr",
  "arguments": {
    "name": "Verein e.V.",
    "iban": "DE89370400440532013000",
    "amount": 25.00,
    "reference": "Mitgliedsbeitrag 2026"
  }
}
```

## Hinweise

**Die Nutzdaten sind positionsbasiert.** Ein Lesegerät nimmt Zeile 6 als Empfänger und Zeile 7 als IBAN, was auch immer dort steht — es gibt keine Feldnamen zum Abgleichen. Zwei vertauschte Zeilen ergeben einen einwandfrei scanbaren Code, der das Falsche bezahlt. Deshalb vergleicht `formats_test.go` die Nutzdaten zeilenweise und nicht per Teilstring.

**Ein Betrag von null ist nicht `EUR0.00`.** Das Betragsfeld bleibt leer, und der Zahlende trägt den Betrag selbst ein. Eine geschriebene Null ergäbe eine Überweisung, die die Bank ablehnt.

**Leere Elemente am Ende entfallen**, was der Standard erlaubt. Elemente in der Mitte werden immer geschrieben, auch leer — ihre Position ist ihre Kennung.

Die IBAN wird nicht auf ihre Prüfsumme geprüft. Ein Tippfehler ergibt einen Code, den die Banking-App zurückweist — keine Überweisung an jemand anderen.
