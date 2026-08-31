# EPC QR Code (GiroCode) Format

This document describes how to generate a GiroCode — the QR code European banking apps scan to pre-fill a SEPA credit transfer.

Standard: European Payments Council **EPC069-12**, payload version `002`.

## Parameters

| CLI Flag | MCP Parameter | Description |
|----------|---------------|-------------|
| `-epc-name` | `name` | **Required.** Beneficiary name (Empfänger). |
| `-epc-iban` | `iban` | **Required.** Beneficiary IBAN. Spaces are stripped automatically. |
| `-epc-bic` | `bic` | BIC. Optional in payload version 002. |
| `-epc-amount` | `amount` | Amount in EUR, e.g. `12.50`. Omit for an open amount. |
| `-epc-ref` | `reference` | Remittance information (Verwendungszweck), max 140 characters. |

## Usage Examples

### CLI Example
```bash
./bin/barcode -epc-name "Example Charity" -epc-iban "DE89 3704 0044 0532 0130 00" \
  -epc-amount 25.00 -epc-ref "Membership 2026" -out donation.png
```

### MCP Tool Example
```json
{
  "name": "generate_epc_qr",
  "arguments": {
    "name": "Example Charity",
    "iban": "DE89370400440532013000",
    "amount": 25.00,
    "reference": "Membership 2026"
  }
}
```

## Notes

**The payload is positional.** A reader takes line 6 as the beneficiary and line 7 as the IBAN, whatever they contain — there are no field names to check against. Swapping two lines produces a perfectly scannable code that pays the wrong thing, which is why `formats_test.go` compares the payload line by line rather than by substring.

**An amount of zero is not `EUR0.00`.** It leaves the amount element empty, and the payer fills it in. Writing a zero would produce a transfer the bank rejects.

**Trailing empty elements are omitted**, which the standard allows. Elements in the middle are always written, even when empty, because their position is what identifies them.

The IBAN is not checked for a valid checksum. A typo produces a code the banking app refuses — it does not produce a transfer to somebody else.
