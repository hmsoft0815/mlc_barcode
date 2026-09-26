// Port of ZXing's PDF417ScanningDecoder (com.google.zxing.pdf417.decoder,
// Apache-2.0, see LICENSE and NOTICE in this directory).

package pdf417decode

import (
	"strconv"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common"
)

const (
	codewordSkewSize = 2
	maxErrors        = 3
	maxECCodewords   = 512
)

func scanDecode(image *gozxing.BitMatrix, imageTopLeft, imageBottomLeft, imageTopRight, imageBottomRight gozxing.ResultPoint,
	minCodewordWidth, maxCodewordWidth int) (*common.DecoderResult, error) {

	box, err := newBoundingBox(image, imageTopLeft, imageBottomLeft, imageTopRight, imageBottomRight)
	if err != nil {
		return nil, err
	}
	var left, right *detectionResultColumn
	var result *detectionResult
	for firstPass := true; ; firstPass = false {
		if imageTopLeft != nil {
			left = rowIndicatorColumnAt(image, box, imageTopLeft, true, minCodewordWidth, maxCodewordWidth)
		}
		if imageTopRight != nil {
			right = rowIndicatorColumnAt(image, box, imageTopRight, false, minCodewordWidth, maxCodewordWidth)
		}
		if result, err = mergeIndicators(left, right); err != nil {
			return nil, err
		}
		if result == nil {
			return nil, gozxing.NewNotFoundException()
		}
		resultBox := result.boundingBox
		if firstPass && resultBox != nil && (resultBox.minY < box.minY || resultBox.maxY > box.maxY) {
			box = resultBox
		} else {
			break
		}
	}
	result.boundingBox = box
	maxBarcodeColumn := result.columnCount + 1
	result.columns[0] = left
	result.columns[maxBarcodeColumn] = right

	leftToRight := left != nil
	for n := 1; n <= maxBarcodeColumn; n++ {
		barcodeColumn := n
		if !leftToRight {
			barcodeColumn = maxBarcodeColumn - n
		}
		if result.columns[barcodeColumn] != nil {
			continue // the opposite row indicator column, already decoded
		}
		var column *detectionResultColumn
		if barcodeColumn == 0 || barcodeColumn == maxBarcodeColumn {
			column = newRowIndicatorColumn(box, barcodeColumn == 0)
		} else {
			column = newDetectionResultColumn(box)
		}
		result.columns[barcodeColumn] = column
		startColumn, previousStartColumn := -1, -1
		for imageRow := box.minY; imageRow <= box.maxY; imageRow++ {
			startColumn = startColumnFor(result, barcodeColumn, imageRow, leftToRight)
			if startColumn < 0 || startColumn > box.maxX {
				if previousStartColumn == -1 {
					continue
				}
				startColumn = previousStartColumn
			}
			cw := detectCodeword(image, box.minX, box.maxX, leftToRight, startColumn, imageRow, minCodewordWidth, maxCodewordWidth)
			if cw != nil {
				column.setCodeword(imageRow, cw)
				previousStartColumn = startColumn
				minCodewordWidth = min(minCodewordWidth, cw.width())
				maxCodewordWidth = max(maxCodewordWidth, cw.width())
			}
		}
	}
	return createDecoderResult(result)
}

func mergeIndicators(left, right *detectionResultColumn) (*detectionResult, error) {
	if left == nil && right == nil {
		return nil, nil
	}
	meta := barcodeMetadataOf(left, right)
	if meta == nil {
		return nil, nil
	}
	leftBox, err := adjustBoundingBox(left)
	if err != nil {
		return nil, err
	}
	rightBox, err := adjustBoundingBox(right)
	if err != nil {
		return nil, err
	}
	box, err := mergeBoundingBoxes(leftBox, rightBox)
	if err != nil {
		return nil, err
	}
	return newDetectionResult(meta, box), nil
}

