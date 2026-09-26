package barcodes

import (
	"fmt"
	"strconv"
)

// Reasons a retail code (EAN-13, EAN-8, UPC-A) is rejected.
const (
	ReasonNonDigit = "non_digit"
	ReasonLength   = "length"
	ReasonChecksum = "checksum"
)

// RetailCheck is the result of checking an EAN/UPC input.
type RetailCheck struct {
	Valid bool
	// Reason is empty when Valid, otherwise one of the Reason* constants.
	Reason string
	// Code is the complete code including the check digit: the input when it
	// was complete, the input plus the computed digit when it was one short,
	// and the corrected code on a checksum mismatch.
	Code string
	// CheckDigitAdded is set when the input lacked the check digit.
	CheckDigitAdded bool
	// Given and Expected are the check digits on a checksum mismatch.
	Given, Expected int
	// Lengths are the accepted input lengths (without / with check digit).
	Lengths [2]int
}

// IsRetail reports whether the type carries a GS1 check digit.
func IsRetail(btype BarcodeType) bool {
	return btype == TypeEAN13 || btype == TypeEAN8 || btype == TypeUPCA
}

func retailLength(btype BarcodeType) int {
	switch btype {
	case TypeEAN13:
		return 13
	case TypeEAN8:
		return 8
	case TypeUPCA:
		return 12
	}
	return 0
}

// CheckRetail validates digits, length and check digit of an EAN-13, EAN-8
// or UPC-A input. Other types are always valid.
func CheckRetail(btype BarcodeType, data string) RetailCheck {
	full := retailLength(btype)
	if full == 0 {
		return RetailCheck{Valid: true, Code: data}
	}
	res := RetailCheck{Lengths: [2]int{full - 1, full}}

	for _, r := range data {
		if r < '0' || r > '9' {
			res.Reason = ReasonNonDigit
			return res
		}
	}

	switch len(data) {
	case full - 1:
		res.Valid = true
		res.CheckDigitAdded = true
		res.Code = data + strconv.Itoa(gs1CheckDigit(data))
	case full:
		body := data[:full-1]
		res.Expected = gs1CheckDigit(body)
		res.Given = int(data[full-1] - '0')
		res.Code = body + strconv.Itoa(res.Expected)
		res.Valid = res.Given == res.Expected
		if !res.Valid {
			res.Reason = ReasonChecksum
		}
	default:
		res.Reason = ReasonLength
	}
	return res
}

// Error describes an invalid result in English (CLI, MCP, batch).
func (c RetailCheck) Error() string {
	switch c.Reason {
	case ReasonNonDigit:
		return "only digits are allowed"
	case ReasonLength:
		return fmt.Sprintf("needs %d digits (or %d without check digit)", c.Lengths[1], c.Lengths[0])
	case ReasonChecksum:
		return fmt.Sprintf("wrong check digit %d, expected %d (%s)", c.Given, c.Expected, c.Code)
	}
	return ""
}

// gs1CheckDigit computes the GS1 mod-10 check digit: from the right, digits
// are weighted 3, 1, 3, 1 …
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
