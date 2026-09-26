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
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/code39"
	"github.com/boombuler/barcode/datamatrix"
	"github.com/boombuler/barcode/ean"
	"github.com/boombuler/barcode/qr"
	"github.com/boombuler/barcode/twooffive"
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
)

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
}

// DefaultOptions returns recommended default options for a barcode type
func DefaultOptions(btype BarcodeType) BarcodeOptions {
	opts := BarcodeOptions{
		QRLevel:         qr.M,
		QRMode:          qr.Auto,
		ShowText:        false,
		ForegroundColor: "black",
		BackgroundColor: "white",
	}

	switch btype {
	case TypeQR, TypeDataMatrix:
		opts.Width = 256
		opts.Height = 256
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

// Generate generates a barcode object
func Generate(btype BarcodeType, data string, opts BarcodeOptions) (barcode.Barcode, error) {
	var bc barcode.Barcode
	var err error

	data = strings.TrimSpace(data)
	if data == "" {
		return nil, errors.New("data must not be empty")
	}

	switch btype {
	case TypeQR:
		bc, err = qr.Encode(data, opts.QRLevel, opts.QRMode)
	case TypeDataMatrix:
		bc, err = datamatrix.Encode(data)
	case TypeCode128:
		bc, err = code128.Encode(data)
	case TypeCode39:
		bc, err = code39.Encode(data, opts.IncludeChecksum, opts.FullASCIICode39)
	case TypeEAN13, TypeEAN8, TypeUPCA:
		check := CheckRetail(btype, data)
		if !check.Valid {
			return nil, fmt.Errorf("invalid %s: %s", btype, check.Error())
		}
		code := check.Code
		if btype == TypeUPCA {
			// A UPC-A is an EAN-13 with a leading 0. Passed on as is, the
			// encoder would read 12 digits as an EAN-13 without check digit
			// and encode a different number.
			code = "0" + code
		}
		bc, err = ean.Encode(code)
	case TypeITF:
		bc, err = twooffive.Encode(data, true)
	default:
		return nil, fmt.Errorf("unsupported barcode type: %s", btype)
	}

	if err != nil {
		return nil, err
	}

	// Scale to requested size
	if opts.Width > 0 && opts.Height > 0 {
		bc, err = barcode.Scale(bc, opts.Width, opts.Height)
		if err != nil {
			return nil, err
		}
	}

	return bc, nil
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
