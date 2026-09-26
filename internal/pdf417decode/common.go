// Port of ZXing's com.google.zxing.pdf417.PDF417Common (Apache-2.0, see
// LICENSE and NOTICE in this directory).

package pdf417decode

import "sort"

const (
	numberOfCodewords     = 929
	maxCodewordsInBarcode = numberOfCodewords - 1
	minRowsInBarcode      = 3
	maxRowsInBarcode      = 90
	modulesInCodeword     = 17
	modulesInStopPattern  = 18
	barsInModule          = 8
)

// getCodeword translates an encoded symbol to its codeword, or -1.
func getCodeword(symbol int) int {
	s := symbol & 0x3FFFF
	i := sort.SearchInts(symbolTable[:], s)
	if i >= len(symbolTable) || symbolTable[i] != s {
		return -1
	}
	return (codewordTable[i] - 1) % numberOfCodewords
}

func sum(values []int) int {
	s := 0
	for _, v := range values {
		s += v
	}
	return s
}
