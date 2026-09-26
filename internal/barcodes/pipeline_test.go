package barcodes

import (
	"testing"

	"github.com/mlcmcp/mlc_barcode/internal/qrformats"
)

// The whole chain an LLM uses through MCP: format a payload, render it as
// QR, read the image, parse the payload — the fields must survive.
func TestFormatEncodeDecodeParse(t *testing.T) {
	tests := []struct {
		payload string
		kind    string
		check   map[string]string
	}{
		{qrformats.FormatEPC(qrformats.EPCOptions{Name: "Muster GmbH", IBAN: "DE89370400440532013000", Amount: 49.9, Reference: "Rechnung 1234"}),
			"epc", map[string]string{"iban": "DE89370400440532013000", "iban_valid": "true", "amount": "49.90", "reference": "Rechnung 1234"}},
		{qrformats.FormatWifi(qrformats.WifiOptions{SSID: "Gäste;WLAN", Password: "geheim:1"}),
			"wifi", map[string]string{"ssid": "Gäste;WLAN", "password": "geheim:1"}},
		{qrformats.FormatVCard(qrformats.VCardOptions{FirstName: "Erika", LastName: "Müller", Organization: "Muster GmbH; Nord"}),
			"vcard", map[string]string{"last_name": "Müller", "org": "Muster GmbH; Nord"}},
		{qrformats.FormatVCalendar(qrformats.VCalendarOptions{Summary: "Fest, draußen", StartTime: "20260905", EndTime: "20260906"}),
			"event", map[string]string{"summary": "Fest, draußen", "all_day": "true"}},
		{qrformats.FormatCrypto(qrformats.CryptoOptions{Address: "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh", Message: "Rechnung 1234"}),
			"crypto", map[string]string{"message": "Rechnung 1234"}},
	}
	for _, tt := range tests {
		found := roundTrip(t, TypeQR, tt.payload, DefaultOptions(TypeQR))
		if len(found) != 1 {
			t.Fatalf("%s: decoded %d codes", tt.kind, len(found))
		}
		p := qrformats.Parse(found[0].Text)
		if p.Kind != tt.kind {
			t.Errorf("kind %q, want %q (text %q)", p.Kind, tt.kind, found[0].Text)
		}
		for k, v := range tt.check {
			if p.Fields[k] != v {
				t.Errorf("%s: field %s = %q, want %q", tt.kind, k, p.Fields[k], v)
			}
		}
	}
}

// Pharmaceutical pack codes: generated as DataMatrix, read back, parsed.
// The GS1 variant carries its separators as GS characters here (our encoder
// cannot write FNC1); real GS1 codes are covered by the parser tests.
func TestPharmaPackCodes(t *testing.T) {
	for _, tt := range []struct {
		payload, format string
	}{
		{"[)>\x1e06\x1d9N111234567842\x1d1TABC123\x1dD280331\x1dSSN0001X9Z\x1e\x04", "ifa"},
		{"010415012345678217280331" + "10ABC123\x1d21SN0001X9Z", "gs1"},
	} {
		found := roundTrip(t, TypeDataMatrix, tt.payload, DefaultOptions(TypeDataMatrix))
		if len(found) != 1 {
			t.Fatalf("%s: decoded %d codes", tt.format, len(found))
		}
		p := qrformats.Parse(found[0].Text)
		want := map[string]string{"format": tt.format, "pzn": "12345678", "pzn_valid": "true", "batch": "ABC123", "expiry": "280331", "serial": "SN0001X9Z"}
		for k, v := range want {
			if p.Kind != "pharma" || p.Fields[k] != v {
				t.Errorf("%s: %s = %q (kind %s), want %q", tt.format, k, p.Fields[k], p.Kind, v)
			}
		}
	}
}
