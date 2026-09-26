package barcodes

import "testing"

func TestCheckRetail(t *testing.T) {
	tests := []struct {
		name   string
		btype  BarcodeType
		data   string
		valid  bool
		reason string
		code   string
		added  bool
	}{
		{"ean13 complete", TypeEAN13, "4006381333931", true, "", "4006381333931", false},
		{"ean13 without check digit", TypeEAN13, "400638133393", true, "", "4006381333931", true},
		{"ean13 wrong check digit", TypeEAN13, "4006381333932", false, ReasonChecksum, "4006381333931", false},
		{"ean13 given 8 digits", TypeEAN13, "96385074", false, ReasonLength, "", false},
		{"ean13 letters", TypeEAN13, "40063813339X", false, ReasonNonDigit, "", false},
		{"ean8 complete", TypeEAN8, "96385074", true, "", "96385074", false},
		{"ean8 without check digit", TypeEAN8, "9638507", true, "", "96385074", true},
		{"upca complete", TypeUPCA, "036000291452", true, "", "036000291452", false},
		{"upca without check digit", TypeUPCA, "03600029145", true, "", "036000291452", true},
		{"upca wrong check digit", TypeUPCA, "036000291453", false, ReasonChecksum, "036000291452", false},
		{"upca given ean13", TypeUPCA, "4006381333931", false, ReasonLength, "", false},
		{"check digit 0", TypeEAN13, "978316148410", true, "", "9783161484100", true},
		{"other types pass", TypeCode128, "anything", true, "", "anything", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckRetail(tt.btype, tt.data)
			if got.Valid != tt.valid || got.Reason != tt.reason || got.Code != tt.code || got.CheckDigitAdded != tt.added {
				t.Errorf("CheckRetail(%s, %q) = %+v", tt.btype, tt.data, got)
			}
		})
	}
}

// A UPC-A must encode the same number as the EAN-13 with a leading 0 —
// not the 12 digits read as an EAN-13 without check digit.
func TestUPCAEncodesAsLeadingZeroEAN13(t *testing.T) {
	bc, err := Generate(TypeUPCA, "036000291452", BarcodeOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if got := bc.Content(); got != "0036000291452" {
		t.Errorf("UPC-A encoded as %s, want 0036000291452", got)
	}
}

func TestGenerateRejectsEmptyData(t *testing.T) {
	for _, bt := range []BarcodeType{TypeQR, TypeCode128, TypeEAN13} {
		if _, err := Generate(bt, "  ", BarcodeOptions{}); err == nil {
			t.Errorf("%s: empty data accepted", bt)
		}
	}
}

func TestGenerateRejectsWrongRetailInput(t *testing.T) {
	if _, err := Generate(TypeEAN13, "96385074", BarcodeOptions{}); err == nil {
		t.Error("EAN-13 accepted 8 digits (would silently become an EAN-8)")
	}
	if _, err := Generate(TypeEAN13, "4006381333932", BarcodeOptions{}); err == nil {
		t.Error("EAN-13 accepted a wrong check digit")
	}
}
