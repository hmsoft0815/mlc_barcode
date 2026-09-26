// Port of ZXing's DecodedBitStreamParser (com.google.zxing.pdf417.decoder)
// and ECIStringBuilder (com.google.zxing.common), Apache-2.0, see LICENSE
// and NOTICE in this directory. Macro PDF417 metadata is parsed for
// validation but not returned.

package pdf417decode

import (
	"math/big"
	"strconv"
	"strings"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common"
	"golang.org/x/text/encoding"
)

const (
	textCompactionModeLatch       = 900
	byteCompactionModeLatch       = 901
	numericCompactionModeLatch    = 902
	byteCompactionModeLatch6      = 924
	eciUserDefined                = 925
	eciGeneralPurpose             = 926
	eciCharset                    = 927
	beginMacroPDF417ControlBlock  = 928
	beginMacroPDF417OptionalField = 923
	macroPDF417Terminator         = 922
	readerInitialization          = 921
	modeShiftToByteCompactionMode = 913
	maxNumericCodewords           = 15
	numberOfSequenceCodewords     = 2
	pl, ll, as, ml, al, ps, pal   = 25, 27, 27, 28, 28, 29, 29
)

var (
	punctChars = []byte(";<>@[\\]_`~!\r\t,:\n-.$/\"|*()?{}'")
	mixedChars = []byte("0123456789&\r\t,:#-.$/+%*=^")
	exp900     = func() []*big.Int {
		e := make([]*big.Int, 16)
		e[0] = big.NewInt(1)
		for i := 1; i < len(e); i++ {
			e[i] = new(big.Int).Mul(e[i-1], big.NewInt(900))
		}
		return e
	}()
)

type textMode int

const (
	modeAlpha textMode = iota
	modeLower
	modeMixed
	modePunct
	modeAlphaShift
	modePunctShift
)

// eciStringBuilder collects bytes in the current character set (ISO-8859-1
// until an ECI switches it) and decodes them when the set changes.
type eciStringBuilder struct {
	current []byte
	result  strings.Builder
	charset encoding.Encoding // nil: ISO-8859-1
}

func (b *eciStringBuilder) appendByte(v byte)     { b.current = append(b.current, v) }
func (b *eciStringBuilder) appendString(s string) { b.current = append(b.current, s...) }

func (b *eciStringBuilder) appendECI(value int) error {
	b.flush()
	eci, err := common.GetCharacterSetECIByValue(value)
	if err != nil || eci == nil {
		return gozxing.NewFormatException()
	}
	b.charset = eci.GetCharset()
	return nil
}

func (b *eciStringBuilder) flush() {
	if len(b.current) == 0 {
		return
	}
	if b.charset == nil {
		for _, c := range b.current {
			b.result.WriteRune(rune(c)) // ISO-8859-1 maps 1:1 to Unicode
		}
	} else if s, err := b.charset.NewDecoder().Bytes(b.current); err == nil {
		b.result.Write(s)
	} else {
		b.result.Write(b.current)
	}
	b.current = b.current[:0]
}

func (b *eciStringBuilder) isEmpty() bool { return len(b.current) == 0 && b.result.Len() == 0 }

func (b *eciStringBuilder) String() string {
	b.flush()
	return b.result.String()
}

func decodeBitStream(codewords []int, ecLevel string) (*common.DecoderResult, error) {
	count := codewords[0]
	result := &eciStringBuilder{}
	codeIndex, err := textCompaction(codewords, 1, result)
	if err != nil {
		return nil, err
	}
	fileIDSeen := false
	for codeIndex < count {
		code := codewords[codeIndex]
		codeIndex++
		switch code {
		case textCompactionModeLatch:
			codeIndex, err = textCompaction(codewords, codeIndex, result)
		case byteCompactionModeLatch, byteCompactionModeLatch6:
			codeIndex, err = byteCompaction(code, codewords, codeIndex, result)
		case modeShiftToByteCompactionMode:
			if codeIndex >= count {
				return nil, gozxing.NewFormatException()
			}
			result.appendByte(byte(codewords[codeIndex]))
			codeIndex++
		case numericCompactionModeLatch:
			codeIndex, err = numericCompaction(codewords, codeIndex, result)
		case eciCharset:
			if codeIndex >= count {
				return nil, gozxing.NewFormatException()
			}
			err = result.appendECI(codewords[codeIndex])
			codeIndex++
		case eciGeneralPurpose:
			codeIndex += 2 // cannot use a generic ECI; skip its two characters
		case eciUserDefined:
			codeIndex++
		case beginMacroPDF417ControlBlock:
			codeIndex, err = decodeMacroBlock(codewords, codeIndex)
			fileIDSeen = err == nil
		case beginMacroPDF417OptionalField, macroPDF417Terminator:
			return nil, gozxing.NewFormatException() // only valid inside a macro block
		default:
			// Missing mode codeword: default to text compaction, as ZXing does.
			codeIndex, err = textCompaction(codewords, codeIndex-1, result)
		}
		if err != nil {
			return nil, err
		}
	}
	if result.isEmpty() && !fileIDSeen {
		return nil, gozxing.NewFormatException()
	}
	return common.NewDecoderResult(nil, result.String(), nil, ecLevel), nil
}

