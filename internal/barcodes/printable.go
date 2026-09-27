package barcodes

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/mlcmcp/mlc_barcode/internal/qrformats"
)

// Codes that can be encoded but not used are refused, with what to do
// instead: a GiroCode a bank app rejects (wrong IBAN, fields beyond the EPC
// limits), a 2D code too dense to scan from paper, an image too small for
// its modules.

// maxPrintableModules is the largest 2D symbol still readable from a
// printed page with a phone: a QR code of version 27 (125×125 modules). At
// 5 cm a module is then 0.4 mm, about what phone cameras resolve.
const maxPrintableModules = 125

// minModulePixels is the smallest module an image may draw: below 1.5 px
// modules come out uneven and blur when printed. The default size of 2D
// codes (256 px) keeps about 2 px even at maxPrintableModules.
const minModulePixels = 1.5

// EPC069-12 field limits (SEPA credit transfer, version 002).
const (
	epcMaxName      = 70
	epcMaxReference = 140
	epcMaxAmount    = 999999999.99
)

// checkEPC refuses a GiroCode payload a banking app would reject.
func checkEPC(data string) error {
	if !strings.HasPrefix(data, "BCD\n") && !strings.HasPrefix(data, "BCD\r\n") {
		return nil
	}
	p := qrformats.Parse(data)
	if p.Kind != "epc" {
		return nil
	}
	f := p.Fields
	iban := f["iban"]
	if iban == "" || f["iban_valid"] != "true" {
		return inputError(ErrEPCIBAN,
			fmt.Sprintf("the IBAN %q is not valid (checksum or length) — a banking app would reject this GiroCode; check it against the invoice, do not change digits", iban),
			map[string]string{"iban": iban})
	}
	if n := len([]rune(f["name"])); n > epcMaxName {
		return inputError(ErrEPCName,
			fmt.Sprintf("the beneficiary name has %d characters, a GiroCode allows at most %d — shorten it", n, epcMaxName),
			map[string]string{"length": strconv.Itoa(n), "max": strconv.Itoa(epcMaxName)})
	}
	if n := len([]rune(f["reference"])); n > epcMaxReference {
		return inputError(ErrEPCReference,
			fmt.Sprintf("the reference has %d characters, a GiroCode allows at most %d — shorten it (invoice and customer number are enough)", n, epcMaxReference),
			map[string]string{"length": strconv.Itoa(n), "max": strconv.Itoa(epcMaxReference)})
	}
	if a := f["amount"]; a != "" {
		if v, err := strconv.ParseFloat(a, 64); err == nil && v > epcMaxAmount {
			return inputError(ErrEPCAmount,
				fmt.Sprintf("the amount %s EUR exceeds the GiroCode maximum of 999999999.99 EUR", a),
				map[string]string{"amount": a})
		}
	}
	return nil
}

// checkPrintable refuses a 2D symbol too dense to scan from paper, and an
// image that gives its modules less than minModulePixels. bc is the
// encoded symbol before scaling (one pixel per module).
func checkPrintable(btype BarcodeType, bc barcode.Barcode, opts BarcodeOptions) error {
	if btype != TypeQR && btype != TypeDataMatrix && btype != TypeAztec {
		return nil
	}
	n := bc.Bounds()
	modules := max(n.Dx(), n.Dy())
	if modules > maxPrintableModules {
		return inputError(ErrTooDense,
			fmt.Sprintf("this %s code needs %d×%d modules — too dense to scan from paper (at most %d); shorten the content or encode a link to it instead", btype, modules, modules, maxPrintableModules),
			map[string]string{"type": string(btype), "modules": strconv.Itoa(modules), "max": strconv.Itoa(maxPrintableModules)})
	}
	if opts.Width <= 0 || opts.Height <= 0 {
		return nil
	}
	qx, qy := 0, 0
	if !opts.NoQuietZone {
		qx, qy = quietModules(btype)
	}
	perModule := math.Min(float64(opts.Width)/float64(n.Dx()+2*qx), float64(opts.Height)/float64(n.Dy()+2*qy))
	if perModule < minModulePixels {
		need := int(math.Ceil(2 * float64(max(n.Dx()+2*qx, n.Dy()+2*qy))))
		return inputError(ErrTooSmall,
			fmt.Sprintf("%d×%d px is too small for this code (%d modules with quiet zone): a module would get %.1f px; use at least %d×%d px", opts.Width, opts.Height, max(n.Dx()+2*qx, n.Dy()+2*qy), perModule, need, need),
			map[string]string{"width": strconv.Itoa(opts.Width), "height": strconv.Itoa(opts.Height), "min": strconv.Itoa(need)})
	}
	return nil
}
