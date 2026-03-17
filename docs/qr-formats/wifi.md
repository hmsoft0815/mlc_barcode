# WIFI QR Code Format

This document describes how to generate QR codes for easy WIFI network access.

## Parameters

| CLI Flag | MCP Parameter | Description |
|----------|---------------|-------------|
| `-wifi-ssid` | `ssid` | **Required.** The name of the WIFI network. |
| `-wifi-pass` | `password` | The WIFI password. |
| `-wifi-enc` | `encryption` | Encryption type: `WPA`, `WEP`, or `nopass` (Default: `WPA`). |
| (N/A) | `hidden` | Set to `true` if the SSID is hidden. |

## Usage Examples

### CLI Example
```bash
./bin/barcode -wifi-ssid "HomeOffice" -wifi-pass "secret123" -out wifi.png
```

### MCP Tool Example
```json
{
  "name": "generate_wifi_qr",
  "arguments": {
    "ssid": "CafeGuest",
    "encryption": "nopass",
    "bg_color": "transparent"
  }
}
```

## Technical Details
The generated string follows the de-facto standard:
`WIFI:T:WPA;S:SSID;P:PASSWORD;;`
Special characters like `:`, `;`, and `\` are automatically escaped.
