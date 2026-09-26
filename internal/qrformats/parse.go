package qrformats

import (
	"math/big"
	"net/url"
	"strings"
)

// Parsed is a decoded payload split into fields. Kind is one of epc, wifi,
// vcard, event, geo, tel, sms, email, crypto, url, pharma or text; Fields holds
// only what the payload carries. It is the inverse of the Format*
// functions, so a round trip returns the same values.
type Parsed struct {
	Kind   string            `json:"kind"`
	Fields map[string]string `json:"fields,omitempty"`
}

// Parse recognises the payload formats this package writes (and the
// common variants other generators use).
func Parse(text string) Parsed {
	t := strings.TrimSpace(text)
	lower := strings.ToLower(t)
	if isPharma(text) { // control characters matter: check before trimming
		return parsePharma(text)
	}
	switch {
	case strings.HasPrefix(t, "BCD\n") || strings.HasPrefix(t, "BCD\r\n"):
		return parseEPC(t)
	case strings.HasPrefix(lower, "wifi:"):
		return parseWifi(t[5:])
	case strings.HasPrefix(lower, "begin:vcard"):
		return Parsed{Kind: "vcard", Fields: parseVCardFields(t)}
	case strings.HasPrefix(lower, "begin:vcalendar") || strings.HasPrefix(lower, "begin:vevent"):
		return Parsed{Kind: "event", Fields: parseEventFields(t)}
	case strings.HasPrefix(lower, "geo:"):
		return parseGeo(t[4:])
	case strings.HasPrefix(lower, "tel:"):
		return Parsed{Kind: "tel", Fields: map[string]string{"phone": t[4:]}}
	case strings.HasPrefix(lower, "smsto:") || strings.HasPrefix(lower, "sms:"):
		return parseSMS(t)
	case strings.HasPrefix(lower, "mailto:"):
		return parseMailto(t[7:])
	case strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://"):
		return Parsed{Kind: "url", Fields: map[string]string{"url": t}}
	}
	if scheme, rest, ok := strings.Cut(t, ":"); ok && cryptoSchemes[strings.ToLower(scheme)] {
		return parseCrypto(strings.ToLower(scheme), rest)
	}
	return Parsed{Kind: "text", Fields: map[string]string{"text": text}}
}

var cryptoSchemes = map[string]bool{"bitcoin": true, "ethereum": true, "solana": true, "dogecoin": true, "litecoin": true}

func put(m map[string]string, key, value string) {
	if value != "" {
		m[key] = value
	}
}

// parseEPC reads an EPC069-12 (GiroCode) payload: fixed lines.
func parseEPC(t string) Parsed {
	lines := strings.Split(strings.ReplaceAll(t, "\r\n", "\n"), "\n")
	line := func(i int) string {
		if i < len(lines) {
			return strings.TrimSpace(lines[i])
		}
		return ""
	}
	f := map[string]string{}
	put(f, "bic", line(4))
	put(f, "name", line(5))
	put(f, "iban", line(6))
	if amount := line(7); len(amount) > 3 {
		f["currency"], f["amount"] = amount[:3], amount[3:]
	}
	put(f, "purpose", line(8))
	put(f, "reference", line(9)) // structured creditor reference
	if line(9) == "" {
		put(f, "reference", line(10))
	}
	put(f, "note", line(11))
	if iban := f["iban"]; iban != "" {
		if ValidIBAN(iban) {
			f["iban_valid"] = "true"
		} else {
			f["iban_valid"] = "false"
		}
	}
	return Parsed{Kind: "epc", Fields: f}
}

// ValidIBAN checks the ISO 13616 mod-97 checksum.
func ValidIBAN(iban string) bool {
	s := strings.ToUpper(strings.ReplaceAll(iban, " ", ""))
	if len(s) < 15 || len(s) > 34 {
		return false
	}
	var digits strings.Builder
	for _, r := range s[4:] + s[:4] {
		switch {
		case r >= '0' && r <= '9':
			digits.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			digits.WriteString(big.NewInt(int64(r - 'A' + 10)).String())
		default:
			return false
		}
	}
	n, ok := new(big.Int).SetString(digits.String(), 10)
	return ok && new(big.Int).Mod(n, big.NewInt(97)).Int64() == 1
}

// parseWifi reads T:…;S:…;P:…;H:…; with backslash escapes.
func parseWifi(body string) Parsed {
	f := map[string]string{}
	for _, part := range splitEscaped(body, ';') {
		key, value, ok := strings.Cut(part, ":")
		if !ok {
			continue
		}
		value = unescapeBackslash(value)
		switch strings.ToUpper(key) {
		case "T":
			put(f, "encryption", value)
		case "S":
			put(f, "ssid", value)
		case "P":
			put(f, "password", value)
		case "H":
			if strings.EqualFold(value, "true") {
				f["hidden"] = "true"
			}
		}
	}
	return Parsed{Kind: "wifi", Fields: f}
}

func parseGeo(body string) Parsed {
	f := map[string]string{}
	coords, query, _ := strings.Cut(body, "?")
	coords, _, _ = strings.Cut(coords, ";") // drop ;crs=… / ;u=…
	parts := strings.Split(coords, ",")
	if len(parts) >= 2 {
		f["latitude"], f["longitude"] = strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
	if q, err := url.ParseQuery(query); err == nil {
		put(f, "query", q.Get("q"))
	}
	return Parsed{Kind: "geo", Fields: f}
}

func parseSMS(t string) Parsed {
	f := map[string]string{}
	if strings.HasPrefix(strings.ToLower(t), "smsto:") {
		num, msg, _ := strings.Cut(t[6:], ":")
		put(f, "phone", num)
		put(f, "message", msg)
	} else { // sms:NUMBER?body=…
		num, query, _ := strings.Cut(t[4:], "?")
		put(f, "phone", num)
		if q, err := url.ParseQuery(query); err == nil {
			put(f, "message", q.Get("body"))
		}
	}
	return Parsed{Kind: "sms", Fields: f}
}

func parseMailto(body string) Parsed {
	f := map[string]string{}
	to, query, _ := strings.Cut(body, "?")
	if dec, err := url.PathUnescape(to); err == nil {
		to = dec
	}
	put(f, "to", to)
	for _, kv := range strings.Split(query, "&") {
		key, value, _ := strings.Cut(kv, "=")
		// RFC 6068: '+' is a literal plus, spaces are %20.
		if dec, err := url.PathUnescape(value); err == nil {
			value = dec
		}
		switch strings.ToLower(key) {
		case "subject":
			put(f, "subject", value)
		case "body":
			put(f, "body", value)
		}
	}
	return Parsed{Kind: "email", Fields: f}
}

func parseCrypto(scheme, rest string) Parsed {
	f := map[string]string{"coin": scheme}
	addr, query, _ := strings.Cut(rest, "?")
	put(f, "address", addr)
	for _, kv := range strings.Split(query, "&") {
		key, value, _ := strings.Cut(kv, "=")
		if dec, err := url.PathUnescape(value); err == nil {
			value = dec
		}
		switch strings.ToLower(key) {
		case "amount", "label", "message":
			put(f, strings.ToLower(key), value)
		}
	}
	return Parsed{Kind: "crypto", Fields: f}
}
