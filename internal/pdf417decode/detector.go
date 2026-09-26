// Port of ZXing's com.google.zxing.pdf417.detector.Detector (Apache-2.0,
// see LICENSE and NOTICE in this directory): finds start and stop patterns
// and returns the eight vertices of each symbol. Java float arithmetic is
// kept as float32.

package pdf417decode

import (
	"math"

	"github.com/makiuchi-d/gozxing"
)

var (
	indexesStartPattern = []int{0, 4, 1, 5}
	indexesStopPattern  = []int{6, 2, 7, 3}
	startPattern        = []int{8, 1, 1, 1, 1, 1, 1, 3}    // 11111111 0 1 0 1 0 1 000
	stopPattern         = []int{7, 1, 1, 3, 1, 1, 1, 2, 1} // 1111111 0 1 000 1 0 1 00 1
	rotations           = []int{0, 180, 270, 90}
)

const (
	maxAvgVariance               float32 = 0.42
	maxIndividualVariance        float32 = 0.8
	maxStopPatternHeightVariance float32 = 0.5
	maxPixelDrift                        = 3
	maxPatternDrift                      = 5
	skippedRowCountMax                   = 25 // too low misses damaged start patterns, too high finds neighbours
	rowStep                              = 5  // >= 3 rows of >= 3 modules: about half of the 9 px minimum
	barcodeMinHeight                     = 10
)

type detectorResult struct {
	bits     *gozxing.BitMatrix
	points   [][]gozxing.ResultPoint
	rotation int
}

func detect(image *gozxing.BinaryBitmap, multiple bool) (*detectorResult, error) {
	original, err := image.GetBlackMatrix()
	if err != nil {
		return nil, err
	}
	for _, rotation := range rotations {
		m := applyRotation(original, rotation)
		if coords := detectIn(multiple, m); len(coords) > 0 {
			return &detectorResult{bits: m, points: coords, rotation: rotation}, nil
		}
	}
	return &detectorResult{bits: original, rotation: 0}, nil
}

func applyRotation(m *gozxing.BitMatrix, rotation int) *gozxing.BitMatrix {
	if rotation%360 == 0 {
		return m
	}
	c := cloneMatrix(m)
	switch rotation {
	case 90:
		c.Rotate90()
	case 180:
		c.Rotate180()
	case 270:
		c.Rotate90()
		c.Rotate180()
	}
	return c
}

func cloneMatrix(m *gozxing.BitMatrix) *gozxing.BitMatrix {
	c, _ := gozxing.NewBitMatrix(m.GetWidth(), m.GetHeight())
	for y := 0; y < m.GetHeight(); y++ {
		for x := 0; x < m.GetWidth(); x++ {
			if m.Get(x, y) {
				c.Set(x, y)
			}
		}
	}
	return c
}

func detectIn(multiple bool, m *gozxing.BitMatrix) [][]gozxing.ResultPoint {
	var coords [][]gozxing.ResultPoint
	row, column := 0, 0
	foundInRow := false
	for row < m.GetHeight() {
		vertices := findVertices(m, row, column)
		if vertices[0] == nil && vertices[3] == nil {
			if !foundInRow {
				break // nothing found: end of search
			}
			// Nothing at this column and row: retry from the first column,
			// slightly below the lowest barcode found so far.
			foundInRow = false
			column = 0
			for _, c := range coords {
				if c[1] != nil {
					row = max(row, int(c[1].GetY()))
				}
				if c[3] != nil {
					row = max(row, int(c[3].GetY()))
				}
			}
			row += rowStep
			continue
		}
		foundInRow = true
		coords = append(coords, vertices)
		if !multiple {
			break
		}
		// Without a right row indicator, continue after the start pattern.
		if vertices[2] != nil {
			column, row = int(vertices[2].GetX()), int(vertices[2].GetY())
		} else {
			column, row = int(vertices[4].GetX()), int(vertices[4].GetY())
		}
	}
	return coords
}

func findVertices(m *gozxing.BitMatrix, startRow, startColumn int) []gozxing.ResultPoint {
	height, width := m.GetHeight(), m.GetWidth()
	result := make([]gozxing.ResultPoint, 8)
	minHeight := barcodeMinHeight
	copyToResult(result, findRowsWithPattern(m, height, width, startRow, startColumn, minHeight, startPattern), indexesStartPattern)
	if result[4] != nil {
		startColumn, startRow = int(result[4].GetX()), int(result[4].GetY())
		if result[5] != nil {
			startPatternHeight := int(result[5].GetY()) - startRow
			minHeight = int(math.Max(float64(float32(startPatternHeight)*maxStopPatternHeightVariance), barcodeMinHeight))
		}
	}
	copyToResult(result, findRowsWithPattern(m, height, width, startRow, startColumn, minHeight, stopPattern), indexesStopPattern)
	return result
}

