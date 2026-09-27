package barcodes

import (
	"strings"
	"testing"

	"github.com/mlcmcp/mlc_barcode/internal/qrformats"
)

func wantCode(t *testing.T, err error, code string) {
	t.Helper()
	ie, ok := AsInputError(err)
	if !ok || ie.Code != code {
		t.Errorf("want error %q, got %v", code, err)
	}
}

// A GiroCode a banking app would reject is refused when it is made.
func TestGiroCodeChecked(t *testing.T) {
	epc := func(o qrformats.EPCOptions) string { return qrformats.FormatEPC(o) }
	good := qrformats.EPCOptions{Name: "Muster GmbH", IBAN: "DE89 3704 0044 0532 0130 00", Amount: 49.90, Reference: "Rechnung 1234"}
	if _, err := Generate(TypeQR, epc(good), DefaultOptions(TypeQR)); err != nil {
		t.Fatalf("valid GiroCode refused: %v", err)
	}
	bad := good
	bad.IBAN = "DE89370400440532013001" // last digit changed
	_, err := Generate(TypeQR, epc(bad), DefaultOptions(TypeQR))
	wantCode(t, err, ErrEPCIBAN)

	long := good
	long.Name = strings.Repeat("N", 71)
	_, err = Generate(TypeQR, epc(long), DefaultOptions(TypeQR))
	wantCode(t, err, ErrEPCName)

	long = good
	long.Reference = strings.Repeat("R", 141)
	_, err = Generate(TypeQR, epc(long), DefaultOptions(TypeQR))
	wantCode(t, err, ErrEPCReference)

	long = good
	long.Amount = 1e9
	_, err = Generate(TypeQR, epc(long), DefaultOptions(TypeQR))
	wantCode(t, err, ErrEPCAmount)

	// The same payload as DataMatrix is checked too.
	_, err = Generate(TypeDataMatrix, epc(bad), DefaultOptions(TypeDataMatrix))
	wantCode(t, err, ErrEPCIBAN)
}

// A 2D code beyond 125 modules cannot be scanned from paper; up to it the
// default size still gives every module about 2 px.
func TestPrintableSize(t *testing.T) {
	_, err := Generate(TypeQR, strings.Repeat("x", 1500), DefaultOptions(TypeQR))
	wantCode(t, err, ErrTooDense)

	// The densest QR still accepted (version 27, 125 modules) at the default size.
	fits := ""
	for n := 400; n < 1500; n += 5 {
		s := strings.Repeat("y", n)
		if _, err := Generate(TypeQR, s, DefaultOptions(TypeQR)); err != nil {
			break
		}
		fits = s
	}
	if len(fits) < 500 {
		t.Fatalf("only %d characters fit before the density limit", len(fits))
	}

	// An image too small for its modules.
	opts := DefaultOptions(TypeQR)
	opts.Width, opts.Height = 100, 100
	_, err = Generate(TypeQR, fits, opts)
	wantCode(t, err, ErrTooSmall)

	// Short content at a small size is fine.
	if _, err := Generate(TypeQR, "https://mlcgo.eu", opts); err != nil {
		t.Errorf("small QR refused: %v", err)
	}
}