func adjustBoundingBox(c *detectionResultColumn) (*boundingBox, error) {
	if c == nil {
		return nil, nil
	}
	rowHeights := c.rowHeights()
	if rowHeights == nil {
		return nil, nil
	}
	maxRowHeight := -1
	for _, h := range rowHeights {
		maxRowHeight = max(maxRowHeight, h)
	}
	missingStartRows := 0
	for _, h := range rowHeights {
		missingStartRows += maxRowHeight - h
		if h > 0 {
			break
		}
	}
	codewords := c.codewords
	for row := 0; missingStartRows > 0 && codewords[row] == nil; row++ {
		missingStartRows--
	}
	missingEndRows := 0
	for row := len(rowHeights) - 1; row >= 0; row-- {
		missingEndRows += maxRowHeight - rowHeights[row]
		if rowHeights[row] > 0 {
			break
		}
	}
	for row := len(codewords) - 1; missingEndRows > 0 && codewords[row] == nil; row-- {
		missingEndRows--
	}
	return c.boundingBox.addMissingRows(missingStartRows, missingEndRows, c.isLeft)
}

func barcodeMetadataOf(left, right *detectionResultColumn) *barcodeMetadata {
	var leftMeta *barcodeMetadata
	if left != nil {
		leftMeta = left.barcodeMetadata()
	}
	if leftMeta == nil {
		if right == nil {
			return nil
		}
		return right.barcodeMetadata()
	}
	var rightMeta *barcodeMetadata
	if right != nil {
		rightMeta = right.barcodeMetadata()
	}
	if rightMeta == nil {
		return leftMeta
	}
	if leftMeta.columnCount != rightMeta.columnCount &&
		leftMeta.errorCorrectionLevel != rightMeta.errorCorrectionLevel &&
		leftMeta.rowCount != rightMeta.rowCount {
		return nil
	}
	return leftMeta
}

func rowIndicatorColumnAt(image *gozxing.BitMatrix, box *boundingBox, startPoint gozxing.ResultPoint,
	leftToRight bool, minCodewordWidth, maxCodewordWidth int) *detectionResultColumn {
	column := newRowIndicatorColumn(box, leftToRight)
	for i := 0; i < 2; i++ {
		increment := 1
		if i == 1 {
			increment = -1
		}
		startColumn := int(startPoint.GetX())
		for imageRow := int(startPoint.GetY()); imageRow <= box.maxY && imageRow >= box.minY; imageRow += increment {
			cw := detectCodeword(image, 0, image.GetWidth(), leftToRight, startColumn, imageRow, minCodewordWidth, maxCodewordWidth)
			if cw != nil {
				column.setCodeword(imageRow, cw)
				if leftToRight {
					startColumn = cw.startX
				} else {
					startColumn = cw.endX
				}
			}
		}
	}
	return column
}

func adjustCodewordCount(result *detectionResult, matrix [][]*barcodeValue) error {
	cell := matrix[0][1]
	numberOfCodewords := cell.value()
	calculated := result.columnCount*result.meta.rowCount - numberOfECCodewords(result.meta.errorCorrectionLevel)
	if len(numberOfCodewords) == 0 {
		if calculated < 1 || calculated > maxCodewordsInBarcode {
			return gozxing.NewNotFoundException()
		}
		cell.setValue(calculated)
	} else if numberOfCodewords[0] != calculated && calculated >= 1 && calculated <= maxCodewordsInBarcode {
		// The value derived from the row indicators is more reliable.
		cell.setValue(calculated)
	}
	return nil
}

func createDecoderResult(result *detectionResult) (*common.DecoderResult, error) {
	matrix := createBarcodeMatrix(result)
	if err := adjustCodewordCount(result, matrix); err != nil {
		return nil, err
	}
	rows, cols := result.meta.rowCount, result.columnCount
	var erasures, ambiguousIndexes []int
	var ambiguousValues [][]int
	codewords := make([]int, rows*cols)
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			values := matrix[row][col+1].value()
			index := row*cols + col
			switch len(values) {
			case 0:
				erasures = append(erasures, index)
			case 1:
				codewords[index] = values[0]
			default:
				ambiguousIndexes = append(ambiguousIndexes, index)
				ambiguousValues = append(ambiguousValues, values)
			}
		}
	}
	return decoderResultFromAmbiguousValues(result.meta.errorCorrectionLevel, codewords, erasures, ambiguousIndexes, ambiguousValues)
}

