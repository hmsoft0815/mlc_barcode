// Port of ZXing's com.google.zxing.pdf417.decoder.ec (Apache-2.0, see
// LICENSE and NOTICE in this directory): Reed-Solomon error correction over
// GF(929), the codeword field of PDF417.

package pdf417decode

import "github.com/makiuchi-d/gozxing"

type modulusGF struct {
	expTable []int
	logTable []int
	zero     *modulusPoly
	one      *modulusPoly
	modulus  int
}

var pdf417GF = newModulusGF(numberOfCodewords, 3)

func newModulusGF(modulus, generator int) *modulusGF {
	f := &modulusGF{modulus: modulus, expTable: make([]int, modulus), logTable: make([]int, modulus)}
	x := 1
	for i := 0; i < modulus; i++ {
		f.expTable[i] = x
		x = (x * generator) % modulus
	}
	for i := 0; i < modulus-1; i++ {
		f.logTable[f.expTable[i]] = i
	}
	f.zero = newModulusPoly(f, []int{0})
	f.one = newModulusPoly(f, []int{1})
	return f
}

func (f *modulusGF) buildMonomial(degree, coefficient int) *modulusPoly {
	if coefficient == 0 {
		return f.zero
	}
	c := make([]int, degree+1)
	c[0] = coefficient
	return newModulusPoly(f, c)
}

func (f *modulusGF) add(a, b int) int      { return (a + b) % f.modulus }
func (f *modulusGF) subtract(a, b int) int { return (f.modulus + a - b) % f.modulus }
func (f *modulusGF) exp(a int) int         { return f.expTable[a] }
func (f *modulusGF) log(a int) int         { return f.logTable[a] } // a != 0
func (f *modulusGF) inverse(a int) int     { return f.expTable[f.modulus-f.logTable[a]-1] }

func (f *modulusGF) multiply(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return f.expTable[(f.logTable[a]+f.logTable[b])%(f.modulus-1)]
}

type modulusPoly struct {
	field        *modulusGF
	coefficients []int
}

func newModulusPoly(field *modulusGF, coefficients []int) *modulusPoly {
	p := &modulusPoly{field: field, coefficients: coefficients}
	if n := len(coefficients); n > 1 && coefficients[0] == 0 {
		// Leading term must be non-zero for anything except the constant "0".
		first := 1
		for first < n && coefficients[first] == 0 {
			first++
		}
		if first == n {
			p.coefficients = []int{0}
		} else {
			p.coefficients = append([]int(nil), coefficients[first:]...)
		}
	}
	return p
}

func (p *modulusPoly) degree() int  { return len(p.coefficients) - 1 }
func (p *modulusPoly) isZero() bool { return p.coefficients[0] == 0 }

func (p *modulusPoly) coefficient(degree int) int {
	return p.coefficients[len(p.coefficients)-1-degree]
}

func (p *modulusPoly) evaluateAt(a int) int {
	if a == 0 {
		return p.coefficient(0)
	}
	if a == 1 {
		result := 0
		for _, c := range p.coefficients {
			result = p.field.add(result, c)
		}
		return result
	}
	result := p.coefficients[0]
	for _, c := range p.coefficients[1:] {
		result = p.field.add(p.field.multiply(a, result), c)
	}
	return result
}

func (p *modulusPoly) add(other *modulusPoly) *modulusPoly {
	if p.isZero() {
		return other
	}
	if other.isZero() {
		return p
	}
	small, large := p.coefficients, other.coefficients
	if len(small) > len(large) {
		small, large = large, small
	}
	sum := make([]int, len(large))
	diff := len(large) - len(small)
	copy(sum, large[:diff])
	for i := diff; i < len(large); i++ {
		sum[i] = p.field.add(small[i-diff], large[i])
	}
	return newModulusPoly(p.field, sum)
}

func (p *modulusPoly) subtract(other *modulusPoly) *modulusPoly {
	if other.isZero() {
		return p
	}
	return p.add(other.negative())
}

func (p *modulusPoly) multiply(other *modulusPoly) *modulusPoly {
	if p.isZero() || other.isZero() {
		return p.field.zero
	}
	a, b := p.coefficients, other.coefficients
	product := make([]int, len(a)+len(b)-1)
	for i, ac := range a {
		for j, bc := range b {
			product[i+j] = p.field.add(product[i+j], p.field.multiply(ac, bc))
		}
	}
	return newModulusPoly(p.field, product)
}

func (p *modulusPoly) negative() *modulusPoly {
	neg := make([]int, len(p.coefficients))
	for i, c := range p.coefficients {
		neg[i] = p.field.subtract(0, c)
	}
	return newModulusPoly(p.field, neg)
}

