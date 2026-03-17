# Event (iCalendar) QR Code Format

This document describes how to generate QR codes for calendar events (RFC 5545).

## Standards
This implementation follows the **iCalendar** standard, as defined in [RFC 5545](https://datatracker.ietf.org/doc/html/rfc5545).

## Parameters

| CLI Flag | MCP Parameter | Description |
|----------|---------------|-------------|
| `-event-summary` | `summary` | **Required.** Title of the event. |
| `-event-start` | `start_time` | **Required.** Start time in `YYYYMMDDTHHMMSS` format. |
| `-event-end` | `end_time` | End time in `YYYYMMDDTHHMMSS` format. |
| (N/A) | `description` | Detailed description of the event. |
| (N/A) | `location` | Physical or virtual location. |
| `-event-tz` | `timezone` | TimeZone identifier (e.g., `Europe/Berlin`). |
| (N/A) | `latitude` | Geographical latitude (decimal). |
| (N/A) | `longitude` | Geographical longitude (decimal). |

## Usage Examples

### CLI Example
```bash
./bin/barcode -event-summary "Project Kickoff" -event-start "20260401T100000" -event-tz "Europe/Berlin" -out kickoff.svg
```

### MCP Tool Example
```json
{
  "name": "generate_event_qr",
  "arguments": {
    "summary": "Team Lunch",
    "start_time": "20260320T120000",
    "location": "Central Park",
    "latitude": 40.785091,
    "longitude": -73.968285
  }
}
```

## Technical Details
The output is a structured string starting with `BEGIN:VCALENDAR` and containing a `BEGIN:VEVENT` block. It supports `DTSTART`, `DTEND`, `SUMMARY`, `LOCATION`, `DESCRIPTION`, and `GEO` fields. Time zones are handled via the `;TZID=` prefix unless UTC (`Z`) is specified.
