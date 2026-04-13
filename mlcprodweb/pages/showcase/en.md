# Showcase

Experience the variety of barcode types and styling options provided by MLC Barcode.

## Real World Examples

### Website Access
A clean QR code for our main website.
```bash
barcode -type qr -data "https://mlcgo.eu" -out mlcgo_qr.png
```
![mlcgo.eu](../../assets/mlcgo_qr.png)

### Retail & Products (EAN-13)
Standard European Article Number for retail products, including the human-readable text below.
```bash
barcode -type ean13 -data "4006381333931" -text -out retail.png
```
![EAN-13](../../assets/barcode_ean13.png)

### Logistics & Shipping (Code 128)
High-density barcode for tracking numbers and serial numbers.
```bash
barcode -type code128 -data "LOGISTICS-12345678" -text -out log.png
```
![Code 128](../../assets/barcode_code128.png)

## Specialized Templates

MLC Barcode simplifies complex data formatting for common use cases.

### WiFi Connection
Generate a QR code that allows mobile devices to join a network instantly.
```bash
barcode -wifi-ssid "MyNet" -wifi-pass "secret" -out wifi.png
```
![WiFi QR](../../assets/qr_wifi.png)

### Business Cards (vCard)
Share contact information in a standardized format.
```bash
barcode -vcard-first "Michael" -vcard-last "Lechner" -vcard-email "mlcgo.eu@michael-lchner.de"
```
![vCard QR](../../assets/qr_vcard.png)

### Calendar Events
Invite guests to events with iCalendar compatibility.
![Event QR](../../assets/qr_event.png)
