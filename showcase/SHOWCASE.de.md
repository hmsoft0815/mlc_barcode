# MLC Barcode Showcase

Dieses Dokument zeigt Beispiele für die verschiedenen Barcode-Typen und Ausgabeformate, die mit dem `mlc_barcode` CLI-Tool generiert werden können.

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

### Transparenter Hintergrund (SVG)
```bash
./bin/barcode -type qr -data "Transparent QR" -out showcase/assets/qr_transparent.svg -bg transparent
```
![Transparent QR SVG](assets/qr_transparent.svg)

### Rot auf Transparent (PNG)
```bash
./bin/barcode -type qr -data "Red on Transparent" -out showcase/assets/qr_red_transparent.png -fg red -bg transparent
```
![Red on Transparent PNG](assets/qr_red_transparent.png)

---

## 2. Code 128

### SVG mit Text
```bash
./bin/barcode -type code128 -data "MLC-BARCODE-123" -out showcase/assets/code128.svg -text
```
![Code 128 SVG](assets/code128.svg)

### Farben (Blau auf Gelb)
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

## 4. EAN-13 (SVG mit Text)
```bash
./bin/barcode -type ean13 -data "4006381333931" -out showcase/assets/ean13.svg -text
```
![EAN-13 SVG](assets/ean13.svg)

---

## 5. Code 39 (PNG mit Text)
```bash
./bin/barcode -type code39 -data "CODE39" -out showcase/assets/code39.png -text
```
![Code 39 PNG](assets/code39.png)

---

## Zusammenfassung der Parameter

| Parameter | Beschreibung |
|-----------|--------------|
| `-type`   | Barcode-Typ (qr, datamatrix, code128, code39, ean13, etc.) |
| `-data`   | Die zu kodierenden Daten |
| `-out`    | Dateiname (Endung bestimmt Format: .svg oder .png) |
| `-text`   | Text unter dem Barcode anzeigen (falls unterstützt) |
| `-fg`     | Vordergrundfarbe (z.B. black, red, #0000ff) |
| `-bg`     | Hintergrundfarbe (z.B. white, transparent, #ffff00) |
| `-width`  | Optionale Breite in Pixeln |
| `-height` | Optionale Höhe in Pixeln |
