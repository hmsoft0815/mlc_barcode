---
name: girocode-from-invoice
description: Create a GiroCode (EPC QR code for SEPA credit transfers) from an invoice so it can be paid by scanning it in a banking app. Use when someone wants a payment QR code, a GiroCode, an EPC code or a "scan to pay" code for a bill or invoice in euro.
---

# GiroCode from an invoice

Tool: `generate_epc_qr` (EPC069-12, SEPA credit transfer, euro only). The prompt `payment_qr_from_invoice` does the same step by step.

## Gather from the invoice

| Field | Argument | Notes |
|---|---|---|
| Beneficiary | `name` | the payee as on the invoice, at most 70 characters |
| IBAN | `iban` | required; spaces are fine |
| BIC | `bic` | optional within the SEPA area |
| Amount | `amount` | in euro, a number: `49.90` (not "49,90 €") |
| Reference | `reference` | invoice number and customer number — what the payee needs to match the payment, at most 140 characters |

## Rules

- **Never guess** IBAN or amount. If the invoice shows several accounts or amounts (net/gross, partial payments, discount), ask which one.
- Amount with a dot as decimal separator; no currency sign. Omit `amount` only if the user wants the payer to fill it in.
- **Check the IBAN:** `generate_epc_qr` does not verify it. Read the generated code back with `decode_barcode` (the PNG, `format: png`): its `content` shows `iban_valid` (mod-97 check) and the fields as a bank app will see them. If `iban_valid` is `false`, the IBAN is mistyped — tell the user and do not hand out the code; never change digits yourself.

## Example

```json
{"name": "Muster GmbH", "iban": "DE89 3704 0044 0532 0130 00", "amount": 49.90,
 "reference": "Rechnung 1234, Kunde 5678", "format": "png", "width": 500}
```