// decodeMacroBlock validates a Macro PDF417 control block and skips it.
func decodeMacroBlock(codewords []int, codeIndex int) (int, error) {
	maxLength := codewords[0]
	if codeIndex+numberOfSequenceCodewords > maxLength {
		return 0, gozxing.NewFormatException()
	}
	if _, err := decodeBase900toBase10(codewords[codeIndex:codeIndex+numberOfSequenceCodewords], numberOfSequenceCodewords); err != nil {
		return 0, err
	}
	codeIndex += numberOfSequenceCodewords
	fileIDStart := codeIndex
	for codeIndex < maxLength && codeIndex < len(codewords) &&
		codewords[codeIndex] != macroPDF417Terminator && codewords[codeIndex] != beginMacroPDF417OptionalField {
		codeIndex++
	}
	if codeIndex == fileIDStart {
		return 0, gozxing.NewFormatException() // at least one file ID codeword (Annex H.2)
	}
	for codeIndex < maxLength {
		var err error
		switch codewords[codeIndex] {
		case beginMacroPDF417OptionalField:
			codeIndex++
			if codeIndex >= maxLength {
				return 0, gozxing.NewFormatException()
			}
			field := &eciStringBuilder{}
			switch codewords[codeIndex] {
			case 0, 3, 4: // file name, sender, addressee
				codeIndex, err = textCompaction(codewords, codeIndex+1, field)
			case 1, 2, 5, 6: // segment count, time stamp, file size, checksum
				codeIndex, err = numericCompaction(codewords, codeIndex+1, field)
				if err == nil {
					if _, perr := strconv.ParseInt(field.String(), 10, 64); perr != nil {
						return 0, gozxing.NewFormatException()
					}
				}
			default:
				return 0, gozxing.NewFormatException()
			}
		case macroPDF417Terminator:
			codeIndex++
		default:
			return 0, gozxing.NewFormatException()
		}
		if err != nil {
			return 0, err
		}
	}
	return codeIndex, nil
}

func textCompaction(codewords []int, codeIndex int, result *eciStringBuilder) (int, error) {
	size := max((codewords[0]-codeIndex)*2, 0)
	textData := make([]int, size)
	byteData := make([]int, size)
	index := 0
	end := false
	subMode := modeAlpha
	for codeIndex < codewords[0] && !end {
		code := codewords[codeIndex]
		codeIndex++
		if code < textCompactionModeLatch {
			textData[index] = code / 30
			textData[index+1] = code % 30
			index += 2
			continue
		}
		switch code {
		case textCompactionModeLatch:
			textData[index] = textCompactionModeLatch
			index++
		case byteCompactionModeLatch, byteCompactionModeLatch6, numericCompactionModeLatch,
			beginMacroPDF417ControlBlock, beginMacroPDF417OptionalField, macroPDF417Terminator:
			codeIndex--
			end = true
		case modeShiftToByteCompactionMode:
			textData[index] = modeShiftToByteCompactionMode
			if codeIndex >= codewords[0] {
				return 0, gozxing.NewFormatException()
			}
			byteData[index] = codewords[codeIndex]
			codeIndex++
			index++
		case eciCharset:
			if codeIndex >= codewords[0] {
				return 0, gozxing.NewFormatException()
			}
			subMode = decodeTextCompaction(textData, byteData, index, result, subMode)
			if err := result.appendECI(codewords[codeIndex]); err != nil {
				return 0, err
			}
			codeIndex++
			size = max((codewords[0]-codeIndex)*2, 0)
			textData, byteData, index = make([]int, size), make([]int, size), 0
		}
	}
	decodeTextCompaction(textData, byteData, index, result, subMode)
	return codeIndex, nil
}

