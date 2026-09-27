---
name: medicine-pack-check
description: Read the code on a medicine pack from a photo and report what it is — PZN, batch, expiry date, serial number — and whether it has expired. Use when someone sends a photo of a medicine box or bottle, asks whether a medicine is still usable, or wants to check a home pharmacy pack by pack.
---

# Check a medicine pack

Tool: `decode_barcode` with `path` (stdio only) or `image_base64` of the photo.

## Read the result

A pack code is reported with `content.kind = "pharma"`:

| Field | Meaning |
|---|---|
| `format` | `gs1` or `ifa` — the securPharm DataMatrix (prescription medicines since 2019); `pzn` — the older Code 39 barcode `-12345678` |
| `pzn`, `pzn_valid` | German pharmacy product number and whether its check digit is right |
| `expiry` | `YYMMDD`; day `00` means **the end of that month** (`280300` → usable until 31 March 2028) |
| `batch`, `serial` | batch (Charge) and the pack's serial number |
| `gtin` / `ppn` | product number with its own check (`gtin_valid`, `ppn_valid`) |

## Rules

- **Expired?** Compare the expiry with today's date: expired if the last usable day lies before today. Say it plainly and suggest returning the pack to a pharmacy.
- A **`pzn` barcode has no expiry and no batch** — the date is only printed as text on the pack ("verw. bis 03/2027", "EXP 2027-03"). Ask the user to read it or send a sharper photo of that text; do not guess.
- A check digit marked invalid (`pzn_valid: false`) means a misread or a damaged code — ask for a new photo rather than reporting the number.
- Several codes in one photo (several packs): report each with its PZN and expiry.
- Nothing found: ask for a photo with the code sharp and filling more of the picture; on a round bottle, turn it so the code faces the camera and avoid glare.
- This checks the code only. It does not identify the medicine by name and is no medical advice; for what a PZN is, the user can look it up in the BfArM database (AMIce, German only).
