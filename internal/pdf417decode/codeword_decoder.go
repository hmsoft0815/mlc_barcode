// Port of ZXing's PDF417CodewordDecoder (com.google.zxing.pdf417.decoder,
// Apache-2.0, see LICENSE and NOTICE in this directory). Java float
// arithmetic is kept as float32 so borderline samples round the same way.

package pdf417decode

import "math"

var ratiosTable = func() [][barsInModule]float32 {
	table := make([][barsInModule]float32, len(symbolTable))
	for i, symbol := range symbolTable {
		current := symbol
		currentBit := current & 0x1
		for j := 0; j < barsInModule; j++ {
			var size float32
			for (current & 0x1) == currentBit {
				size++
				current >>= 1
			}
			currentBit = current & 0x1
			table[i][barsInModule-j-1] = size / modulesInCodeword
		}
	}
	return table
}()

func decodedValue(moduleBitCount []int) int {
	if v := decodedCodewordValue(sampleBitCounts(moduleBitCount)); v != -1 {
		return v
	}
	return closestDecodedValue(moduleBitCount)
}

func sampleBitCounts(moduleBitCount []int) []int {
	bitCountSum := float32(sum(moduleBitCount))
	result := make([]int, barsInModule)
	bitCountIndex, sumPreviousBits := 0, 0
	for i := 0; i < modulesInCodeword; i++ {
		sampleIndex := bitCountSum/(2*modulesInCodeword) + (float32(i)*bitCountSum)/modulesInCodeword
		if float32(sumPreviousBits+moduleBitCount[bitCountIndex]) <= sampleIndex {
			sumPreviousBits += moduleBitCount[bitCountIndex]
			bitCountIndex++
		}
		result[bitCountIndex]++
	}
	return result
}

func decodedCodewordValue(moduleBitCount []int) int {
	v := bitValue(moduleBitCount)
	if getCodeword(v) == -1 {
		return -1
	}
	return v
}

func bitValue(moduleBitCount []int) int {
	var result int64
	for i, count := range moduleBitCount {
		for bit := 0; bit < count; bit++ {
			result <<= 1
			if i%2 == 0 {
				result |= 1
			}
		}
	}
	return int(int32(result))
}

func closestDecodedValue(moduleBitCount []int) int {
	bitCountSum := sum(moduleBitCount)
	var ratios [barsInModule]float32
	if bitCountSum > 1 {
		for i := range ratios {
			ratios[i] = float32(moduleBitCount[i]) / float32(bitCountSum)
		}
	}
	bestMatchError := float32(math.MaxFloat32)
	bestMatch := -1
	for j, row := range ratiosTable {
		var e float32
		for k := 0; k < barsInModule; k++ {
			diff := row[k] - ratios[k]
			e += diff * diff
			if e >= bestMatchError {
				break
			}
		}
		if e < bestMatchError {
			bestMatchError = e
			bestMatch = symbolTable[j]
		}
	}
	return bestMatch
}
