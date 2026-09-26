package main

import "testing"

func TestDecodeBase64Header(t *testing.T) {
	tests := []struct {
		in   string
		want string
		ok   bool
	}{
		{"=?base64?Z2VuZXJhdGVfYmFyY29kZQ==?=", "generate_barcode", true},
		{"=?base64?Z2VuZXJhdGVfYmFyY29kZQ?=", "generate_barcode", true}, // unpadded
		{"=?base64?w6TDtsO8?=", "äöü", true},
		{"generate_barcode", "", false},
		{"=?base64?not base64!?=", "", false},
		{"=?base64?/w==?=", "", false}, // not UTF-8
	}
	for _, tt := range tests {
		got, ok := decodeBase64Header(tt.in)
		if got != tt.want || ok != tt.ok {
			t.Errorf("decodeBase64Header(%q) = %q, %v; want %q, %v", tt.in, got, ok, tt.want, tt.ok)
		}
	}
}
