# Email QR Format

Opens the default mail client with recipient, subject and body pre-filled.

Standard: **RFC 6068** (`mailto:` URI).

## Parameters

| CLI Flag | MCP Parameter | Description |
|----------|---------------|-------------|
| `-mail-to` | `to` | **Required.** Recipient address. |
| `-mail-subject` | `subject` | Subject line. |
| `-mail-body` | `body` | Message body. |

## Usage Examples

### CLI Example
```bash
./bin/barcode -mail-to "info@example.com" -mail-subject "Anfrage" \
  -mail-body "Guten Tag," -out mail.png
```

### MCP Tool Example
```json
{
  "name": "generate_email_qr",
  "arguments": {
    "to": "info@example.com",
    "subject": "Anfrage",
    "body": "Guten Tag,"
  }
}
```

## Notes

Subject and body are percent-encoded, so umlauts, line breaks and ampersands survive intact. A space becomes `%20`, not `+`: RFC 6068 gives `+` no special meaning, and a client that follows it would show a literal plus.

A long body makes the QR code denser and harder to scan from a distance. For anything beyond a few lines, a link to a page is the better carrier.