func decodeTextCompaction(textData, byteData []int, length int, result *eciStringBuilder, startMode textMode) textMode {
	subMode, priorToShift, latched := startMode, startMode, startMode
	for i := 0; i < length; i++ {
		c := textData[i]
		var ch byte
		switch subMode {
		case modeAlpha, modeLower:
			base := byte('A')
			if subMode == modeLower {
				base = 'a'
			}
			if c < 26 {
				ch = base + byte(c)
				break
			}
			switch {
			case c == 26:
				ch = ' '
			case c == ll && subMode == modeAlpha:
				subMode, latched = modeLower, modeLower
			case c == as && subMode == modeLower:
				priorToShift, subMode = subMode, modeAlphaShift
			case c == ml:
				subMode, latched = modeMixed, modeMixed
			case c == ps:
				priorToShift, subMode = subMode, modePunctShift
			case c == modeShiftToByteCompactionMode:
				result.appendByte(byte(byteData[i]))
			case c == textCompactionModeLatch:
				subMode, latched = modeAlpha, modeAlpha
			}
		case modeMixed:
			if c < pl {
				ch = mixedChars[c]
				break
			}
			switch c {
			case pl:
				subMode, latched = modePunct, modePunct
			case 26:
				ch = ' '
			case ll:
				subMode, latched = modeLower, modeLower
			case al, textCompactionModeLatch:
				subMode, latched = modeAlpha, modeAlpha
			case ps:
				priorToShift, subMode = subMode, modePunctShift
			case modeShiftToByteCompactionMode:
				result.appendByte(byte(byteData[i]))
			}
		case modePunct:
			if c < pal {
				ch = punctChars[c]
				break
			}
			switch c {
			case pal, textCompactionModeLatch:
				subMode, latched = modeAlpha, modeAlpha
			case modeShiftToByteCompactionMode:
				result.appendByte(byte(byteData[i]))
			}
		case modeAlphaShift:
			subMode = priorToShift
			if c < 26 {
				ch = 'A' + byte(c)
				break
			}
			switch c {
			case 26:
				ch = ' '
			case textCompactionModeLatch:
				subMode = modeAlpha
			}
		case modePunctShift:
			subMode = priorToShift
			if c < pal {
				ch = punctChars[c]
				break
			}
			switch c {
			case pal, textCompactionModeLatch:
				subMode = modeAlpha
			case modeShiftToByteCompactionMode:
				// PS before shift-to-byte is padding, see 5.4.2.4 of the spec
				result.appendByte(byte(byteData[i]))
			}
		}
		if ch != 0 {
			result.appendByte(ch)
		}
	}
	return latched
}

func byteCompaction(mode int, codewords []int, codeIndex int, result *eciStringBuilder) (int, error) {
	end := false
	for codeIndex < codewords[0] && !end {
		for codeIndex < codewords[0] && codewords[codeIndex] == eciCharset { // leading ECIs
			if codeIndex+1 >= codewords[0] {
				return 0, gozxing.NewFormatException()
			}
			codeIndex++
			if err := result.appendECI(codewords[codeIndex]); err != nil {
				return 0, err
			}
			codeIndex++
		}
		if codeIndex >= codewords[0] || codewords[codeIndex] >= textCompactionModeLatch {
			end = true
			continue
		}
		// decode one block of 5 codewords to 6 bytes
		var value int64
		count := 0
		for {
			value = 900*value + int64(codewords[codeIndex])
			codeIndex++
			count++
			if !(count < 5 && codeIndex < codewords[0] && codewords[codeIndex] < textCompactionModeLatch) {
				break
			}
		}
		if count == 5 && (mode == byteCompactionModeLatch6 ||
			codeIndex < codewords[0] && codewords[codeIndex] < textCompactionModeLatch) {
			for i := 0; i < 6; i++ {
				result.appendByte(byte(value >> (8 * (5 - i))))
			}
			continue
		}
		codeIndex -= count
		for codeIndex < codewords[0] && !end {
			code := codewords[codeIndex]
			codeIndex++
			switch {
			case code < textCompactionModeLatch:
				result.appendByte(byte(code))
			case code == eciCharset:
				if codeIndex >= codewords[0] {
					return 0, gozxing.NewFormatException()
				}
				if err := result.appendECI(codewords[codeIndex]); err != nil {
					return 0, err
				}
				codeIndex++
			default:
				codeIndex--
				end = true
			}
		}
	}
	return codeIndex, nil
}

func numericCompaction(codewords []int, codeIndex int, result *eciStringBuilder) (int, error) {
	count := 0
	end := false
	numeric := make([]int, maxNumericCodewords)
	for codeIndex < codewords[0] && !end {
		code := codewords[codeIndex]
		codeIndex++
		if codeIndex == codewords[0] {
			end = true
		}
		if code < textCompactionModeLatch {
			numeric[count] = code
			count++
		} else {
			switch code {
			case textCompactionModeLatch, byteCompactionModeLatch, byteCompactionModeLatch6,
				beginMacroPDF417ControlBlock, beginMacroPDF417OptionalField, macroPDF417Terminator, eciCharset:
				codeIndex--
				end = true
			}
		}
		if (count%maxNumericCodewords == 0 || code == numericCompactionModeLatch || end) && count > 0 {
			// 902 inside numeric mode ends the current group and starts a new one (5.4.4.2).
			s, err := decodeBase900toBase10(numeric, count)
			if err != nil {
				return 0, err
			}
			result.appendString(s)
			count = 0
		}
	}
	return codeIndex, nil
}

// decodeBase900toBase10 turns up to 15 base-900 codewords into digits; the
// leading 1 the encoder prepends is dropped.
func decodeBase900toBase10(codewords []int, count int) (string, error) {
	result := new(big.Int)
	for i := 0; i < count; i++ {
		result.Add(result, new(big.Int).Mul(exp900[count-i-1], big.NewInt(int64(codewords[i]))))
	}
	s := result.String()
	if s[0] != '1' {
		return "", gozxing.NewFormatException()
	}
	return s[1:], nil
}
