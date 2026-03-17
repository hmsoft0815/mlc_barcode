# vCard QR Code Format

This document describes how to generate QR codes for digital business cards (vCard 3.0).

## Standards
This implementation follows the **vCard 3.0** standard, referencing [RFC 6350](https://www.rfc-editor.org/rfc/rfc6350.html) for general vCard structures.

## Parameters

| CLI Flag | MCP Parameter | Description |
|----------|---------------|-------------|
| `-vcard-first` | `first_name` | **Required.** First name of the contact. |
| `-vcard-last` | `last_name` | **Required.** Last name of the contact. |
| (N/A) | `org` | Organization or company name. |
| (N/A) | `title` | Job title. |
| `-vcard-tel` | `phone` | Phone number. |
| `-vcard-email` | `email` | Email address. |
| (N/A) | `address` | Street address. |
| (N/A) | `city` | City. |
| (N/A) | `zip` | ZIP / Postal code. |
| (N/A) | `country` | Country. |
| (N/A) | `url` | Website URL. |

## Usage Examples

### CLI Example
```bash
./bin/barcode -vcard-first "John" -vcard-last "Doe" -vcard-email "john.doe@example.com" -vcard-tel "+123456789" -out contact.svg
```

### MCP Tool Example
```json
{
  "name": "generate_vcard_qr",
  "arguments": {
    "first_name": "Jane",
    "last_name": "Smith",
    "org": "Tech Corp",
    "email": "jane@tech.com",
    "url": "https://tech.com"
  }
}
```

## Technical Details
The output is a structured string starting with `BEGIN:VCARD` and ending with `END:VCARD`. It includes properly formatted `N`, `FN`, `ORG`, `TEL`, and `EMAIL` fields.
