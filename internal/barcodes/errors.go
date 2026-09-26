package barcodes

import (
	"errors"
	"strconv"
)

// Error codes of InputError. They are a stable contract: the GUI translates
// by code and parameters, CLI and MCP print the English message.
const (
	ErrEmpty           = "empty"
	ErrRetailNonDigit  = "retail_non_digit"
	ErrRetailLength    = "retail_length"
	ErrRetailChecksum  = "retail_checksum"
	ErrCode39Charset   = "code39_charset"
	ErrASCIIOnly       = "ascii_only"
	ErrITFDigits       = "itf_digits"
	ErrITFEven         = "itf_even"
	ErrCapacity        = "capacity"
	ErrUnsupportedType = "unsupported_type"
	ErrEncoder         = "encoder"
	ErrNothingFound    = "nothing_found"
	ErrAztecCharset    = "aztec_charset"
	ErrImageFormat     = "image_format"
	ErrImageTooLarge   = "image_too_large"
)

// ErrorCodes lists every code; a test checks the GUI translates each one.
var ErrorCodes = []string{
	ErrEmpty, ErrRetailNonDigit, ErrRetailLength, ErrRetailChecksum, ErrCode39Charset,
	ErrASCIIOnly, ErrITFDigits, ErrITFEven, ErrCapacity, ErrUnsupportedType, ErrEncoder, ErrNothingFound, ErrAztecCharset,
	ErrImageFormat, ErrImageTooLarge,
}

// InputError is an input the symbology cannot encode. Error() is the
// English text; Code and Params let a UI word it in its own language.
type InputError struct {
	Code   string
	Params map[string]string
	msg    string
}

func (e *InputError) Error() string { return e.msg }

func inputError(code, msg string, params map[string]string) *InputError {
	if params == nil {
		params = map[string]string{}
	}
	return &InputError{Code: code, Params: params, msg: msg}
}

// AsInputError returns the InputError inside err, if any.
func AsInputError(err error) (*InputError, bool) {
	var ie *InputError
	ok := errors.As(err, &ie)
	return ie, ok
}

func itoa(n int) string { return strconv.Itoa(n) }
