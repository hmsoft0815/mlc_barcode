# SMS QR Format

Opens a new SMS with recipient and text pre-filled.

Format: `smsto:<number>:<text>`, the spelling Android and iOS both understand.

## Parameters

| CLI Flag | MCP Parameter | Description |
|----------|---------------|-------------|
| `-sms` | `phone_number` | **Required.** Recipient number. |
| `-sms-text` | `message` | Message text. |

## Usage Examples

### CLI Example
```bash
./bin/barcode -sms "+493012345678" -sms-text "On my way" -out sms.png
```

### MCP Tool Example
```json
{
  "name": "generate_sms_qr",
  "arguments": {
    "phone_number": "+493012345678",
    "message": "On my way"
  }
}
```

## Notes

The message is inserted verbatim — nothing is escaped. A colon inside the text therefore produces a third colon in the payload; readers split on the first two and are unbothered. Pinned as behaviour in `formats_test.go`, because changing it would mean inventing an escaping scheme no reader has agreed to.

Nothing is sent by scanning. The message is composed and waits for the user.