func decoderResultFromAmbiguousValues(ecLevel int, codewords, erasures, ambiguousIndexes []int, ambiguousValues [][]int) (*common.DecoderResult, error) {
	counts := make([]int, len(ambiguousIndexes))
	for tries := 100; tries > 0; tries-- {
		for i, idx := range ambiguousIndexes {
			codewords[idx] = ambiguousValues[i][counts[i]]
		}
		res, err := decodeCodewords(codewords, ecLevel, erasures)
		if err == nil {
			return res, nil
		}
		if _, isChecksum := err.(gozxing.ChecksumException); !isChecksum {
			return nil, err
		}
		if len(counts) == 0 {
			return nil, gozxing.NewChecksumException()
		}
		for i := range counts {
			if counts[i] < len(ambiguousValues[i])-1 {
				counts[i]++
				break
			}
			counts[i] = 0
			if i == len(counts)-1 {
				return nil, gozxing.NewChecksumException()
			}
		}
	}
	return nil, gozxing.NewChecksumException()
}

func createBarcodeMatrix(result *detectionResult) [][]*barcodeValue {
	matrix := make([][]*barcodeValue, result.meta.rowCount)
	for row := range matrix {
		matrix[row] = make([]*barcodeValue, result.columnCount+2)
		for col := range matrix[row] {
			matrix[row][col] = newBarcodeValue()
		}
	}
	for col, column := range result.detectionResultColumns() {
		if column == nil {
			continue
		}
		for _, cw := range column.codewords {
			if cw == nil || cw.rowNumber < 0 || cw.rowNumber >= len(matrix) {
				continue
			}
			matrix[cw.rowNumber][col].setValue(cw.value)
		}
	}
	return matrix
}

func isValidBarcodeColumn(result *detectionResult, col int) bool {
	return col >= 0 && col <= result.columnCount+1 && result.columns[col] != nil
}

func startColumnFor(result *detectionResult, barcodeColumn, imageRow int, leftToRight bool) int {
	offset := 1
	if !leftToRight {
		offset = -1
	}
	endOrStart := func(cw *codeword, end bool) int {
		if end {
			return cw.endX
		}
		return cw.startX
	}
	if isValidBarcodeColumn(result, barcodeColumn-offset) {
		if cw := result.columns[barcodeColumn-offset].codeword(imageRow); cw != nil {
			return endOrStart(cw, leftToRight)
		}
	}
	if cw := result.columns[barcodeColumn].codewordNearby(imageRow); cw != nil {
		return endOrStart(cw, !leftToRight)
	}
	if isValidBarcodeColumn(result, barcodeColumn-offset) {
		if cw := result.columns[barcodeColumn-offset].codewordNearby(imageRow); cw != nil {
			return endOrStart(cw, leftToRight)
		}
	}
	skippedColumns := 0
	for isValidBarcodeColumn(result, barcodeColumn-offset) {
		barcodeColumn -= offset
		for _, prev := range result.columns[barcodeColumn].codewords {
			if prev != nil {
				return endOrStart(prev, leftToRight) + offset*skippedColumns*(prev.endX-prev.startX)
			}
		}
		skippedColumns++
	}
	if leftToRight {
		return result.boundingBox.minX
	}
	return result.boundingBox.maxX
}

func detectCodeword(image *gozxing.BitMatrix, minColumn, maxColumn int, leftToRight bool, startColumn, imageRow,
	minCodewordWidth, maxCodewordWidth int) *codeword {
	startColumn = adjustCodewordStartColumn(image, minColumn, maxColumn, leftToRight, startColumn, imageRow)
	moduleBitCount := moduleBitCountAt(image, minColumn, maxColumn, leftToRight, startColumn, imageRow)
	if moduleBitCount == nil {
		return nil
	}
	codewordBitCount := sum(moduleBitCount)
	var endColumn int
	if leftToRight {
		endColumn = startColumn + codewordBitCount
	} else {
		for i, j := 0, len(moduleBitCount)-1; i < j; i, j = i+1, j-1 {
			moduleBitCount[i], moduleBitCount[j] = moduleBitCount[j], moduleBitCount[i]
		}
		endColumn = startColumn
		startColumn = endColumn - codewordBitCount
	}
	if !checkCodewordSkew(codewordBitCount, minCodewordWidth, maxCodewordWidth) {
		return nil
	}
	decoded := decodedValue(moduleBitCount)
	cwValue := getCodeword(decoded)
	if cwValue == -1 {
		return nil
	}
	return newCodeword(startColumn, endColumn, codewordBucketNumber(bitCountForCodeword(decoded)), cwValue)
}

