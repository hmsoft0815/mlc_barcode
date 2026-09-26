package qrformats

import (
	"reflect"
	"strings"
	"testing"
)

// Round trip: Format* → Parse must give back what went in — including the
// characters each format has to escape.
func TestParseRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		kind    string
		want    map[string]string
	}{
		{"epc", FormatEPC(EPCOptions{Name: "Muster GmbH", IBAN: "DE89 3704 0044 0532 0130 00", BIC: "COBADEFFXXX", Amount: 49.9, Reference: "Rechnung 1234"}),
			"epc", map[string]string{"name": "Muster GmbH", "iban": "DE89370400440532013000", "iban_valid": "true", "bic": "COBADEFFXXX", "currency": "EUR", "amount": "49.90", "reference": "Rechnung 1234"}},
		{"epc wrong iban", FormatEPC(EPCOptions{Name: "X", IBAN: "DE89370400440532013001", Amount: 1}),
			"epc", map[string]string{"name": "X", "iban": "DE89370400440532013001", "iban_valid": "false", "currency": "EUR", "amount": "1.00"}},
		{"wifi with specials", FormatWifi(WifiOptions{SSID: `Gast;Netz:1`, Password: `p\a,s"s`, Encryption: "WPA", Hidden: true}),
			"wifi", map[string]string{"ssid": `Gast;Netz:1`, "password": `p\a,s"s`, "encryption": "WPA", "hidden": "true"}},
		{"vcard with specials", FormatVCard(VCardOptions{FirstName: "Erika", LastName: "Mustermann", Organization: "Muster GmbH; Filiale Nord", Title: "Leitung, Einkauf", Phone: "+49 30 123", Email: "erika@example.org", Address: "Hauptstr. 1", City: "Berlin", Zip: "10115", Country: "DE", URL: "https://example.org"}),
			"vcard", map[string]string{"first_name": "Erika", "last_name": "Mustermann", "full_name": "Erika Mustermann", "org": "Muster GmbH; Filiale Nord", "title": "Leitung, Einkauf", "phone": "+49 30 123", "email": "erika@example.org", "street": "Hauptstr. 1", "city": "Berlin", "zip": "10115", "country": "DE", "url": "https://example.org"}},
		{"event with specials", FormatVCalendar(VCalendarOptions{Summary: "Sommerfest, draußen; mit Grill", Location: "Hof\nHaus B", Description: `Pfad C:\Fest`, StartTime: "20260901T180000", EndTime: "20260901T230000", TimeZone: "Europe/Berlin"}),
			"event", map[string]string{"summary": "Sommerfest, draußen; mit Grill", "location": "Hof\nHaus B", "description": `Pfad C:\Fest`, "start": "20260901T180000", "end": "20260901T230000", "timezone": "Europe/Berlin"}},
		{"event all day", FormatVCalendar(VCalendarOptions{Summary: "Ausflug", StartTime: "20260905", EndTime: "20260906"}),
			"event", map[string]string{"summary": "Ausflug", "start": "20260905", "end": "20260906", "all_day": "true"}},
		{"geo", FormatGeo(GeoOptions{Latitude: 52.52, Longitude: 13.405, Query: "Brandenburger Tor"}),
			"geo", map[string]string{"latitude": "52.520000", "longitude": "13.405000", "query": "Brandenburger Tor"}},
		{"tel", FormatTel(TelOptions{PhoneNumber: "+49 (30) 123-45"}), "tel", map[string]string{"phone": "+493012345"}},
		{"sms", FormatSMS(SMSOptions{PhoneNumber: "+49 30 123", Message: "Hallo: bis gleich"}),
			"sms", map[string]string{"phone": "+4930123", "message": "Hallo: bis gleich"}},
		{"email", FormatEmail(EmailOptions{To: "info@example.org", Subject: "Frage zu 1+1", Body: "Hallo Welt"}),
			"email", map[string]string{"to": "info@example.org", "subject": "Frage zu 1+1", "body": "Hallo Welt"}},
		{"crypto", FormatCrypto(CryptoOptions{Coin: "btc", Address: "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh", Amount: 0.001, Label: "Muster GmbH", Message: "Rechnung 1234"}),
			"crypto", map[string]string{"coin": "bitcoin", "address": "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh", "amount": "0.001", "label": "Muster GmbH", "message": "Rechnung 1234"}},
		{"url", "https://mlcgo.eu/produkte", "url", map[string]string{"url": "https://mlcgo.eu/produkte"}},
		{"text", "Nur ein Text", "text", map[string]string{"text": "Nur ein Text"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.payload)
			if got.Kind != tt.kind || !reflect.DeepEqual(got.Fields, tt.want) {
				t.Errorf("Parse(%q)\n got  %s %v\n want %s %v", tt.payload, got.Kind, got.Fields, tt.kind, tt.want)
			}
		})
	}
}

func TestValidIBAN(t *testing.T) {
	for iban, want := range map[string]bool{
		"DE89370400440532013000":      true,
		"DE89 3704 0044 0532 0130 00": true,
		"GB82WEST12345698765432":      true,
		"DE89370400440532013001":      false,
		"DE00":                        false,
		"DE89-3704":                   false,
	} {
		if got := ValidIBAN(iban); got != want {
			t.Errorf("ValidIBAN(%q) = %v", iban, got)
		}
	}
}

// vCard (RFC 2426) and iCalendar (RFC 5545) text values must escape
// backslash, comma, semicolon and newline — a lenient parser would not
// notice, but address books split "GmbH; Filiale" into two org units.
func TestFormattersEscapeText(t *testing.T) {
	vcard := FormatVCard(VCardOptions{FirstName: "A,B", LastName: "C;D", Organization: "Muster GmbH; Filiale Nord", Title: "Leitung, Einkauf", City: "Berlin; Mitte"})
	for _, want := range []string{`N:C\;D;A\,B`, `ORG:Muster GmbH\; Filiale Nord`, `TITLE:Leitung\, Einkauf`, `Berlin\; Mitte`} {
		if !strings.Contains(vcard, want) {
			t.Errorf("vCard lacks %q:\n%s", want, vcard)
		}
	}
	event := FormatVCalendar(VCalendarOptions{Summary: "a, b; c", Location: "Hof\nHaus B", Description: `C:\x`})
	for _, want := range []string{`SUMMARY:a\, b\; c`, `LOCATION:Hof\nHaus B`, `DESCRIPTION:C:\\x`} {
		if !strings.Contains(event, want) {
			t.Errorf("iCalendar lacks %q:\n%s", want, event)
		}
	}
	if crypto := FormatCrypto(CryptoOptions{Address: "bc1q", Message: "Rechnung 1234"}); strings.Contains(crypto, "+") {
		t.Errorf("BIP 21 needs %%20 for spaces, got %s", crypto)
	}
}
