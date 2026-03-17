# Termin (iCalendar) QR-Code Format

Dieses Dokument beschreibt, wie QR-Codes für Kalendertermine (RFC 5545) generiert werden.

## Standards
Diese Implementierung folgt dem **iCalendar** Standard, wie in [RFC 5545](https://datatracker.ietf.org/doc/html/rfc5545) definiert.

## Parameter

| CLI Flag | MCP Parameter | Beschreibung |
|----------|---------------|-------------|
| `-event-summary` | `summary` | **Erforderlich.** Titel des Termins. |
| `-event-start` | `start_time` | **Erforderlich.** Startzeit im Format `YYYYMMDDTHHMMSS`. |
| `-event-end` | `end_time` | Endzeit im Format `YYYYMMDDTHHMMSS`. |
| (N/A) | `description` | Detaillierte Beschreibung des Termins. |
| (N/A) | `location` | Physischer oder virtueller Ort. |
| `-event-tz` | `timezone` | Zeitzonen-Kennung (z. B. `Europe/Berlin`). |
| (N/A) | `latitude` | Geografische Breite (Dezimalwert). |
| (N/A) | `longitude` | Geografische Länge (Dezimalwert). |

## Nutzungsbeispiele

### CLI Beispiel
```bash
./bin/barcode -event-summary "Projekt Kickoff" -event-start "20260401T100000" -event-tz "Europe/Berlin" -out kickoff.svg
```

### MCP Tool Beispiel
```json
{
  "name": "generate_event_qr",
  "arguments": {
    "summary": "Team Mittagessen",
    "start_time": "20260320T120000",
    "location": "Stadtpark",
    "latitude": 48.137154,
    "longitude": 11.576124
  }
}
```

## Technische Details
Die Ausgabe ist ein strukturierter String, der mit `BEGIN:VCALENDAR` beginnt und einen `BEGIN:VEVENT` Block enthält. Es werden die Felder `DTSTART`, `DTEND`, `SUMMARY`, `LOCATION`, `DESCRIPTION`, und `GEO` unterstützt. Zeitzonen werden über das `;TZID=` Präfix gehandhabt, sofern nicht UTC (`Z`) angegeben ist.
