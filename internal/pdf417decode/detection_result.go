// Port of ZXing's DetectionResultRowIndicatorColumn and DetectionResult
// (com.google.zxing.pdf417.decoder, Apache-2.0, see LICENSE and NOTICE in
// this directory).

package pdf417decode

// --- row indicator columns -------------------------------------------------

func (c *detectionResultColumn) setRowNumbers() {
	for _, cw := range c.codewords {
		if cw != nil {
			cw.setRowNumberAsRowIndicatorColumn()
		}
	}
}

func (c *detectionResultColumn) indicatorRange() (firstRow, lastRow int) {
	top, bottom := c.boundingBox.topRight, c.boundingBox.bottomRight
	if c.isLeft {
		top, bottom = c.boundingBox.topLeft, c.boundingBox.bottomLeft
	}
	return c.imageRowToCodewordIndex(int(top.GetY())), c.imageRowToCodewordIndex(int(bottom.GetY()))
}

func (c *detectionResultColumn) adjustCompleteIndicatorColumnRowNumbers(meta *barcodeMetadata) {
	codewords := c.codewords
	c.setRowNumbers()
	c.removeIncorrectCodewords(meta)
	firstRow, lastRow := c.indicatorRange()
	barcodeRow, maxRowHeight, currentRowHeight := -1, 1, 0
	for row := firstRow; row < lastRow; row++ {
		cw := codewords[row]
		if cw == nil {
			continue
		}
		rowDifference := cw.rowNumber - barcodeRow
		switch {
		case rowDifference == 0:
			currentRowHeight++
		case rowDifference == 1:
			maxRowHeight = max(maxRowHeight, currentRowHeight)
			currentRowHeight = 1
			barcodeRow = cw.rowNumber
		case rowDifference < 0 || cw.rowNumber >= meta.rowCount || rowDifference > row:
			codewords[row] = nil
		default:
			checkedRows := rowDifference
			if maxRowHeight > 2 {
				checkedRows = (maxRowHeight - 2) * rowDifference
			}
			closePreviousFound := checkedRows >= row
			for i := 1; i <= checkedRows && !closePreviousFound; i++ {
				closePreviousFound = codewords[row-i] != nil
			}
			if closePreviousFound {
				codewords[row] = nil
			} else {
				barcodeRow = cw.rowNumber
				currentRowHeight = 1
			}
		}
	}
}

func (c *detectionResultColumn) rowHeights() []int {
	meta := c.barcodeMetadata()
	if meta == nil {
		return nil
	}
	c.adjustIncompleteIndicatorColumnRowNumbers(meta)
	result := make([]int, meta.rowCount)
	for _, cw := range c.codewords {
		if cw != nil && cw.rowNumber < len(result) {
			result[cw.rowNumber]++
		}
	}
	return result
}

func (c *detectionResultColumn) adjustIncompleteIndicatorColumnRowNumbers(meta *barcodeMetadata) {
	firstRow, lastRow := c.indicatorRange()
	codewords := c.codewords
	barcodeRow, maxRowHeight, currentRowHeight := -1, 1, 0
	for row := firstRow; row < lastRow; row++ {
		cw := codewords[row]
		if cw == nil {
			continue
		}
		cw.setRowNumberAsRowIndicatorColumn()
		rowDifference := cw.rowNumber - barcodeRow
		switch {
		case rowDifference == 0:
			currentRowHeight++
		case rowDifference == 1:
			maxRowHeight = max(maxRowHeight, currentRowHeight)
			currentRowHeight = 1
			barcodeRow = cw.rowNumber
		case cw.rowNumber >= meta.rowCount:
			codewords[row] = nil
		default:
			barcodeRow = cw.rowNumber
			currentRowHeight = 1
		}
	}
	_ = maxRowHeight // kept as in ZXing
}

func (c *detectionResultColumn) barcodeMetadata() *barcodeMetadata {
	columnCount, rowCountUpper, rowCountLower, ecLevel := newBarcodeValue(), newBarcodeValue(), newBarcodeValue(), newBarcodeValue()
	for _, cw := range c.codewords {
		if cw == nil {
			continue
		}
		cw.setRowNumberAsRowIndicatorColumn()
		value := cw.value % 30
		row := cw.rowNumber
		if !c.isLeft {
			row += 2
		}
		switch row % 3 {
		case 0:
			rowCountUpper.setValue(value*3 + 1)
		case 1:
			ecLevel.setValue(value / 3)
			rowCountLower.setValue(value % 3)
		case 2:
			columnCount.setValue(value + 1)
		}
	}
	cc, ru, rl, ec := columnCount.value(), rowCountUpper.value(), rowCountLower.value(), ecLevel.value()
	if len(cc) == 0 || len(ru) == 0 || len(rl) == 0 || len(ec) == 0 || cc[0] < 1 ||
		ru[0]+rl[0] < minRowsInBarcode || ru[0]+rl[0] > maxRowsInBarcode {
		return nil
	}
	meta := newBarcodeMetadata(cc[0], ru[0], rl[0], ec[0])
	c.removeIncorrectCodewords(meta)
	return meta
}

func (c *detectionResultColumn) removeIncorrectCodewords(meta *barcodeMetadata) {
	for i, cw := range c.codewords {
		if cw == nil {
			continue
		}
		value := cw.value % 30
		row := cw.rowNumber
		if row > meta.rowCount {
			c.codewords[i] = nil
			continue
		}
		if !c.isLeft {
			row += 2
		}
		switch row % 3 {
		case 0:
			if value*3+1 != meta.rowCountUpperPart {
				c.codewords[i] = nil
			}
		case 1:
			if value/3 != meta.errorCorrectionLevel || value%3 != meta.rowCountLowerPart {
				c.codewords[i] = nil
			}
		case 2:
			if value+1 != meta.columnCount {
				c.codewords[i] = nil
			}
		}
	}
}

