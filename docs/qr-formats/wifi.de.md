# WLAN QR-Code Format

Dieses Dokument beschreibt, wie QR-Codes für einen einfachen WLAN-Zugriff generiert werden.

## Parameter

| CLI Flag | MCP Parameter | Beschreibung |
|----------|---------------|-------------|
| `-wifi-ssid` | `ssid` | **Erforderlich.** Der Name des WLAN-Netzwerks. |
| `-wifi-pass` | `password` | Das WLAN-Passwort. |
| `-wifi-enc` | `encryption` | Verschlüsselungstyp: `WPA`, `WEP`, oder `nopass` (Standard: `WPA`). |
| (N/A) | `hidden` | Auf `true` setzen, wenn die SSID versteckt ist. |

## Nutzungsbeispiele

### CLI Beispiel
```bash
./bin/barcode -wifi-ssid "HomeOffice" -wifi-pass "geheim123" -out wlan.png
```

### MCP Tool Beispiel
```json
{
  "name": "generate_wifi_qr",
  "arguments": {
    "ssid": "CafeGast",
    "encryption": "nopass",
    "bg_color": "transparent"
  }
}
```

## Technische Details
Der generierte String folgt dem De-facto-Standard:
`WIFI:T:WPA;S:SSID;P:PASSWORD;;`
Sonderzeichen wie `:`, `;`, und `\` werden automatisch maskiert.
