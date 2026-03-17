# QR-Code Formate

QR-Codes funktionieren intern wie ein Container für Text. Je nachdem, welches Präfix (Start-Kürzel) am Anfang steht, weiß dein Smartphone, ob es eine Website öffnen, einen Kontakt speichern oder sich mit dem WLAN verbinden soll.

## 1. URL-Links (Speisekarten, Flyer)

Dies ist das einfachste Format. Es enthält lediglich die vollständige Adresse.

**Struktur:** `https://www.beispiel.de`

**Besonderheit:** Moderne Smartphones erkennen das `https://` und öffnen sofort den Browser.

## 2. vCard (Digitale Visitenkarte)

Dieses Format ist komplexer, da es viele Felder (Name, Tel, E-Mail) in einem Block speichert. Meist wird der vCard 3.0 Standard genutzt.

**Struktur:**
```text
BEGIN:VCARD
VERSION:3.0
N:Mustermann;Max
FN:Max Mustermann
ORG:Beispiel GmbH
TEL;TYPE=WORK,VOICE:0123456789
EMAIL:max@beispiel.de
ADR;TYPE=WORK:;;Musterstraße 1;Berlin;;10115;Germany
URL:https://www.beispiel.de
END:VCARD
```

**Hinweis:** Je mehr Text (z. B. lange Adressen) du hinzufügst, desto "feiner" und schwieriger zu scannen wird das Muster des QR-Codes.

## 3. WLAN-Zugangsdaten

Dieses Format wurde ursprünglich von ZXing entwickelt und ist heute Standard für Cafés und Hotels.

**Struktur:** `WIFI:T:WPA;S:NetzwerkName;P:Passwort123;;`

**Felder:**
- `T:` Verschlüsselung (WPA, WEP oder nopass)
- `S:` SSID (der Name deines WLANs)
- `P:` Passwort
- `;;` (Doppel-Semikolon schließt den Code ab)

## 4. Einladungen (vCalendar / iCalendar)

Für Einladungen mit Ort und Zeit nutzt man das vCalendar-Format. Dein Smartphone erkennt dies und schlägt vor, den Termin im Kalender zu speichern.

**Struktur:**
```text
BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
SUMMARY:Geburtstagsparty von Max
DTSTART:20240520T180000Z
DTEND:20240521T020000Z
LOCATION:Musterbar, Hauptstraße 10, Berlin
DESCRIPTION:Hier gibt es mehr Infos: https://party-website.de
END:VEVENT
END:VCALENDAR
```

**Details:**
- `DTSTART:` Startzeit im Format YYYYMMDDTHHMMSS (das Z steht für UTC-Zeit)
- `LOCATION:` Textuelle Ortsangabe
- `DESCRIPTION:` Hier kannst du zusätzliche Details oder einen Link zur Website unterbringen