package barcodes

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// The messages here are read by people and by language models calling the
// MCP server. Each one says what is wrong (which character, how much), what
// is allowed and what to do instead — so a model can correct its next call.
// None of them echoes the full input.

const code39Charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 -.$/+%"

// checkCharset rejects input a symbology cannot encode, before the encoder
// sees it.
func checkCharset(btype BarcodeType, data string, opts BarcodeOptions) error {
	switch btype {
	case TypeCode39:
		if opts.FullASCIICode39 {
			return checkASCII(btype, data, "use qr or datamatrix for other characters")
		}
		pos, r := firstRune(data, func(r rune) bool { return !strings.ContainsRune(code39Charset, r) })
		if pos == 0 {
			return nil
		}
		// "upper": upper-casing alone fixes it; "other": the character has no
		// Code 39 form at all.
		fix, hint := "other", "use code128 for full ASCII or qr for any text"
		if strings.ToUpper(data) != data && firstRuneIn(strings.ToUpper(data), code39Charset) == 0 {
			fix, hint = "upper", "convert the text to upper case, or use code128 to keep lower case"
		}
		return inputError(ErrCode39Charset,
			fmt.Sprintf("code39 cannot encode %s at position %d; allowed are A-Z, 0-9, space and - . $ / + %% — %s", quoteRune(r), pos, hint),
			map[string]string{"char": string(r), "pos": itoa(pos), "fix": fix})

	case TypeCode128:
		return checkASCII(btype, data, "use qr, datamatrix or aztec for umlauts and other Unicode text")

	case TypeITF:
		pos, r := firstRune(data, func(r rune) bool { return r < '0' || r > '9' })
		if pos != 0 {
			return inputError(ErrITFDigits,
				fmt.Sprintf("itf encodes digits only; %s at position %d is not a digit — use code128 for letters", quoteRune(r), pos),
				map[string]string{"char": string(r), "pos": itoa(pos)})
		}
		if len(data)%2 != 0 {
			return inputError(ErrITFEven,
				fmt.Sprintf("itf needs an even number of digits, got %d — add a leading 0 (0%s)", len(data), abbreviate(data)),
				map[string]string{"count": itoa(len(data)), "suggestion": "0" + abbreviate(data)})
		}
	}
	return nil
}

func checkASCII(btype BarcodeType, data, hint string) error {
	pos, r := firstRune(data, func(r rune) bool { return r > 127 })
	if pos == 0 {
		return nil
	}
	return inputError(ErrASCIIOnly,
		fmt.Sprintf("%s encodes ASCII only; %s at position %d is not ASCII — %s", btype, quoteRune(r), pos, hint),
		map[string]string{"type": string(btype), "char": string(r), "pos": itoa(pos)})
}

// capacityError explains an encoder failure that is caused by too much
// data. It finds the longest prefix the symbology still accepts, so the
// limit is exact for this very content (digits pack denser than text).
// It returns nil when the input fails for another reason.
func capacityError(btype BarcodeType, data string, encode func(string) error) error {
	runes := []rune(data)
	if len(runes) < 2 || encode(string(runes[:1])) != nil {
		return nil
	}
	lo, hi := 1, len(runes)-1 // lo fits, hi+1 = len(runes) does not
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if encode(string(runes[:mid])) == nil {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return inputError(ErrCapacity,
		fmt.Sprintf("too much data for %s: %d characters (%d bytes), at most %d characters of this content fit — %s",
			btype, len(runes), len(data), lo, capacityHint(btype)),
		map[string]string{"type": string(btype), "chars": itoa(len(runes)), "bytes": itoa(len(data)), "max": itoa(lo)})
}

func capacityHint(btype BarcodeType) string {
	switch btype {
	case TypeDataMatrix:
		return "shorten it, or use qr, which holds more"
	case TypePDF417:
		return "shorten it, or use qr"
	case TypeQR:
		return "shorten it, split it into several codes, or use aztec, which holds more text"
	default:
		return "shorten it or split it into several codes"
	}
}

// encoderError wraps any other library error without repeating the input.
func encoderError(btype BarcodeType, err error) error {
	msg := err.Error()
	if i := strings.IndexByte(msg, '"'); i >= 0 {
		msg = strings.TrimSpace(msg[:i]) + " (input omitted)"
	}
	return inputError(ErrEncoder, fmt.Sprintf("%s cannot encode this input: %s", btype, msg),
		map[string]string{"type": string(btype), "detail": msg})
}

// firstRune returns the 1-based position and value of the first rune that
// matches, or 0 when none does.
func firstRune(s string, match func(rune) bool) (int, rune) {
	pos := 0
	for _, r := range s {
		pos++
		if match(r) {
			return pos, r
		}
	}
	return 0, 0
}

func firstRuneIn(s, charset string) int {
	pos, _ := firstRune(s, func(r rune) bool { return !strings.ContainsRune(charset, r) })
	return pos
}

func quoteRune(r rune) string {
	if r < 32 || r == utf8.RuneError {
		return fmt.Sprintf("character U+%04X", r)
	}
	return fmt.Sprintf("'%c'", r)
}

func abbreviate(s string) string {
	if r := []rune(s); len(r) > 20 {
		return string(r[:20]) + "…"
	}
	return s
}
