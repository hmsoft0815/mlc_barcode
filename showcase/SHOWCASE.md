# MLC Barcode Showcase

This document provides examples of the various barcode types and output formats that can be generated using the `mlc_barcode` CLI tool.

## 1. QR Code

### SVG Format (Standard)
```bash
./bin/barcode -type qr -data "https://github.com/mlcmcp/mlc_barcode" -out showcase/assets/qr.svg
```
![QR Code SVG](assets/qr.svg)

### PNG Format (Standard)
```bash
./bin/barcode -type qr -data "Gemini CLI" -out showcase/assets/qr.png
```
![QR Code PNG](assets/qr.png)

### Transparent Background (SVG)
```bash
./bin/barcode -type qr -data "Transparent QR" -out showcase/assets/qr_transparent.svg -bg transparent
```
![Transparent QR SVG](assets/qr_transparent.svg)

### Red on Transparent (PNG)
```bash
./bin/barcode -type qr -data "Red on Transparent" -out showcase/assets/qr_red_transparent.png -fg red -bg transparent
```
![Red on Transparent PNG](assets/qr_red_transparent.png)

---

## 2. Code 128

### SVG with Text
```bash
./bin/barcode -type code128 -data "MLC-BARCODE-123" -out showcase/assets/code128.svg -text
```
![Code 128 SVG](assets/code128.svg)

### Colors (Blue on Yellow)
```bash
./bin/barcode -type code128 -data "Blue on Yellow" -out showcase/assets/code128_colors.svg -fg blue -bg "#ffff00" -text
```
![Blue on Yellow](assets/code128_colors.svg)

---

## 3. DataMatrix (SVG)
```bash
./bin/barcode -type datamatrix -data "DataMatrix Example" -out showcase/assets/datamatrix.svg
```
![DataMatrix SVG](assets/datamatrix.svg)

---

## 4. EAN-13 (SVG with Text)
```bash
./bin/barcode -type ean13 -data "4006381333931" -out showcase/assets/ean13.svg -text
```
![EAN-13 SVG](assets/ean13.svg)

---

## 5. Code 39 (PNG with Text)
```bash
./bin/barcode -type code39 -data "CODE39" -out showcase/assets/code39.png -text
```
![Code 39 PNG](assets/code39.png)

---

## 6. Specialized QR Codes

These examples use the new formatting templates for common tasks like WIFI access, contact sharing, and event invitations.

### WIFI Access
```bash
# Data format: WIFI:T:WPA;S:ShowcaseNet;P:password123;;
./bin/barcode -type qr -data "WIFI:T:WPA;S:ShowcaseNet;P:password123;;" -out showcase/assets/qr_wifi.png
```
![WIFI QR](assets/qr_wifi.png)

### vCard 3.0 (Contact)
```bash
# Data format: BEGIN:VCARD...
./bin/barcode -type qr -data "BEGIN:VCARD..." -out showcase/assets/qr_vcard.png
```
![vCard QR](assets/qr_vcard.png)

### iCalendar (Event)
```bash
# Data format: BEGIN:VCALENDAR...
./bin/barcode -type qr -data "BEGIN:VCALENDAR..." -out showcase/assets/qr_event.png
```
![Event QR](assets/qr_event.png)

---

## Summary of Parameters

| Parameter | Description |
|-----------|-------------|
| `-type`   | Barcode type (qr, datamatrix, code128, code39, ean13, etc.) |
| `-data`   | The data to encode |
| `-out`    | Filename (extension determines format: .svg or .png) |
| `-text`   | Display text below the barcode (if supported) |
| `-fg`     | Foreground color (e.g. black, red, #0000ff) |
| `-bg`     | Background color (e.g. white, transparent, #ffff00) |
| `-width`  | Optional width in pixels |
| `-height` | Optional height in pixels |