func copyToResult(result, tmp []gozxing.ResultPoint, destinationIndexes []int) {
	for i, d := range destinationIndexes {
		result[d] = tmp[i]
	}
}

func findRowsWithPattern(m *gozxing.BitMatrix, height, width, startRow, startColumn, minHeight int, pattern []int) []gozxing.ResultPoint {
	result := make([]gozxing.ResultPoint, 4)
	found := false
	minStartRow := startRow
	counters := make([]int, len(pattern))
	for ; startRow < height; startRow += rowStep {
		loc := findGuardPattern(m, startColumn, startRow, width, pattern, counters)
		if loc == nil {
			continue
		}
		for startRow > minStartRow+1 {
			startRow--
			prev := findGuardPattern(m, startColumn, startRow, width, pattern, counters)
			if prev == nil {
				startRow++
				break
			}
			loc = prev
		}
		result[0] = gozxing.NewResultPoint(float64(loc[0]), float64(startRow))
		result[1] = gozxing.NewResultPoint(float64(loc[1]), float64(startRow))
		found = true
		break
	}
	stopRow := startRow + 1
	if found { // last row of this symbol that contains the pattern
		skipped := 0
		prevLoc := []int{int(result[0].GetX()), int(result[1].GetX())}
		for ; stopRow < height; stopRow++ {
			loc := findGuardPattern(m, prevLoc[0], stopRow, width, pattern, counters)
			// Only a pattern at nearly the same position belongs to the same barcode.
			if loc != nil && abs(prevLoc[0]-loc[0]) < maxPatternDrift && abs(prevLoc[1]-loc[1]) < maxPatternDrift {
				prevLoc = loc
				skipped = 0
			} else if skipped > skippedRowCountMax {
				break
			} else {
				skipped++
			}
		}
		stopRow -= skipped + 1
		result[2] = gozxing.NewResultPoint(float64(prevLoc[0]), float64(stopRow))
		result[3] = gozxing.NewResultPoint(float64(prevLoc[1]), float64(stopRow))
	}
	if stopRow-startRow < minHeight {
		for i := range result {
			result[i] = nil
		}
	}
	return result
}

func findGuardPattern(m *gozxing.BitMatrix, column, row, width int, pattern, counters []int) []int {
	for i := range counters {
		counters[i] = 0
	}
	patternStart := column
	pixelDrift := 0
	// Black pixels left of the start: shift left, by at most maxPixelDrift.
	for pixel(m, patternStart, row) && patternStart > 0 && pixelDrift < maxPixelDrift {
		pixelDrift++
		patternStart--
	}
	x := patternStart
	counterPosition := 0
	patternLength := len(pattern)
	for isWhite := false; x < width; x++ {
		if pixel(m, x, row) != isWhite {
			counters[counterPosition]++
			continue
		}
		if counterPosition == patternLength-1 {
			if patternMatchVariance(counters, pattern) < maxAvgVariance {
				return []int{patternStart, x}
			}
			patternStart += counters[0] + counters[1]
			copy(counters, counters[2:counterPosition+1])
			counters[counterPosition-1] = 0
			counters[counterPosition] = 0
			counterPosition--
		} else {
			counterPosition++
		}
		counters[counterPosition] = 1
		isWhite = !isWhite
	}
	if counterPosition == patternLength-1 && patternMatchVariance(counters, pattern) < maxAvgVariance {
		return []int{patternStart, x - 1}
	}
	return nil
}

func patternMatchVariance(counters, pattern []int) float32 {
	total, patternLength := 0, 0
	for i := range counters {
		total += counters[i]
		patternLength += pattern[i]
	}
	if total < patternLength {
		return float32(math.Inf(1)) // less than one pixel per module: too small
	}
	unitBarWidth := float32(total) / float32(patternLength)
	maxVariance := maxIndividualVariance * unitBarWidth
	var totalVariance float32
	for i, counter := range counters {
		scaled := float32(pattern[i]) * unitBarWidth
		variance := float32(counter) - scaled
		if variance < 0 {
			variance = -variance
		}
		if variance > maxVariance {
			return float32(math.Inf(1))
		}
		totalVariance += variance
	}
	return totalVariance / float32(total)
}
