---
name: wifi-guest-access
description: Create a Wi-Fi QR code so guests join a network by scanning instead of typing the password — for a guest network, café, office or holiday flat. Use when someone wants to share Wi-Fi access, a network password or a WLAN code.
---

# Wi-Fi access as a QR code

Tool: `generate_wifi_qr`.

## Gather

| Field | Argument | Notes |
|---|---|---|
| Network name | `ssid` | required, exactly as the router shows it (case and spaces matter) |
| Password | `password` | exactly as typed on the router |
| Security | `encryption` | `WPA` for WPA, WPA2 and WPA3 (almost every network today); `WEP` only for very old routers; `nopass` for open networks (then no password) |
| Hidden network | `hidden` | `true` only if the network does not broadcast its name |

## Rules

- Copy SSID and password character for character. Special characters (`; , : " \`) are escaped by the tool — pass them as they are, do not escape them yourself.
- Do not "correct" a password (letters that look alike: `l`/`1`, `O`/`0`). If the user typed it from a photo or sticker and a character is unclear, ask.
- Recommend a **guest network** for codes that hang on a wall — anyone who photographs the code has the password.
- For printing: `format: png`, `width` 600 or more; a caption with the network name (`caption`) helps people pick the right code.

## Example

```json
{"ssid": "Cafe Sonne Gast", "password": "kaffee&kuchen2026", "encryption": "WPA",
 "format": "png", "width": 600, "caption": "WLAN: Cafe Sonne Gast"}
```
