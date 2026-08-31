package qrformats

import (
	"fmt"
	"strings"
)

// TelOptions holds configuration for telephone call QR codes.
// Standard: RFC 3966.
type TelOptions struct {
	PhoneNumber string
}

// FormatTel returns a standard RFC 3966 telephone URI (e.g. "tel:+49123456789").
func FormatTel(opts TelOptions) string {
	num := strings.TrimSpace(opts.PhoneNumber)
	// Strip spaces and dashes for clean dialing
	num = strings.ReplaceAll(num, " ", "")
	num = strings.ReplaceAll(num, "-", "")
	num = strings.ReplaceAll(num, "(", "")
	num = strings.ReplaceAll(num, ")", "")
	return fmt.Sprintf("tel:%s", num)
}

// SMSOptions holds configuration for SMS message QR codes.
type SMSOptions struct {
	PhoneNumber string
	Message     string
}

// FormatSMS returns a standard "smsto:" URI (widely supported across iOS & Android camera apps).
func FormatSMS(opts SMSOptions) string {
	num := strings.TrimSpace(opts.PhoneNumber)
	num = strings.ReplaceAll(num, " ", "")
	num = strings.ReplaceAll(num, "-", "")
	msg := strings.TrimSpace(opts.Message)
	if msg != "" {
		return fmt.Sprintf("smsto:%s:%s", num, msg)
	}
	return fmt.Sprintf("smsto:%s", num)
}
