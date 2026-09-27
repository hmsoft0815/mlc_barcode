---
name: label-sheet
description: Turn a list of products, parts or items into barcodes ready for labels — the right symbology per use and a readable caption. Use when someone wants labels, shelf tags, inventory codes, asset tags or a batch of barcodes from a list or spreadsheet.
---

# Barcodes for a label sheet

Tools: `generate_barcode` once per item; the prompt `product_labels` walks through a list `code;label` line by line.

## Choose the symbology

| Content | `type` | Why |
|---|---|---|
| Retail product with a GTIN/EAN (13 digits) | `ean13` | what tills scan; the tool computes or checks the check digit |
| Small retail pack (8 digits) | `ean8` | |
| Internal article number, SKU, letters and digits | `code128` | compact, any ASCII |
| Only capitals, digits, `-.$/+%`, older scanners | `code39` | |
| Numbers with an even count of digits (cartons, logistics) | `itf` | |
| A link or longer text, scanned by phones | `qr` | |
| Very small labels, many characters | `datamatrix` | smallest per character |

## Rules

- One call per item; keep the order of the list so labels match the sheet.
- Put the human-readable text in `caption` (product name, price, location); `text: true` shows the encoded value instead.
- For printing use `format: png` with an explicit `width`/`height` matching the label (e.g. 600×300 for a 1D code, 400×400 for QR/DataMatrix), or `svg` for sharp output at any size. Keep `quiet_zone` on — scanners need the margin.
- An EAN with a wrong check digit is rejected and the error names the right one: report that to the user instead of silently fixing their data.
- Do not mix symbologies within one sheet unless the list asks for it.

## Example

```json
{"type": "ean13", "data": "4006381333931", "caption": "Schraube M4 · 0,12 €", "format": "png", "width": 600, "height": 300}
```
