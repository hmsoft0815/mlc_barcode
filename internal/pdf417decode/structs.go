// Port of ZXing's Codeword, BarcodeValue, BarcodeMetadata, BoundingBox and
// DetectionResultColumn (com.google.zxing.pdf417.decoder, Apache-2.0, see
// LICENSE and NOTICE in this directory).

package pdf417decode

import (
	"math"
	"sort"

	"github.com/makiuchi-d/gozxing"
)

const barcodeRowUnknown = -1

type codeword struct {
	startX, endX, bucket, value int
	rowNumber                   int
}

func newCodeword(startX, endX, bucket, value int) *codeword {
	return &codeword{startX: startX, endX: endX, bucket: bucket, value: value, rowNumber: barcodeRowUnknown}
}

func (c *codeword) hasValidRowNumber() bool { return c.isValidRowNumber(c.rowNumber) }

func (c *codeword) isValidRowNumber(rowNumber int) bool {
	return rowNumber != barcodeRowUnknown && c.bucket == (rowNumber%3)*3
}

func (c *codeword) setRowNumberAsRowIndicatorColumn() {
	c.rowNumber = (c.value/30)*3 + c.bucket/3
}

func (c *codeword) width() int { return c.endX - c.startX }

// barcodeValue counts how often each value was read for one position.
type barcodeValue struct {
	values map[int]int
}

func newBarcodeValue() *barcodeValue { return &barcodeValue{values: map[int]int{}} }

func (b *barcodeValue) setValue(value int) { b.values[value]++ }

// value returns the values with the highest confidence, ascending (Java's
// HashMap iterates small Integer keys in that order).
func (b *barcodeValue) value() []int {
	keys := make([]int, 0, len(b.values))
	for k := range b.values {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	maxConfidence := -1
	var result []int
	for _, k := range keys {
		switch c := b.values[k]; {
		case c > maxConfidence:
			maxConfidence = c
			result = []int{k}
		case c == maxConfidence:
			result = append(result, k)
		}
	}
	return result
}

func (b *barcodeValue) confidence(value int) int { return b.values[value] }

type barcodeMetadata struct {
	columnCount, errorCorrectionLevel    int
	rowCountUpperPart, rowCountLowerPart int
	rowCount                             int
}

func newBarcodeMetadata(columnCount, rowCountUpperPart, rowCountLowerPart, errorCorrectionLevel int) *barcodeMetadata {
	return &barcodeMetadata{
		columnCount: columnCount, errorCorrectionLevel: errorCorrectionLevel,
		rowCountUpperPart: rowCountUpperPart, rowCountLowerPart: rowCountLowerPart,
		rowCount: rowCountUpperPart + rowCountLowerPart,
	}
}

type boundingBox struct {
	image                                      *gozxing.BitMatrix
	topLeft, bottomLeft, topRight, bottomRight gozxing.ResultPoint
	minX, maxX, minY, maxY                     int
}

func newBoundingBox(image *gozxing.BitMatrix, topLeft, bottomLeft, topRight, bottomRight gozxing.ResultPoint) (*boundingBox, error) {
	leftUnspecified := topLeft == nil || bottomLeft == nil
	rightUnspecified := topRight == nil || bottomRight == nil
	if leftUnspecified && rightUnspecified {
		return nil, gozxing.NewNotFoundException()
	}
	if leftUnspecified {
		topLeft = gozxing.NewResultPoint(0, topRight.GetY())
		bottomLeft = gozxing.NewResultPoint(0, bottomRight.GetY())
	} else if rightUnspecified {
		topRight = gozxing.NewResultPoint(float64(image.GetWidth()-1), topLeft.GetY())
		bottomRight = gozxing.NewResultPoint(float64(image.GetWidth()-1), bottomLeft.GetY())
	}
	return &boundingBox{
		image: image, topLeft: topLeft, bottomLeft: bottomLeft, topRight: topRight, bottomRight: bottomRight,
		minX: int(math.Min(topLeft.GetX(), bottomLeft.GetX())),
		maxX: int(math.Max(topRight.GetX(), bottomRight.GetX())),
		minY: int(math.Min(topLeft.GetY(), topRight.GetY())),
		maxY: int(math.Max(bottomLeft.GetY(), bottomRight.GetY())),
	}, nil
}

func (b *boundingBox) clone() *boundingBox {
	c := *b
	return &c
}

func mergeBoundingBoxes(left, right *boundingBox) (*boundingBox, error) {
	if left == nil {
		return right, nil
	}
	if right == nil {
		return left, nil
	}
	return newBoundingBox(left.image, left.topLeft, left.bottomLeft, right.topRight, right.bottomRight)
}

func (b *boundingBox) addMissingRows(missingStartRows, missingEndRows int, isLeft bool) (*boundingBox, error) {
	newTopLeft, newBottomLeft, newTopRight, newBottomRight := b.topLeft, b.bottomLeft, b.topRight, b.bottomRight
	if missingStartRows > 0 {
		top := b.topRight
		if isLeft {
			top = b.topLeft
		}
		newMinY := max(int(top.GetY())-missingStartRows, 0)
		newTop := gozxing.NewResultPoint(top.GetX(), float64(newMinY))
		if isLeft {
			newTopLeft = newTop
		} else {
			newTopRight = newTop
		}
	}
	if missingEndRows > 0 {
		bottom := b.bottomRight
		if isLeft {
			bottom = b.bottomLeft
		}
		newMaxY := min(int(bottom.GetY())+missingEndRows, b.image.GetHeight()-1)
		newBottom := gozxing.NewResultPoint(bottom.GetX(), float64(newMaxY))
		if isLeft {
			newBottomLeft = newBottom
		} else {
			newBottomRight = newBottom
		}
	}
	return newBoundingBox(b.image, newTopLeft, newBottomLeft, newTopRight, newBottomRight)
}

const maxNearbyDistance = 5

// detectionResultColumn holds the codewords found in one column, indexed
// by image row. rowIndicator is set for the left/right indicator columns.
type detectionResultColumn struct {
	boundingBox  *boundingBox
	codewords    []*codeword
	rowIndicator bool
	isLeft       bool
}

func newDetectionResultColumn(box *boundingBox) *detectionResultColumn {
	return &detectionResultColumn{boundingBox: box.clone(), codewords: make([]*codeword, box.maxY-box.minY+1)}
}

func newRowIndicatorColumn(box *boundingBox, isLeft bool) *detectionResultColumn {
	c := newDetectionResultColumn(box)
	c.rowIndicator, c.isLeft = true, isLeft
	return c
}

func (c *detectionResultColumn) codewordNearby(imageRow int) *codeword {
	if cw := c.codeword(imageRow); cw != nil {
		return cw
	}
	for i := 1; i < maxNearbyDistance; i++ {
		if near := c.imageRowToCodewordIndex(imageRow) - i; near >= 0 {
			if cw := c.codewords[near]; cw != nil {
				return cw
			}
		}
		if near := c.imageRowToCodewordIndex(imageRow) + i; near < len(c.codewords) {
			if cw := c.codewords[near]; cw != nil {
				return cw
			}
		}
	}
	return nil
}

func (c *detectionResultColumn) imageRowToCodewordIndex(imageRow int) int {
	return imageRow - c.boundingBox.minY
}

func (c *detectionResultColumn) setCodeword(imageRow int, cw *codeword) {
	c.codewords[c.imageRowToCodewordIndex(imageRow)] = cw
}

func (c *detectionResultColumn) codeword(imageRow int) *codeword {
	return c.codewords[c.imageRowToCodewordIndex(imageRow)]
}
