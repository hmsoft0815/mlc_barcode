---
name: contact-card
description: Turn contact details — an e-mail signature, a business card, a note — into a vCard QR code that saves the contact in one scan. Use when someone wants a QR code for a business card, a contact, a name badge or "save my number".
---

# Contact card (vCard) as a QR code

Tool: `generate_vcard_qr` (vCard 3.0).

## Gather

| Field | Argument | Notes |
|---|---|---|
| Name | `first_name`, `last_name` | both required; split "Dr. Anna Meier" as first `Anna`, last `Meier`, put titles into `title` |
| Company, role | `org`, `title` | |
| Phone | `phone` | international form `+49 30 1234567` works in every country |
| E-mail | `email` | |
| Address | `address`, `zip`, `city`, `country` | street and number in `address` |
| Website | `url` | with `https://` |

## Rules

- Take values from what the user gave; do not invent a company, title or website.
- Only fields that are known — every extra field makes the code denser and harder to scan on a small badge.
- For signatures with several numbers, ask which one belongs on the card (usually the mobile or direct line).
- Show the resulting `encoded_data` so the user can check spelling before printing.

## Example

```json
{"first_name": "Anna", "last_name": "Meier", "org": "Praxis Dr. Meier", "title": "Ärztin",
 "phone": "+49 30 1234567", "email": "anna.meier@example.org", "url": "https://example.org",
 "format": "png", "width": 500}
```
