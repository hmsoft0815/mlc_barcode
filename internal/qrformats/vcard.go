package qrformats

import (
	"fmt"
	"strings"
)

// VCardOptions holds configuration for a vCard 3.0
type VCardOptions struct {
	FirstName    string
	LastName     string
	Organization string
	Title        string
	Phone        string
	Email        string
	Address      string
	City         string
	Zip          string
	Country      string
	URL          string
}

// FormatVCard returns a formatted vCard 3.0 string (referencing RFC 6350 for vCard 4.0 standards)
func FormatVCard(opts VCardOptions) string {
	var sb strings.Builder
	sb.WriteString("BEGIN:VCARD\n")
	sb.WriteString("VERSION:3.0\n")
	
	if opts.LastName != "" || opts.FirstName != "" {
		fmt.Fprintf(&sb, "N:%s;%s\n", opts.LastName, opts.FirstName)
		fmt.Fprintf(&sb, "FN:%s %s\n", opts.FirstName, opts.LastName)
	}
	
	if opts.Organization != "" {
		fmt.Fprintf(&sb, "ORG:%s\n", opts.Organization)
	}
	
	if opts.Title != "" {
		fmt.Fprintf(&sb, "TITLE:%s\n", opts.Title)
	}
	
	if opts.Phone != "" {
		fmt.Fprintf(&sb, "TEL;TYPE=WORK,VOICE:%s\n", opts.Phone)
	}
	
	if opts.Email != "" {
		fmt.Fprintf(&sb, "EMAIL:%s\n", opts.Email)
	}
	
	if opts.Address != "" || opts.City != "" || opts.Zip != "" || opts.Country != "" {
		fmt.Fprintf(&sb, "ADR;TYPE=WORK:;;%s;%s;;%s;%s\n", opts.Address, opts.City, opts.Zip, opts.Country)
	}
	
	if opts.URL != "" {
		fmt.Fprintf(&sb, "URL:%s\n", opts.URL)
	}
	
	sb.WriteString("END:VCARD")
	return sb.String()
}
