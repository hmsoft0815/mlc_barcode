---
name: event-invite
description: Turn an invitation, appointment or meeting into a calendar QR code that phones add to their calendar in one scan — with time zone, place and organizer. Use when someone wants a QR code for an event, appointment, invitation, meeting or deadline.
---

# Event invitation as a QR code

Tool: `generate_event_qr` (iCalendar, RFC 5545).

## Gather

| Field | Argument | Notes |
|---|---|---|
| Title | `summary` | required, short: "Sprechstunde Dr. Meier", "Team lunch" |
| Start | `start_time` | required, `YYYYMMDDTHHMMSS` local time, e.g. `20261015T140000` |
| End | `end_time` | same format; if unknown use start + 1 hour |
| Time zone | `timezone` | IANA name, e.g. `Europe/Berlin` — always set it for timed events, otherwise phones treat the time as floating |
| Place | `location` | address or room; add `latitude`/`longitude` when known so maps can route |
| Organizer | `organizer`, `organizer_email` | shown as "Last name, First name"; the e-mail is optional |
| Details | `description` | agenda, dial-in link, what to bring |

## Rules

- **All-day events:** `start_time` = `YYYYMMDD`; `end_time` = the **day after** the last day (a one-day event on 15 Oct: `20261015` → `20261016`). No time zone.
- Ask when date or time is ambiguous ("next Friday" across a month end, 12-hour times without am/pm). Do not guess the year — use the coming occurrence.
- Keep the text short: a QR code with a long description gets dense and hard to scan from paper; put long agendas behind a link.
- Show the user the resulting `encoded_data` (the iCalendar text) so they can check date and time before printing.

## Example

```json
{"summary": "Sprechstunde", "start_time": "20261015T140000", "end_time": "20261015T143000",
 "timezone": "Europe/Berlin", "location": "Praxis Dr. Meier, Hauptstr. 1, Berlin",
 "organizer": "Meier, Anna", "organizer_email": "praxis@example.org", "format": "png", "width": 600}
```