// --- detection result ------------------------------------------------------

const adjustRowNumberSkip = 2

type detectionResult struct {
	meta        *barcodeMetadata
	columns     []*detectionResultColumn
	boundingBox *boundingBox
	columnCount int
}

func newDetectionResult(meta *barcodeMetadata, box *boundingBox) *detectionResult {
	return &detectionResult{
		meta: meta, boundingBox: box, columnCount: meta.columnCount,
		columns: make([]*detectionResultColumn, meta.columnCount+2),
	}
}

func (d *detectionResult) detectionResultColumns() []*detectionResultColumn {
	d.adjustIndicatorColumnRowNumbers(d.columns[0])
	d.adjustIndicatorColumnRowNumbers(d.columns[d.columnCount+1])
	unadjusted := maxCodewordsInBarcode
	for {
		previous := unadjusted
		unadjusted = d.adjustRowNumbers()
		if !(unadjusted > 0 && unadjusted < previous) {
			break
		}
	}
	return d.columns
}

func (d *detectionResult) adjustIndicatorColumnRowNumbers(c *detectionResultColumn) {
	if c != nil {
		c.adjustCompleteIndicatorColumnRowNumbers(d.meta)
	}
}

func (d *detectionResult) adjustRowNumbers() int {
	unadjusted := d.adjustRowNumbersByRow()
	if unadjusted == 0 {
		return 0
	}
	for col := 1; col < d.columnCount+1; col++ {
		codewords := d.columns[col].codewords
		for row, cw := range codewords {
			if cw != nil && !cw.hasValidRowNumber() {
				d.adjustRowNumbersAt(col, row, codewords)
			}
		}
	}
	return unadjusted
}

func (d *detectionResult) adjustRowNumbersByRow() int {
	d.adjustRowNumbersFromBothRI()
	return d.adjustRowNumbersFromLRI() + d.adjustRowNumbersFromRRI()
}

func (d *detectionResult) adjustRowNumbersFromBothRI() {
	left, right := d.columns[0], d.columns[d.columnCount+1]
	if left == nil || right == nil {
		return
	}
	for row := range left.codewords {
		l, r := left.codewords[row], right.codewords[row]
		if l == nil || r == nil || l.rowNumber != r.rowNumber {
			continue
		}
		for col := 1; col <= d.columnCount; col++ {
			cw := d.columns[col].codewords[row]
			if cw == nil {
				continue
			}
			cw.rowNumber = l.rowNumber
			if !cw.hasValidRowNumber() {
				d.columns[col].codewords[row] = nil
			}
		}
	}
}

func (d *detectionResult) adjustRowNumbersFromRRI() int {
	right := d.columns[d.columnCount+1]
	if right == nil {
		return 0
	}
	unadjusted := 0
	for row, indicator := range right.codewords {
		if indicator == nil {
			continue
		}
		invalid := 0
		for col := d.columnCount + 1; col > 0 && invalid < adjustRowNumberSkip; col-- {
			if cw := d.columns[col].codewords[row]; cw != nil {
				invalid = adjustRowNumberIfValid(indicator.rowNumber, invalid, cw)
				if !cw.hasValidRowNumber() {
					unadjusted++
				}
			}
		}
	}
	return unadjusted
}

func (d *detectionResult) adjustRowNumbersFromLRI() int {
	left := d.columns[0]
	if left == nil {
		return 0
	}
	unadjusted := 0
	for row, indicator := range left.codewords {
		if indicator == nil {
			continue
		}
		invalid := 0
		for col := 1; col < d.columnCount+1 && invalid < adjustRowNumberSkip; col++ {
			if cw := d.columns[col].codewords[row]; cw != nil {
				invalid = adjustRowNumberIfValid(indicator.rowNumber, invalid, cw)
				if !cw.hasValidRowNumber() {
					unadjusted++
				}
			}
		}
	}
	return unadjusted
}

func adjustRowNumberIfValid(rowIndicatorRowNumber, invalidRowCounts int, cw *codeword) int {
	if cw.hasValidRowNumber() {
		return invalidRowCounts
	}
	if cw.isValidRowNumber(rowIndicatorRowNumber) {
		cw.rowNumber = rowIndicatorRowNumber
		return 0
	}
	return invalidRowCounts + 1
}

func (d *detectionResult) adjustRowNumbersAt(col, row int, codewords []*codeword) {
	cw := codewords[row]
	prevColumn := d.columns[col-1]
	if prevColumn == nil { // ZXing would dereference null here
		return
	}
	prev := prevColumn.codewords
	next := prev
	if d.columns[col+1] != nil {
		next = d.columns[col+1].codewords
	}
	var others [14]*codeword
	others[2] = prev[row]
	others[3] = next[row]
	if row > 0 {
		others[0], others[4], others[5] = codewords[row-1], prev[row-1], next[row-1]
	}
	if row > 1 {
		others[8], others[10], others[11] = codewords[row-2], prev[row-2], next[row-2]
	}
	if row < len(codewords)-1 {
		others[1], others[6], others[7] = codewords[row+1], prev[row+1], next[row+1]
	}
	if row < len(codewords)-2 {
		others[9], others[12], others[13] = codewords[row+2], prev[row+2], next[row+2]
	}
	for _, other := range others {
		if other != nil && other.hasValidRowNumber() && other.bucket == cw.bucket {
			cw.rowNumber = other.rowNumber
			return
		}
	}
}
