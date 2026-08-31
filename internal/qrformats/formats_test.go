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

/*
TestFormatEPC checks the payload line by line, not by substring.

EPC069-12 is positional: the reader takes line 6 as the beneficiary and line 7
as the IBAN, whatever they contain. A `strings.Contains` check passes just as
happily when those two are swapped — and that is the one mistake in a payment
format nobody notices until money moves. Found by making exactly that swap and
watching the old test stay green.
*/
func TestFormatEPC(t *testing.T) {
	got := FormatEPC(EPCOptions{
		Name:      "Max Mustermann",
		IBAN:      "DE89 3704 0044 0532 0130 00",
		BIC:       "GENODEFFXXX",
		Amount:    12.50,
		Reference: "Rechnung-1002",
	})

	// Every element the standard defines, in the order it defines them.
	want := []string{
		"BCD",                    // 1 service tag
		"002",                    // 2 version
		"1",                      // 3 UTF-8
		"SCT",                    // 4 SEPA credit transfer
		"GENODEFFXXX",            // 5 BIC
		"Max Mustermann",         // 6 beneficiary
		"DE89370400440532013000", // 7 IBAN, spaces stripped
		"EUR12.50",               // 8 amount
		"",                       // 9 purpose code
		"",                       // 10 structured reference
		"Rechnung-1002",          // 11 unstructured remittance
	}
	lines := strings.Split(got, "\n")
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d:\n%q", len(lines), len(want), got)
	}
	for i, w := range want {
		if lines[i] != w {
			t.Errorf("line %d = %q, want %q", i+1, lines[i], w)
		}
	}
}

/*
TestFormatEPCOmitsTrailingElements pins what the payload looks like with
nothing optional filled in.

The standard allows trailing empty elements to be dropped, and FormatEPC does
so with a single TrimRight. That is fine for the last elements and would be
wrong for one in the middle — the test above is what guards the middle.
*/
func TestFormatEPCOmitsTrailingElements(t *testing.T) {
	got := FormatEPC(EPCOptions{Name: "Verein e.V.", IBAN: "DE89370400440532013000"})
	want := "BCD\n002\n1\nSCT\n\nVerein e.V.\nDE89370400440532013000"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

/*
TestFormatEPCWithoutAmountLeavesItOpen — an amount of zero is not "EUR0.00",
it is an open amount the payer fills in. Writing a zero would produce a
transfer the bank rejects.
*/
func TestFormatEPCWithoutAmountLeavesItOpen(t *testing.T) {
	got := FormatEPC(EPCOptions{
		Name: "Verein e.V.", IBAN: "DE89370400440532013000", Reference: "Spende",
	})
	lines := strings.Split(got, "\n")
	if lines[7] != "" {
		t.Errorf("amount line = %q, want empty for an open amount", lines[7])
	}
	if lines[10] != "Spende" {
		t.Errorf("remittance line = %q — the reference moved when the amount was empty", lines[10])
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
	if !strings.Contains(email, "subject=Question") || !strings.Contains(email, "body=Please%20help%20me") {
		t.Errorf("missing mail params: %s", email)
	}
}

/*
TestFormatEmailEncodesSpaceAsPercent20 pins the difference between the two
encodings that both look like "URL encoding".

url.Values.Encode writes a space as "+", which is HTML form encoding. RFC 6068
gives "+" no such meaning: in a mailto URI it is a plus sign, and a client that
follows the spec puts "Guten+Tag" in the subject line. Found while documenting
the format — the code looked right and the output was not.
*/
func TestFormatEmailEncodesSpaceAsPercent20(t *testing.T) {
	got := FormatEmail(EmailOptions{To: "a@b.de", Subject: "Guten Tag", Body: "Zeile eins"})
	if strings.Contains(got, "+") {
		t.Errorf("got %q — a space was encoded as '+', which reads as a plus sign", got)
	}
	for _, want := range []string{"subject=Guten%20Tag", "body=Zeile%20eins"} {
		if !strings.Contains(got, want) {
			t.Errorf("got %q, want it to contain %q", got, want)
		}
	}
}

// TestFormatEmailKeepsARealPlus — the rewrite above must not turn an encoded
// plus back into a space. Encode writes it as %2B, which contains no "+".
func TestFormatEmailKeepsARealPlus(t *testing.T) {
	got := FormatEmail(EmailOptions{To: "a@b.de", Subject: "1+1"})
	if !strings.Contains(got, "%2B") {
		t.Errorf("got %q, want the plus preserved as %%2B", got)
	}
}

/*
TestFormatSMSDoesNotEscapeTheMessage records what actually happens rather than
what one would hope.

The payload is "smsto:<number>:<text>" and the text is inserted verbatim, so a
colon inside it produces a third colon. Readers split on the first two and are
unbothered, which is why this is pinned as behaviour and not fixed: changing it
would mean picking an escaping scheme that no reader has agreed to.
*/
func TestFormatSMSDoesNotEscapeTheMessage(t *testing.T) {
	got := FormatSMS(SMSOptions{PhoneNumber: "+4930123", Message: "Treffen: 18:00"})
	if got != "smsto:+4930123:Treffen: 18:00" {
		t.Errorf("got %q", got)
	}
}
