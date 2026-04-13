# Showcase

Erleben Sie die Vielfalt der Barcode-Typen und Styling-Optionen von MLC Barcode.

## Praxisbeispiele

### Webseitenzugriff
Ein klarer QR-Code für unsere Hauptwebseite.
```bash
barcode -type qr -data "https://mlcgo.eu" -out mlcgo_qr.png
```
![mlcgo.eu](../../assets/mlcgo_qr.png)

### Handel & Produkte (EAN-13)
Standardmäßig verwendete European Article Number für Handelsprodukte, inklusive des lesbaren Textes unter dem Barcode.
```bash
barcode -type ean13 -data "4006381333931" -text -out retail.png
```
![EAN-13](../../assets/barcode_ean13.png)

### Logistik & Versand (Code 128)
Ein Barcode mit hoher Dichte für Sendungsnummern und Seriennummern.
```bash
barcode -type code128 -data "LOGISTICS-12345678" -text -out log.png
```
![Code 128](../../assets/barcode_code128.png)

## Spezialisierte Vorlagen

MLC Barcode vereinfacht komplexe Datenformatierungen für gängige Anwendungsfälle.

### WLAN-Verbindung
Erzeugen Sie einen QR-Code, mit dem mobile Geräte sofort einem Netzwerk beitreten können.
```bash
barcode -wifi-ssid "MeinNetz" -wifi-pass "geheim" -out wifi.png
```
![WLAN-QR](../../assets/qr_wifi.png)

### Visitenkarten (vCard)
Teilen Sie Kontaktinformationen in einem standardisierten Format.
```bash
barcode -vcard-first "Michael" -vcard-last "Lechner" -vcard-email "mlcgo.eu@michael-lchner.de"
```
![vCard-QR](../../assets/qr_vcard.png)

### Kalenderereignisse
Laden Sie Gäste zu Veranstaltungen mit iCalendar-Kompatibilität ein.
![Event-QR](../../assets/qr_event.png)
