/*
 * MLC Barcode — barcode engine
 * Copyright (c) 2026 Michael Lechner
 *
 * This source code is licensed under the MIT license with attribution
 * clause (MIT with Attribution) found in the LICENSE file in the root
 * directory of this source tree.
 */
package barcodes

import (
	"errors"
	"fmt"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/aztec"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/code39"
	"github.com/boombuler/barcode/datamatrix"
	"github.com/boombuler/barcode/ean"
	"github.com/boombuler/barcode/qr"
	"github.com/boombuler/barcode/twooffive"
	"github.com/mlcmcp/mlc_barcode/internal/pdf417"
)

// BarcodeType represents the supported symbology types
type BarcodeType string

const (
	TypeQR         BarcodeType = "qr"
	TypeDataMatrix BarcodeType = "datamatrix"
	TypeCode128    BarcodeType = "code128"
	TypeCode39     BarcodeType = "code39"
	TypeEAN13      BarcodeType = "ean13"
	TypeEAN8       BarcodeType = "ean8"
	TypeUPCA       BarcodeType = "upca"
	TypeITF        BarcodeType = "itf"
	TypeAztec      BarcodeType = "aztec"
	TypePDF417     BarcodeType = "pdf417"
)

// AllTypes lists every supported symbology, in the order shown to users.
var AllTypes = []BarcodeType{
	TypeQR, TypeDataMatrix, TypeAztec, TypePDF417,
	TypeCode128, TypeCode39, TypeEAN13, TypeEAN8, TypeUPCA, TypeITF,
}

// BarcodeOptions holds configuration for barcode generation
type BarcodeOptions struct {
	Width           int
	Height          int
	ShowText        bool
	CustomText      string // If set, rendered under barcode instead of raw content
	FontSize        int    // Caption size in px; 0 = automatic
	ForegroundColor string // SVG color (e.g. "black", "#000")
	BackgroundColor string // SVG color (e.g. "white", "#fff")
	// For QR codes
	QRLevel qr.ErrorCorrectionLevel
	QRMode  qr.Encoding
	// For Code39
	IncludeChecksum bool
	FullASCIICode39 bool
	// For Aztec: minimum error correction in percent (0 = 23, the ISO default)
	AztecECCPercent int
	// For PDF417: security level 0-8 (-1 = chosen by data length)
	PDF417Security int
}

// DefaultOptions returns recommended default options for a barcode type
func DefaultOptions(btype BarcodeType) BarcodeOptions {
	opts := BarcodeOptions{
		QRLevel:         qr.M,
		QRMode:          qr.Auto,
		PDF417Security:  -1,
		ShowText:        false,
		ForegroundColor: "black",
		BackgroundColor: "white",
	}

	switch btype {
	case TypeQR, TypeDataMatrix, TypeAztec:
		opts.Width = 256
		opts.Height = 256
	case TypePDF417:
		opts.Width = 600
		opts.Height = 200
	case TypeEAN13, TypeUPCA:
		opts.Width = 450
		opts.Height = 150
	case TypeEAN8:
		opts.Width = 320
		opts.Height = 150
	default:
		opts.Width = 600
		opts.Height = 150
	}

	return opts
}

// Generate validates the input, encodes it and scales it to the requested
// size. Errors explain what is wrong and what to do (see inputcheck.go).
func Generate(btype BarcodeType, data string, opts BarcodeOptions) (barcode.Barcode, error) {
	data = strings.TrimSpace(data)
	if data == "" {
		return nil, errors.New("data must not be empty")
	}

	if IsRetail(btype) {
		check := CheckRetail(btype, data)
		if !check.Valid {
			return nil, fmt.Errorf("invalid %s: %s", btype, check.Error())
		}
		data = check.Code
	}
	if err := checkCharset(btype, data, opts); err != nil {
		return nil, err
	}

	bc, err := encode(btype, data, opts)
	if err != nil {
		if errors.Is(err, errUnsupportedType) {
			return nil, err
		}
		fits := func(s string) error { _, e := encode(btype, s, opts); return e }
		if capErr := capacityError(btype, data, fits); capErr != nil {
			return nil, capErr
		}
		return nil, encoderError(btype, err)
	}

	// Scale to the requested size, but never below the code's own size:
	// barcode.Scale cannot shrink, and a long Code 128 is wider than 600 px.
	if opts.Width > 0 && opts.Height > 0 {
		b := bc.Bounds()
		bc, err = barcode.Scale(bc, max(opts.Width, b.Dx()), max(opts.Height, b.Dy()))
		if err != nil {
			return nil, err
		}
	}
	return bc, nil
}

var errUnsupportedType = errors.New("unsupported barcode type")

// encode calls the symbology's encoder. data is already validated; for
// EAN/UPC it carries the check digit.
func encode(btype BarcodeType, data string, opts BarcodeOptions) (barcode.Barcode, error) {
	switch btype {
	case TypeQR:
		return qr.Encode(data, opts.QRLevel, opts.QRMode)
	case TypeDataMatrix:
		return datamatrix.Encode(data)
	case TypeAztec:
		ecc := opts.AztecECCPercent
		if ecc <= 0 {
			ecc = 23
		}
		return aztec.Encode([]byte(data), ecc, 0)
	case TypePDF417:
		return pdf417.Encode(data, pdf417Security(opts.PDF417Security, data))
	case TypeCode128:
		return code128.Encode(data)
	case TypeCode39:
		return code39.Encode(data, opts.IncludeChecksum, opts.FullASCIICode39)
	case TypeEAN13, TypeEAN8:
		return ean.Encode(data)
	case TypeUPCA:
		// A UPC-A is an EAN-13 with a leading 0. Passed on as is, the
		// encoder would read 12 digits as an EAN-13 without check digit
		// and encode a different number.
		return ean.Encode("0" + data)
	case TypeITF:
		return twooffive.Encode(data, true)
	}
	return nil, fmt.Errorf("%w: %s (supported: %s)", errUnsupportedType, btype, typeList())
}

// pdf417Security returns the requested level, or the minimum ISO 15438
// recommends for the amount of data.
func pdf417Security(level int, data string) byte {
	if level >= 0 && level <= 8 {
		return byte(level)
	}
	switch n := len(data); {
	case n <= 40:
		return 2
	case n <= 160:
		return 3
	case n <= 320:
		return 4
	default:
		return 5
	}
}

func typeList() string {
	names := make([]string, len(AllTypes))
	for i, t := range AllTypes {
		names[i] = string(t)
	}
	return strings.Join(names, ", ")
}

// CreateBarcodeToSVG is a helper function for backward compatibility
func CreateBarcodeToSVG(barcodetype string, data string, width, height int) (string, error) {
	btype := BarcodeType(strings.ToLower(barcodetype))
	opts := DefaultOptions(btype)
	if width > 0 {
		opts.Width = width
	}
	if height > 0 {
		opts.Height = height
	}
	return GenerateSVG(btype, data, opts)
}
