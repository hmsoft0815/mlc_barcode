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
		sb.WriteString(fmt.Sprintf("N:%s;%s\n", opts.LastName, opts.FirstName))
		sb.WriteString(fmt.Sprintf("FN:%s %s\n", opts.FirstName, opts.LastName))
	}
	
	if opts.Organization != "" {
		sb.WriteString(fmt.Sprintf("ORG:%s\n", opts.Organization))
	}
	
	if opts.Title != "" {
		sb.WriteString(fmt.Sprintf("TITLE:%s\n", opts.Title))
	}
	
	if opts.Phone != "" {
		sb.WriteString(fmt.Sprintf("TEL;TYPE=WORK,VOICE:%s\n", opts.Phone))
	}
	
	if opts.Email != "" {
		sb.WriteString(fmt.Sprintf("EMAIL:%s\n", opts.Email))
	}
	
	if opts.Address != "" || opts.City != "" || opts.Zip != "" || opts.Country != "" {
		sb.WriteString(fmt.Sprintf("ADR;TYPE=WORK:;;%s;%s;;%s;%s\n", opts.Address, opts.City, opts.Zip, opts.Country))
	}
	
	if opts.URL != "" {
		sb.WriteString(fmt.Sprintf("URL:%s\n", opts.URL))
	}
	
	sb.WriteString("END:VCARD")
	return sb.String()
}
