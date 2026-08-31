package qrformats

import (
	"fmt"
	"net/url"
	"strings"
)

// EmailOptions holds configuration for pre-filled email QR codes.
// Standard: RFC 6068 (mailto: URI scheme).
type EmailOptions struct {
	To      string // Recipient email address
	Subject string // Optional subject line
	Body    string // Optional email body
}

// FormatEmail returns an RFC 6068 mailto URI (e.g. "mailto:support@example.com?subject=Hello&body=World").
func FormatEmail(opts EmailOptions) string {
	to := strings.TrimSpace(opts.To)
	if to == "" {
		return ""
	}

	params := url.Values{}
	if opts.Subject != "" {
		params.Set("subject", opts.Subject)
	}
	if opts.Body != "" {
		params.Set("body", opts.Body)
	}

	queryString := params.Encode()
	if queryString != "" {
		return fmt.Sprintf("mailto:%s?%s", to, queryString)
	}
	return fmt.Sprintf("mailto:%s", to)
}
