# Telephone Number QR Format

Opens the dialler with the number filled in. The user still has to press call.

Standard: **RFC 3966** (`tel:` URI).

## Parameters

| CLI Flag | MCP Parameter | Description |
|----------|---------------|-------------|
| `-tel` | `phone_number` | **Required.** The number to dial. |

## Usage Examples

### CLI Example
```bash
./bin/barcode -tel "+49 30 123456" -out call.png
```

### MCP Tool Example
```json
{
  "name": "generate_tel_qr",
  "arguments": {
    "phone_number": "+493012345678"
  }
}
```

## Notes

Spaces and formatting characters are stripped, so `+49 30 123456` and `+493012345678` produce the same payload.

Use the international form with a leading `+`. A national number works only for someone dialling from that country.
