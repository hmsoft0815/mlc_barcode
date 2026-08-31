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
			t.Errorf("FormatVCard() = %v, must contain %v", got, s)
		}
	}
}

func TestFormatVCalendar(t *testing.T) {
	opts := VCalendarOptions{
		Summary:     "Team Meeting",
		Description: "Discuss roadmap",
		Location:    "Conference Room A",
		StartTime:   "20260420T100000Z",
		EndTime:     "20260420T110000Z",
	}
	got := FormatVCalendar(opts)

	expect := []string{
		"BEGIN:VCALENDAR",
		"VERSION:2.0",
		"BEGIN:VEVENT",
		"SUMMARY:Team Meeting",
		"DESCRIPTION:Discuss roadmap",
		"LOCATION:Conference Room A",
		"DTSTART:20260420T100000Z",
		"DTEND:20260420T110000Z",
		"END:VEVENT",
		"END:VCALENDAR",
	}

	for _, s := range expect {
		if !strings.Contains(got, s) {
			t.Errorf("FormatVCalendar() = %v, must contain %v", got, s)
		}
	}
}

func TestFormatEPC(t *testing.T) {
	opts := EPCOptions{
		Name:      "Max Mustermann",
		IBAN:      "DE89 3704 0044 0532 0130 00",
		BIC:       "GENODEFFXXX",
		Amount:    12.50,
		Reference: "Rechnung-1002",
	}
	got := FormatEPC(opts)

	expect := []string{
		"BCD\n002\n1\nSCT",
		"GENODEFFXXX",
		"Max Mustermann",
		"DE89370400440532013000",
		"EUR12.50",
		"Rechnung-1002",
	}

	for _, s := range expect {
		if !strings.Contains(got, s) {
			t.Errorf("FormatEPC() = %v, must contain %v", got, s)
		}
	}
}

func TestFormatCrypto(t *testing.T) {
	// Bitcoin
	btc := FormatCrypto(CryptoOptions{
		Coin:    "btc",
		Address: "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
		Amount:  0.05,
		Label:   "Donation",
	})
	if !strings.HasPrefix(btc, "bitcoin:1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa?") {
		t.Errorf("unexpected btc URI: %s", btc)
	}
	if !strings.Contains(btc, "amount=0.05") || !strings.Contains(btc, "label=Donation") {
		t.Errorf("missing params in btc URI: %s", btc)
	}

	// Ethereum
	eth := FormatCrypto(CryptoOptions{
		Coin:    "eth",
		Address: "0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae",
	})
	if eth != "ethereum:0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae" {
		t.Errorf("unexpected eth URI: %s", eth)
	}
}

func TestFormatGeo(t *testing.T) {
	geo := FormatGeo(GeoOptions{
		Latitude:  52.520008,
		Longitude: 13.404954,
		Query:     "Berlin Fernsehturm",
	})
	if !strings.HasPrefix(geo, "geo:52.520008,13.404954?q=") {
		t.Errorf("unexpected geo URI: %s", geo)
	}
}

func TestFormatTelAndSMS(t *testing.T) {
	tel := FormatTel(TelOptions{PhoneNumber: "+49 (0) 123 456-789"})
	if tel != "tel:+490123456789" {
		t.Errorf("unexpected tel URI: %s", tel)
	}

	sms := FormatSMS(SMSOptions{
		PhoneNumber: "+49 123 456789",
		Message:     "Hello World",
	})
	if sms != "smsto:+49123456789:Hello World" {
		t.Errorf("unexpected sms URI: %s", sms)
	}
}

func TestFormatEmail(t *testing.T) {
	email := FormatEmail(EmailOptions{
		To:      "support@example.com",
		Subject: "Question",
		Body:    "Please help me",
	})
	if !strings.HasPrefix(email, "mailto:support@example.com?") {
		t.Errorf("unexpected mailto URI: %s", email)
	}
	if !strings.Contains(email, "subject=Question") || !strings.Contains(email, "body=Please+help+me") {
		t.Errorf("missing mail params: %s", email)
	}
}
