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
	StartTime   string // Format: YYYYMMDDTHHMMSS (local) or YYYYMMDDTHHMMSSZ (UTC)
	EndTime     string // Format: YYYYMMDDTHHMMSS (local) or YYYYMMDDTHHMMSSZ (UTC)
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
		sb.WriteString(fmt.Sprintf("SUMMARY:%s\n", opts.Summary))
	}
	
	tzPrefix := ""
	if opts.TimeZone != "" {
		tzPrefix = fmt.Sprintf(";TZID=%s", opts.TimeZone)
	}

	if opts.StartTime != "" {
		// If TimeZone is set and StartTime doesn't end with Z, use TZID
		if opts.TimeZone != "" && !strings.HasSuffix(opts.StartTime, "Z") {
			sb.WriteString(fmt.Sprintf("DTSTART%s:%s\n", tzPrefix, opts.StartTime))
		} else {
			sb.WriteString(fmt.Sprintf("DTSTART:%s\n", opts.StartTime))
		}
	}
	
	if opts.EndTime != "" {
		if opts.TimeZone != "" && !strings.HasSuffix(opts.EndTime, "Z") {
			sb.WriteString(fmt.Sprintf("DTEND%s:%s\n", tzPrefix, opts.EndTime))
		} else {
			sb.WriteString(fmt.Sprintf("DTEND:%s\n", opts.EndTime))
		}
	}
	
	if opts.Location != "" {
		sb.WriteString(fmt.Sprintf("LOCATION:%s\n", opts.Location))
	}

	if opts.Latitude != 0 || opts.Longitude != 0 {
		sb.WriteString(fmt.Sprintf("GEO:%f;%f\n", opts.Latitude, opts.Longitude))
	}
	
	if opts.Description != "" {
		sb.WriteString(fmt.Sprintf("DESCRIPTION:%s\n", opts.Description))
	}
	
	sb.WriteString("END:VEVENT\n")
	sb.WriteString("END:VCALENDAR")
	return sb.String()
}
