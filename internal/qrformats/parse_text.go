package qrformats

import "strings"

// Text-value escaping shared by vCard (RFC 2426) and iCalendar (RFC 5545):
// backslash, comma, semicolon and newline are escaped with a backslash.

func escapeText(s string) string {
	r := strings.NewReplacer(`\`, `\\`, ",", `\,`, ";", `\;`, "\r\n", `\n`, "\n", `\n`)
	return r.Replace(s)
}

func unescapeBackslash(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			i++
			if s[i] == 'n' || s[i] == 'N' {
				b.WriteByte('\n')
			} else {
				b.WriteByte(s[i])
			}
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

// splitEscaped splits at sep, skipping separators preceded by a backslash.
// The parts keep their escapes.
func splitEscaped(s string, sep byte) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\\':
			i++
		case sep:
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	return append(parts, s[start:])
}

// contentLines unfolds RFC 2426/5545 content lines (a line starting with a
// space or tab continues the previous one) and returns name, parameters
// and raw value of each.
type contentLine struct {
	name   string
	params string
	value  string
}

func contentLines(t string) []contentLine {
	raw := strings.Split(strings.ReplaceAll(t, "\r\n", "\n"), "\n")
	var unfolded []string
	for _, l := range raw {
		if (strings.HasPrefix(l, " ") || strings.HasPrefix(l, "\t")) && len(unfolded) > 0 {
			unfolded[len(unfolded)-1] += l[1:]
			continue
		}
		unfolded = append(unfolded, l)
	}
	var out []contentLine
	for _, l := range unfolded {
		head, value, ok := strings.Cut(l, ":")
		if !ok {
			continue
		}
		name, params, _ := strings.Cut(head, ";")
		out = append(out, contentLine{name: strings.ToUpper(name), params: params, value: value})
	}
	return out
}

func param(params, key string) string {
	for _, p := range strings.Split(params, ";") {
		k, v, ok := strings.Cut(p, "=")
		if ok && strings.EqualFold(k, key) {
			return v
		}
	}
	return ""
}

func parseVCardFields(t string) map[string]string {
	f := map[string]string{}
	for _, cl := range contentLines(t) {
		switch cl.name {
		case "N":
			parts := splitEscaped(cl.value, ';')
			if len(parts) > 0 {
				put(f, "last_name", unescapeBackslash(parts[0]))
			}
			if len(parts) > 1 {
				put(f, "first_name", unescapeBackslash(parts[1]))
			}
		case "FN":
			put(f, "full_name", unescapeBackslash(cl.value))
		case "ORG":
			put(f, "org", unescapeBackslash(cl.value))
		case "TITLE":
			put(f, "title", unescapeBackslash(cl.value))
		case "TEL":
			put(f, "phone", unescapeBackslash(cl.value))
		case "EMAIL":
			put(f, "email", unescapeBackslash(cl.value))
		case "URL":
			put(f, "url", unescapeBackslash(cl.value))
		case "ADR": // PO box;extended;street;city;region;zip;country
			parts := splitEscaped(cl.value, ';')
			for i, key := range []string{"", "", "street", "city", "region", "zip", "country"} {
				if key != "" && i < len(parts) {
					put(f, key, unescapeBackslash(parts[i]))
				}
			}
		}
	}
	return f
}

func parseEventFields(t string) map[string]string {
	f := map[string]string{}
	for _, cl := range contentLines(t) {
		switch cl.name {
		case "SUMMARY":
			put(f, "summary", unescapeBackslash(cl.value))
		case "DESCRIPTION":
			put(f, "description", unescapeBackslash(cl.value))
		case "LOCATION":
			put(f, "location", unescapeBackslash(cl.value))
		case "DTSTART", "DTEND":
			key := "start"
			if cl.name == "DTEND" {
				key = "end"
			}
			put(f, key, cl.value)
			put(f, "timezone", param(cl.params, "TZID"))
			if strings.EqualFold(param(cl.params, "VALUE"), "DATE") || len(cl.value) == 8 {
				f["all_day"] = "true"
			}
		case "GEO":
			if lat, lon, ok := strings.Cut(cl.value, ";"); ok {
				f["latitude"], f["longitude"] = lat, lon
			}
		}
	}
	return f
}