func moduleBitCountAt(image *gozxing.BitMatrix, minColumn, maxColumn int, leftToRight bool, startColumn, imageRow int) []int {
	imageColumn := startColumn
	moduleBitCount := make([]int, 8)
	moduleNumber := 0
	increment := 1
	if !leftToRight {
		increment = -1
	}
	previousPixelValue := leftToRight
	inRange := func() bool {
		if leftToRight {
			return imageColumn < maxColumn
		}
		return imageColumn >= minColumn
	}
	for inRange() && moduleNumber < len(moduleBitCount) {
		if pixel(image, imageColumn, imageRow) == previousPixelValue {
			moduleBitCount[moduleNumber]++
			imageColumn += increment
		} else {
			moduleNumber++
			previousPixelValue = !previousPixelValue
		}
	}
	edge := minColumn
	if leftToRight {
		edge = maxColumn
	}
	if moduleNumber == len(moduleBitCount) || (imageColumn == edge && moduleNumber == len(moduleBitCount)-1) {
		return moduleBitCount
	}
	return nil
}

// pixel reads the matrix like Java's BitMatrix.get; positions outside it
// count as white instead of panicking.
func pixel(image *gozxing.BitMatrix, x, y int) bool {
	if x < 0 || y < 0 || x >= image.GetWidth() || y >= image.GetHeight() {
		return false
	}
	return image.Get(x, y)
}

func numberOfECCodewords(ecLevel int) int { return 2 << ecLevel }

func adjustCodewordStartColumn(image *gozxing.BitMatrix, minColumn, maxColumn int, leftToRight bool, codewordStartColumn, imageRow int) int {
	corrected := codewordStartColumn
	increment := -1
	if !leftToRight {
		increment = 1
	}
	for i := 0; i < 2; i++ {
		for (leftToRight && corrected >= minColumn || !leftToRight && corrected < maxColumn) &&
			leftToRight == pixel(image, corrected, imageRow) {
			if abs(codewordStartColumn-corrected) > codewordSkewSize {
				return codewordStartColumn
			}
			corrected += increment
		}
		increment = -increment
		leftToRight = !leftToRight
	}
	return corrected
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func checkCodewordSkew(codewordSize, minCodewordWidth, maxCodewordWidth int) bool {
	return minCodewordWidth-codewordSkewSize <= codewordSize && codewordSize <= maxCodewordWidth+codewordSkewSize
}

func decodeCodewords(codewords []int, ecLevel int, erasures []int) (*common.DecoderResult, error) {
	if len(codewords) == 0 {
		return nil, gozxing.NewFormatException()
	}
	numEC := 1 << (ecLevel + 1)
	if len(erasures) > numEC/2+maxErrors || numEC < 0 || numEC > maxECCodewords {
		return nil, gozxing.NewChecksumException() // too many errors or corrupted EC level
	}
	corrected, err := correctErrors(codewords, numEC)
	if err != nil {
		return nil, err
	}
	if err := verifyCodewordCount(codewords, numEC); err != nil {
		return nil, err
	}
	res, err := decodeBitStream(codewords, strconv.Itoa(ecLevel))
	if err != nil {
		return nil, err
	}
	res.SetErrorsCorrected(corrected)
	res.SetErasures(len(erasures))
	return res, nil
}

func verifyCodewordCount(codewords []int, numEC int) error {
	if len(codewords) < 4 {
		// length descriptor, at least one data and two EC codewords
		return gozxing.NewFormatException()
	}
	n := codewords[0]
	if n > len(codewords) {
		return gozxing.NewFormatException()
	}
	if n == 0 {
		if numEC >= len(codewords) {
			return gozxing.NewFormatException()
		}
		codewords[0] = len(codewords) - numEC
	}
	return nil
}

func bitCountForCodeword(cw int) []int {
	result := make([]int, 8)
	previousValue := 0
	i := len(result) - 1
	for {
		if cw&0x1 != previousValue {
			previousValue = cw & 0x1
			i--
			if i < 0 {
				break
			}
		}
		result[i]++
		cw >>= 1
	}
	return result
}

func codewordBucketNumber(moduleBitCount []int) int {
	return (moduleBitCount[0] - moduleBitCount[2] + moduleBitCount[4] - moduleBitCount[6] + 9) % 9
}
