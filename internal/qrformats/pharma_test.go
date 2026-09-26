package qrformats

import (
	"reflect"
	"testing"
)

func TestParsePharma(t *testing.T) {
	tests := []struct {
		name string
		text string
		want map[string]string
	}{
		{"gs1 as gozxing reports it (leading FNC1 = GS)",
			"\x1d010415012345678217280331" + "10ABC123\x1d21SN0001X9Z",
			map[string]string{"format": "gs1", "gtin": "04150123456782", "gtin_valid": "true", "pzn": "12345678", "pzn_valid": "true",
				"expiry": "280331", "batch": "ABC123", "serial": "SN0001X9Z"}},
		{"gs1 without leading FNC1, serial before batch",
			"0104150123456782" + "21SN0001X9Z\x1d" + "17280300" + "10ABC123",
			map[string]string{"format": "gs1", "gtin": "04150123456782", "gtin_valid": "true", "pzn": "12345678", "pzn_valid": "true",
				"expiry": "280300", "batch": "ABC123", "serial": "SN0001X9Z"}},
		{"ifa format 06",
			"[)>\x1e06\x1d9N111234567842\x1d1TABC123\x1dD280331\x1dSSN0001X9Z\x1e\x04",
			map[string]string{"format": "ifa", "ppn": "111234567842", "ppn_valid": "true", "pzn": "12345678", "pzn_valid": "true",
				"expiry": "280331", "batch": "ABC123", "serial": "SN0001X9Z"}},
		{"wrong check digits are reported, not hidden",
			"[)>\x1e06\x1d9N111234567843\x1d1TX\x1dD280331\x1dS1\x1e\x04",
			map[string]string{"format": "ifa", "ppn": "111234567843", "ppn_valid": "false", "pzn": "12345678", "pzn_valid": "true",
				"expiry": "280331", "batch": "X", "serial": "1"}},
		{"gs1 with wrong GTIN and PZN check digits",
			"\x1d010415012345679817280331" + "10A\x1d21B",
			map[string]string{"format": "gs1", "gtin": "04150123456798", "gtin_valid": "false", "pzn": "12345679", "pzn_valid": "false",
				"expiry": "280331", "batch": "A", "serial": "B"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.text)
			if got.Kind != "pharma" || !reflect.DeepEqual(got.Fields, tt.want) {
				t.Errorf("got %s %v\nwant pharma %v", got.Kind, got.Fields, tt.want)
			}
		})
	}
	for _, text := range []string{"0123", "01 is a number", "4006381333931"} {
		if k := Parse(text).Kind; k == "pharma" {
			t.Errorf("%q taken for a pharma code", text)
		}
	}
}

func TestPZNAndPPNChecks(t *testing.T) {
	if !validPZN("12345678") || validPZN("12345679") || validPZN("1234567") {
		t.Error("PZN check")
	}
	// IFA example: PZN 12345678 → PPN 11 12345678 42
	if !validPPN("111234567842") || validPPN("111234567841") {
		t.Error("PPN check")
	}
}
