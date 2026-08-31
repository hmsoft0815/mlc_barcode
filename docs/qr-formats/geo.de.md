# Geo-Koordinaten-QR-Format

Öffnet Koordinaten in der Kartenanwendung des Telefons — Google Maps, Apple Maps oder OpenStreetMap, je nach Einstellung des Geräts.

Standard: **RFC 5870** (`geo:`-URI).

## Parameter

| CLI-Flag | MCP-Parameter | Beschreibung |
|----------|---------------|-------------|
| `-geo-lat` | `latitude` | **Pflicht.** Breitengrad, z. B. `52.5200`. |
| `-geo-lon` | `longitude` | **Pflicht.** Längengrad, z. B. `13.4050`. |
| `-geo-query` | `query` | Optionaler Ortsname, angehängt als `?q=`. |

## Anwendungsbeispiele

### CLI-Beispiel
```bash
./bin/barcode -geo-lat 52.5200 -geo-lon 13.4050 -geo-query "Brandenburger Tor" -out ort.png
```

### MCP-Tool-Beispiel
```json
{
  "name": "generate_geo_qr",
  "arguments": {
    "latitude": 52.5200,
    "longitude": 13.4050,
    "query": "Brandenburger Tor"
  }
}
```

## Hinweise

Koordinaten sind Dezimalgrad, mit Punkt als Dezimaltrenner und führendem Minus für Süd und West. Grad-Minuten-Sekunden wird nicht angenommen.

Ist `query` gesetzt, zeigen die meisten Anwendungen den Namen an und zentrieren darauf; was den Punkt tatsächlich bestimmt, bleiben die Koordinaten.