func (p *modulusPoly) multiplyScalar(scalar int) *modulusPoly {
	if scalar == 0 {
		return p.field.zero
	}
	if scalar == 1 {
		return p
	}
	product := make([]int, len(p.coefficients))
	for i, c := range p.coefficients {
		product[i] = p.field.multiply(c, scalar)
	}
	return newModulusPoly(p.field, product)
}

func (p *modulusPoly) multiplyByMonomial(degree, coefficient int) *modulusPoly {
	if coefficient == 0 {
		return p.field.zero
	}
	product := make([]int, len(p.coefficients)+degree)
	for i, c := range p.coefficients {
		product[i] = p.field.multiply(c, coefficient)
	}
	return newModulusPoly(p.field, product)
}

// correctErrors corrects received in place and returns the number of
// corrected codewords. Erasures are ignored, as in ZXing.
func correctErrors(received []int, numECCodewords int) (int, error) {
	field := pdf417GF
	if len(received) > field.modulus {
		return 0, gozxing.NewChecksumException()
	}
	poly := newModulusPoly(field, received)
	syndromes := make([]int, numECCodewords)
	hasError := false
	for i := numECCodewords; i > 0; i-- {
		eval := poly.evaluateAt(field.exp(i))
		syndromes[numECCodewords-i] = eval
		if eval != 0 {
			hasError = true
		}
	}
	if !hasError {
		return 0, nil
	}
	syndrome := newModulusPoly(field, syndromes)
	sigma, omega, err := runEuclideanAlgorithm(field.buildMonomial(numECCodewords, 1), syndrome, numECCodewords)
	if err != nil {
		return 0, err
	}
	locations, err := findErrorLocations(sigma)
	if err != nil {
		return 0, err
	}
	magnitudes := findErrorMagnitudes(omega, sigma, locations)
	for i, loc := range locations {
		position := len(received) - 1 - field.log(loc)
		if position < 0 {
			return 0, gozxing.NewChecksumException()
		}
		received[position] = field.subtract(received[position], magnitudes[i])
	}
	return len(locations), nil
}

func runEuclideanAlgorithm(a, b *modulusPoly, r int) (sigma, omega *modulusPoly, err error) {
	field := a.field
	if a.degree() < b.degree() {
		a, b = b, a
	}
	rLast, rCur := a, b
	tLast, t := field.zero, field.one
	for rCur.degree() >= r/2 {
		rLastLast, tLastLast := rLast, tLast
		rLast, tLast = rCur, t
		if rLast.isZero() {
			return nil, nil, gozxing.NewChecksumException()
		}
		rCur = rLastLast
		q := field.zero
		dltInverse := field.inverse(rLast.coefficient(rLast.degree()))
		for rCur.degree() >= rLast.degree() && !rCur.isZero() {
			degreeDiff := rCur.degree() - rLast.degree()
			scale := field.multiply(rCur.coefficient(rCur.degree()), dltInverse)
			q = q.add(field.buildMonomial(degreeDiff, scale))
			rCur = rCur.subtract(rLast.multiplyByMonomial(degreeDiff, scale))
		}
		t = q.multiply(tLast).subtract(tLastLast).negative()
	}
	sigmaTildeAtZero := t.coefficient(0)
	if sigmaTildeAtZero == 0 {
		return nil, nil, gozxing.NewChecksumException()
	}
	inverse := field.inverse(sigmaTildeAtZero)
	return t.multiplyScalar(inverse), rCur.multiplyScalar(inverse), nil
}

func findErrorLocations(errorLocator *modulusPoly) ([]int, error) {
	field := errorLocator.field
	numErrors := errorLocator.degree()
	result := make([]int, 0, numErrors)
	for i := 1; i < field.modulus && len(result) < numErrors; i++ {
		if errorLocator.evaluateAt(i) == 0 {
			result = append(result, field.inverse(i))
		}
	}
	if len(result) != numErrors {
		return nil, gozxing.NewChecksumException()
	}
	return result, nil
}

func findErrorMagnitudes(errorEvaluator, errorLocator *modulusPoly, errorLocations []int) []int {
	field := errorLocator.field
	d := errorLocator.degree()
	if d < 1 {
		return nil
	}
	derivative := make([]int, d)
	for i := 1; i <= d; i++ {
		derivative[d-i] = field.multiply(i, errorLocator.coefficient(i))
	}
	formalDerivative := newModulusPoly(field, derivative)
	result := make([]int, len(errorLocations))
	for i, loc := range errorLocations {
		xiInverse := field.inverse(loc)
		numerator := field.subtract(0, errorEvaluator.evaluateAt(xiInverse))
		denominator := field.inverse(formalDerivative.evaluateAt(xiInverse))
		result[i] = field.multiply(numerator, denominator)
	}
	return result
}
