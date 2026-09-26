package qrformats

import (
	"strconv"
	"strings"
)

// Pharmaceutical pack codes (securPharm, EU Falsified Medicines Directive):
// a DataMatrix in one of two formats.
//
//   - GS1: Application Identifiers 01 (GTIN/NTIN), 17 (expiry YYMMDD),
//     10 (batch), 21 (serial), optionally 710 (German NHRN = PZN). Readers
//     report the FNC1 characters as ASCII 29 (GS), the first one leading.
//   - IFA: ISO/IEC 15434 format 06, "[)>" RS "06" GS, data identifiers 9N
//     (PPN), 1T (batch), D (expiry), S (serial), separated by GS, closed by
//     RS EOT.
//
// The German PZN is derived from the PPN ("11" + PZN + 2 check digits) or
// the NTIN ("04150" + PZN + check digit); all check digits are verified.

const (
	gs  = "\x1d"
	rs  = "\x1e"
	eot = "\x04"
)

// isPharma reports whether text looks like a securPharm pack code.
func isPharma(text string) bool {
	if strings.HasPrefix(text, "[)>"+rs+"06") {
		return true
	}
	t := strings.TrimPrefix(text, gs)
	return len(t) >= 16 && strings.HasPrefix(t, "01") && allDigits(t[2:16]) &&
		(strings.HasPrefix(text, gs) || strings.Contains(t[16:], gs) ||
			strings.HasPrefix(t[16:], "17") || strings.HasPrefix(t[16:], "10") || strings.HasPrefix(t[16:], "21"))
}

func parsePharma(text string) Parsed {
	f := map[string]string{}
	if strings.HasPrefix(text, "[)>"+rs+"06") {
		f["format"] = "ifa"
		parseIFA(text, f)
	} else {
		f["format"] = "gs1"
		parseGS1(strings.TrimPrefix(text, gs), f)
	}
	derivePZN(f)
	return Parsed{Kind: "pharma", Fields: f}
}

func parseIFA(text string, f map[string]string) {
	body := strings.TrimPrefix(text, "[)>"+rs+"06")
	body = strings.TrimSuffix(strings.TrimSuffix(body, eot), rs)
	for _, part := range strings.Split(body, gs) {
		switch {
		case strings.HasPrefix(part, "9N"):
			put(f, "ppn", part[2:])
		case strings.HasPrefix(part, "1T"):
			put(f, "batch", part[2:])
		case strings.HasPrefix(part, "8P"):
			put(f, "gtin", part[2:])
		case strings.HasPrefix(part, "D"):
			put(f, "expiry", part[1:])
		case strings.HasPrefix(part, "S"):
			put(f, "serial", part[1:])
		}
	}
}

// gs1Fixed lists AIs of fixed length (without the AI) used on packs; all
// others run to the next GS or the end.
var gs1Fixed = map[string]int{"01": 14, "11": 6, "17": 6}

var gs1Names = map[string]string{"01": "gtin", "10": "batch", "11": "production_date", "17": "expiry", "21": "serial", "710": "nhrn"}

func parseGS1(t string, f map[string]string) {
	for len(t) > 0 {
		ai := ""
		for _, candidate := range []string{"710", "711", "712", "713", "714", "01", "10", "11", "17", "21"} {
			if strings.HasPrefix(t, candidate) {
				ai = candidate
				break
			}
		}
		if ai == "" {
			f["unparsed"] = t
			return
		}
		t = t[len(ai):]
		var value string
		if n, fixed := gs1Fixed[ai]; fixed {
			value, t = t[:min(n, len(t))], t[min(n, len(t)):]
		} else if i := strings.Index(t, gs); i >= 0 {
			value, t = t[:i], t[i:]
		} else {
			value, t = t, ""
		}
		t = strings.TrimPrefix(t, gs)
		name := gs1Names[ai]
		if name == "" {
			name = "ai_" + ai
		}
		put(f, name, value)
	}
}

// derivePZN adds pzn and the check-digit verdicts.
func derivePZN(f map[string]string) {
	if ppn := f["ppn"]; ppn != "" {
		f["ppn_valid"] = strconv.FormatBool(validPPN(ppn))
		if strings.HasPrefix(ppn, "11") && len(ppn) == 12 {
			f["pzn"] = ppn[2:10]
		}
	}
	if gtin := f["gtin"]; gtin != "" {
		f["gtin_valid"] = strconv.FormatBool(len(gtin) == 14 && allDigits(gtin) && gs1CheckDigit(gtin[:13]) == int(gtin[13]-'0'))
		if strings.HasPrefix(gtin, "04150") && f["pzn"] == "" {
			f["pzn"] = gtin[5:13]
		}
	}
	if nhrn := f["nhrn"]; nhrn != "" && f["pzn"] == "" {
		f["pzn"] = nhrn
	}
	if pzn := f["pzn"]; pzn != "" {
		f["pzn_valid"] = strconv.FormatBool(validPZN(pzn))
	}
}

// validPZN checks a PZN8: digits 1-7 weighted 1..7, sum mod 11 is digit 8
// (a remainder of 10 is never issued).
func validPZN(pzn string) bool {
	if len(pzn) != 8 || !allDigits(pzn) {
		return false
	}
	sum := 0
	for i := 0; i < 7; i++ {
		sum += int(pzn[i]-'0') * (i + 1)
	}
	return sum%11 != 10 && sum%11 == int(pzn[7]-'0')
}

// validPPN checks the IFA Pharmacy Product Number: character codes weighted
// 2, 3, 4 … from the left, sum mod 97 gives the last two digits.
func validPPN(ppn string) bool {
	if len(ppn) < 3 || !allDigits(ppn[len(ppn)-2:]) {
		return false
	}
	body := ppn[:len(ppn)-2]
	sum := 0
	for i := 0; i < len(body); i++ {
		sum += int(body[i]) * (i + 2)
	}
	check, _ := strconv.Atoi(ppn[len(ppn)-2:])
	return sum%97 == check
}

// gs1CheckDigit is the GS1 mod-10 check digit (weights 3, 1 from the right).
func gs1CheckDigit(body string) int {
	sum := 0
	for i := len(body) - 1; i >= 0; i-- {
		d := int(body[i] - '0')
		if (len(body)-1-i)%2 == 0 {
			d *= 3
		}
		sum += d
	}
	return (10 - sum%10) % 10
}

func allDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return s != ""
}
