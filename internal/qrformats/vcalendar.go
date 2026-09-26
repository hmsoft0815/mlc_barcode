package qrformats

import (
	"fmt"
	"strings"
)

// VCalendarOptions holds configuration for a vCalendar event
type VCalendarOptions struct {
	Summary     string
	Description string
	Location    string
	StartTime   string // YYYYMMDDTHHMMSS (local), YYYYMMDDTHHMMSSZ (UTC) or YYYYMMDD (all day)
	EndTime     string // same formats; for all-day events the day after the last day
	TimeZone    string // e.g. Europe/Berlin
	Latitude    float64
	Longitude   float64
}

// FormatVCalendar returns a formatted iCalendar 2.0 string (RFC 5545)
func FormatVCalendar(opts VCalendarOptions) string {
	var sb strings.Builder
	sb.WriteString("BEGIN:VCALENDAR\n")
	sb.WriteString("VERSION:2.0\n")
	sb.WriteString("BEGIN:VEVENT\n")

	if opts.Summary != "" {
		fmt.Fprintf(&sb, "SUMMARY:%s\n", escapeText(opts.Summary))
	}

	writeDateTime(&sb, "DTSTART", opts.StartTime, opts.TimeZone)
	writeDateTime(&sb, "DTEND", opts.EndTime, opts.TimeZone)

	if opts.Location != "" {
		fmt.Fprintf(&sb, "LOCATION:%s\n", escapeText(opts.Location))
	}

	if opts.Latitude != 0 || opts.Longitude != 0 {
		fmt.Fprintf(&sb, "GEO:%f;%f\n", opts.Latitude, opts.Longitude)
	}

	if opts.Description != "" {
		fmt.Fprintf(&sb, "DESCRIPTION:%s\n", escapeText(opts.Description))
	}

	sb.WriteString("END:VEVENT\n")
	sb.WriteString("END:VCALENDAR")
	return sb.String()
}

// writeDateTime writes DTSTART/DTEND. An 8-digit date (YYYYMMDD) is an
// all-day value and gets VALUE=DATE instead of a time zone; a UTC time
// (suffix Z) never gets a TZID.
func writeDateTime(sb *strings.Builder, prop, value, tz string) {
	switch {
	case value == "":
		return
	case len(value) == 8:
		fmt.Fprintf(sb, "%s;VALUE=DATE:%s\n", prop, value)
	case tz != "" && !strings.HasSuffix(value, "Z"):
		fmt.Fprintf(sb, "%s;TZID=%s:%s\n", prop, tz, value)
	default:
		fmt.Fprintf(sb, "%s:%s\n", prop, value)
	}
}
