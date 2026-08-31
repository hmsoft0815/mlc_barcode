package qrformats

import (
	"fmt"
	"strings"
)

// EPCOptions holds configuration for an EPC-QR-Code (GiroCode / SEPA Credit Transfer).
// Standard: European Payments Council EPC069-08.
type EPCOptions struct {
	Name      string  // Beneficiary Name (Empfänger, max 70 chars)
	IBAN      string  // Beneficiary IBAN (max 34 chars)
	BIC       string  // Beneficiary BIC (max 11 chars, optional in SEPA version 002)
	Amount    float64 // Amount in EUR (e.g. 12.50)
	Reference string  // Remittance Information / Verwendungszweck (max 140 chars)
	Purpose   string  // Optional Purpose Code (e.g. "CHAR", "GDDS", max 4 chars)
}

// FormatEPC returns a standard European Payments Council Quick Response Code (GiroCode) payload.
func FormatEPC(opts EPCOptions) string {
	iban := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(opts.IBAN), " ", ""))
	bic := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(opts.BIC), " ", ""))
	name := strings.TrimSpace(opts.Name)
	ref := strings.TrimSpace(opts.Reference)
	purpose := strings.TrimSpace(opts.Purpose)

	amountStr := ""
	if opts.Amount > 0 {
		amountStr = fmt.Sprintf("EUR%.2f", opts.Amount)
	}

	// EPC069-08 Format Specification:
	// Line 1: Service Tag: "BCD"
	// Line 2: Version: "002"
	// Line 3: Character set: "1" (UTF-8)
	// Line 4: Identification: "SCT" (SEPA Credit Transfer)
	// Line 5: BIC (optional in v002)
	// Line 6: Beneficiary Name
	// Line 7: Beneficiary IBAN
	// Line 8: Amount (e.g. "EUR12.50" or "")
	// Line 9: Purpose code (optional)
	// Line 10: Structured Reference (optional)
	// Line 11: Unstructured Remittance info (Verwendungszweck)
	// Line 12: Beneficiary to originator information (optional)
	var sb strings.Builder
	sb.WriteString("BCD\n")
	sb.WriteString("002\n")
	sb.WriteString("1\n")
	sb.WriteString("SCT\n")
	sb.WriteString(bic + "\n")
	sb.WriteString(name + "\n")
	sb.WriteString(iban + "\n")
	sb.WriteString(amountStr + "\n")
	sb.WriteString(purpose + "\n")
	sb.WriteString("\n") // Structured ref empty
	sb.WriteString(ref + "\n")

	return strings.TrimRight(sb.String(), "\n")
}
