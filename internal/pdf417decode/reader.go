// Package pdf417decode reads PDF417 symbols. It is a Go port of the PDF417
// reader of ZXing (https://github.com/zxing/zxing, core/.../pdf417),
// licensed under the Apache License 2.0 — see LICENSE and NOTICE in this
// directory. Changes against the Java original: Go idioms (errors instead
// of exceptions), Macro PDF417 metadata is validated but not returned,
// out-of-range pixel reads count as white, and a null-column dereference in
// DetectionResult is guarded. The package plugs into
// github.com/makiuchi-d/gozxing as a gozxing.Reader.
package pdf417decode

import (
	"github.com/makiuchi-d/gozxing"
)

// Reader decodes PDF417 symbols; it implements gozxing.Reader.
type Reader struct{}

// NewReader returns a PDF417 reader.
func NewReader() *Reader { return &Reader{} }

func (r *Reader) DecodeWithoutHints(image *gozxing.BinaryBitmap) (*gozxing.Result, error) {
	return r.Decode(image, nil)
}

// Decode returns the first PDF417 symbol found.
func (r *Reader) Decode(image *gozxing.BinaryBitmap, _ map[gozxing.DecodeHintType]interface{}) (*gozxing.Result, error) {
	results, err := decodeAll(image, false)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, gozxing.NewNotFoundException()
	}
	return results[0], nil
}

// DecodeMultiple returns every PDF417 symbol found.
func (r *Reader) DecodeMultiple(image *gozxing.BinaryBitmap, _ map[gozxing.DecodeHintType]interface{}) ([]*gozxing.Result, error) {
	results, err := decodeAll(image, true)
	if err != nil || len(results) == 0 {
		return nil, gozxing.NewNotFoundException()
	}
	return results, nil
}

func (r *Reader) Reset() {}

func decodeAll(image *gozxing.BinaryBitmap, multiple bool) ([]*gozxing.Result, error) {
	det, err := detect(image, multiple)
	if err != nil {
		return nil, err
	}
	var results []*gozxing.Result
	for _, p := range det.points {
		dr, err := scanDecode(det.bits, p[4], p[5], p[6], p[7], minCodewordWidth(p), maxCodewordWidth(p))
		if err != nil {
			return nil, err
		}
		res := gozxing.NewResult(dr.GetText(), dr.GetRawBytes(), p, gozxing.BarcodeFormat_PDF_417)
		res.PutMetadata(gozxing.ResultMetadataType_ERROR_CORRECTION_LEVEL, dr.GetECLevel())
		res.PutMetadata(gozxing.ResultMetadataType_ORIENTATION, det.rotation)
		results = append(results, res)
	}
	return results, nil
}

func maxWidth(p1, p2 gozxing.ResultPoint) int32 {
	if p1 == nil || p2 == nil {
		return 0
	}
	return int32(absF(p1.GetX() - p2.GetX()))
}

func minWidth(p1, p2 gozxing.ResultPoint) int32 {
	if p1 == nil || p2 == nil {
		return 1<<31 - 1 // Integer.MAX_VALUE
	}
	return int32(absF(p1.GetX() - p2.GetX()))
}

func absF(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

// The width arithmetic stays in int32: ZXing multiplies Integer.MAX_VALUE by
// 17 and relies on Java's wrap-around.
func maxCodewordWidth(p []gozxing.ResultPoint) int {
	return int(max(
		max(maxWidth(p[0], p[4]), maxWidth(p[6], p[2])*modulesInCodeword/modulesInStopPattern),
		max(maxWidth(p[1], p[5]), maxWidth(p[7], p[3])*modulesInCodeword/modulesInStopPattern)))
}

func minCodewordWidth(p []gozxing.ResultPoint) int {
	return int(min(
		min(minWidth(p[0], p[4]), minWidth(p[6], p[2])*modulesInCodeword/modulesInStopPattern),
		min(minWidth(p[1], p[5]), minWidth(p[7], p[3])*modulesInCodeword/modulesInStopPattern)))
}
