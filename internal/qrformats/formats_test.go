package qrformats

import (
	"strings"
	"testing"
)

func TestFormatWifi(t *testing.T) {
	tests := []struct {
		name     string
		opts     WifiOptions
		contains []string
	}{
		{
			"Standard WPA",
			WifiOptions{SSID: "MyNetwork", Password: "mypassword", Encryption: "WPA"},
			[]string{"WIFI:T:WPA", "S:MyNetwork", "P:mypassword"},
		},
		{
			"Hidden Network",
			WifiOptions{SSID: "Hidden", Password: "secret", Hidden: true},
			[]string{"S:Hidden", "P:secret", "H:true"},
		},
		{
			"Special Characters",
			WifiOptions{SSID: "My:Net;work", Password: "pass:word;"},
			[]string{"S:My\\:Net\\;work", "P:pass\\:word\\;"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatWifi(tt.opts)
			for _, s := range tt.contains {
				if !strings.Contains(got, s) {
					t.Errorf("FormatWifi() = %v, must contain %v", got, s)
				}
			}
		})
	}
}

func TestFormatVCard(t *testing.T) {
	opts := VCardOptions{
		FirstName:    "Max",
		LastName:     "Mustermann",
		Organization: "Example Corp",
		Email:        "max@example.com",
		Phone:        "+491234567",
	}
	got := FormatVCard(opts)

	expect := []string{
		"BEGIN:VCARD",
		"VERSION:3.0",
		"N:Mustermann;Max",
		"FN:Max Mustermann",
		"ORG:Example Corp",
		"EMAIL:max@example.com",
		"TEL;TYPE=WORK,VOICE:+491234567",
		"END:VCARD",
	}

	for _, s := range expect {
		if !strings.Contains(got, s) {
			t.Errorf("FormatVCard() missing %v", s)
		}
	}
}

func TestFormatVCalendar(t *testing.T) {
	opts := VCalendarOptions{
		Summary:   "Meeting",
		StartTime: "20260317T100000",
		EndTime:   "20260317T110000",
		TimeZone:  "Europe/Berlin",
		Latitude:  52.5200,
		Longitude: 13.4050,
	}
	got := FormatVCalendar(opts)

	expect := []string{
		"BEGIN:VCALENDAR",
		"VERSION:2.0",
		"BEGIN:VEVENT",
		"SUMMARY:Meeting",
		"DTSTART;TZID=Europe/Berlin:20260317T100000",
		"DTEND;TZID=Europe/Berlin:20260317T110000",
		"GEO:52.520000;13.405000",
		"END:VEVENT",
		"END:VCALENDAR",
	}

	for _, s := range expect {
		if !strings.Contains(got, s) {
			t.Errorf("FormatVCalendar() missing %v", s)
		}
	}
}
