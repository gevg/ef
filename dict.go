package ef

import (
	"errors"
	"math/bits"
)

// Dict is a dictionary, built upon a bit sequence, compressed
// with the Elias-Fano compression algorithm.
type Dict struct {
	b          []uint // Elias-Fano representation of the bit sequence
	sizeLValue int    // bit length of a lValue
	sizeH      int    // bit length of the higher bits array H
	n          int    // number of entries in the dictionary
	lMask      uint   // mask to isolate lValue
	maxValue   uint   // universe of the Elias-Fano algorithm
	pValue     uint   // previously added value: helper value to check monotonicity
	cap        int    // number of entries that can be stored in the dictionary
}

// New constructs a dictionary with a capacity of n entries. The dictionary
// can accept positive numbers, smaller than or equal to maxValue. The values
// needed to be in monotonic, non-decreasing way.
func New(cap int, maxValue uint) (Dict, error) {
	if cap == 0 {
		return Dict{}, errors.New("ef: dictionary does not support an empty values list")
	}

	sizeLVal := max0(bits.Len(maxValue/uint(cap)) - 1) // Max(Floor(Log2(u/n)), 0)
	sizeH := cap + int(maxValue>>sizeLVal)

	return Dict{
		b:          make([]uint, (sizeH+cap*sizeLVal+63)>>6),
		sizeLValue: sizeLVal,
		sizeH:      sizeH,
		lMask:      1<<sizeLVal - 1,
		maxValue:   maxValue,
		cap:        cap,
	}, nil
}

// From constructs a dictionary from a list of values.
func From(values []uint) (*Dict, error) {
	l := len(values)
	if l == 0 {
		return nil, errors.New("ef: dictionary does not support an empty values list")
	}

	d, err := New(l, values[l-1])
	if err != nil {
		return nil, err
	}

	return d.build(values)
}

// Append appends a value to the dictionary and returns the key. If something
// goes wrong, Append returns -1. In that case, no value is added.
func (d *Dict) Append(value uint) int { // TODO: Return vervangen door error. Definieer een err var voor elk soort error.
	// check values
	if d.n != 0 && (d.cap <= d.n || value < d.pValue || d.maxValue < value) {
		return -1
	}
	d.pValue = value
	d.n++

	// higher bits processing
	hValue := value>>d.sizeLValue + uint(d.n-1)
	d.b[hValue>>6] |= 1 << (hValue & 63)
	if d.sizeLValue == 0 {
		return d.n - 1
	}

	// lower bits processing
	lValue := value & d.lMask
	offset := d.sizeH + (d.n-1)*d.sizeLValue
	d.b[offset>>6] |= lValue << (offset & 63)
	if msb := lValue >> (64 - offset&63); msb != 0 {
		d.b[offset>>6+1] = msb
	}
	return d.n - 1
}

// Cap returns the maximum number of entries in the dictionary.
func Cap(d *Dict) int {
	return d.cap
}

// Len returns the current number of entries in the dictionary.
// When the dictionary is built with From, Len and Cap always
// return the same result.
func Len(d *Dict) int {
	return d.n
}

// build encodes a monotone increasing array of positive integers
// into an Elias-Fano bit sequence. build will return an error
// when the values are not increasing monotonically.
func (d *Dict) build(values []uint) (*Dict, error) {
	var vmin uint
	offset := d.sizeH

	for i, v := range values {
		// check values
		if v < vmin {
			return nil, errors.New("ef: dictionary requires an array that increases monotonically")
		}
		vmin = v

		// higher bits processing
		hValue := v>>d.sizeLValue + uint(i)
		d.b[hValue>>6] |= 1 << (hValue & 63)
		if d.sizeLValue == 0 {
			continue
		}

		// lower bits processing
		lValue := v & d.lMask
		d.b[offset>>6] |= lValue << (offset & 63)
		if msb := lValue >> (64 - offset&63); msb != 0 {
			d.b[offset>>6+1] = msb
		}
		offset += d.sizeLValue
	}

	d.n = len(values)

	return d, nil
}

// max0 returns the maximum of x and 0.
func max0(x int) int {
	if x <= 0 {
		return 0
	}
	return x
}

// Value returns the k-th value in the dictionary.
func (d *Dict) Value(k int) (uint, bool) {
	if k < 0 || d.n <= k {
		return 0, false
	}
	return d.hValue(k)<<d.sizeLValue | d.lValue(k), true
}

// Value2 returns the k-th value in the dictionary.
func (d *Dict) Value2(k int) (uint, bool) {
	if k < 0 || d.n <= k {
		return 0, false
	}
	return uint((d.select1(k)-k)<<d.sizeLValue) | d.lValue(k), true
}

// RangeValues returns a contiguous range of values,
// from index start to end non-inclusive.
func (d *Dict) RangeValues(start, end int) []uint {
	// TODO
	return nil
}

// Values reads all numbers from the dictionary.
func (d *Dict) Values() []uint {
	values := make([]uint, d.n)
	var k int

	if d.sizeLValue == 0 {
		for p, b := range d.b {
			p64 := p << 6
			for b != 0 {
				values[k] = uint(p64 + bits.TrailingZeros(b) - k)
				b &= b - 1
				k++
			}
		}
		return values
	}

	lValFilter := d.sizeH
	for p, b := range d.b[:(d.sizeH+63)>>6] {
		b &= 1<<lValFilter - 1
		p64 := p << 6
		for b != 0 {
			hvalue := uint(p64 + bits.TrailingZeros(b) - k)
			values[k] = hvalue<<d.sizeLValue | d.lValue(k)
			b &= b - 1
			k++
		}
		lValFilter -= 64
	}
	return values
}

// hValue2 returns the higher bits (bucket value) of the k-th value.
func (d *Dict) hValue2(k int) uint {
	// hValue(k) == #zeros in [0, select1(k)] == select1(k) - k
	return uint(d.select1(k) - k)
}

// lValue returns the lower bits of the i-th number.
func (d *Dict) lValue(k int) uint { // TODO: Dit kan crashen wanneer d.sizeLValue == 0!!!!
	offset := d.sizeH + k*d.sizeLValue
	off63 := offset & 63
	val := d.b[offset>>6] >> off63 & d.lMask

	if off63+d.sizeLValue > 64 {
		val |= d.b[offset>>6+1] << (64 - off63) & d.lMask
	}

	return val
}

// NextGEQ returns the first value in the dictionary equal to or larger than
// the given value. NextGEQ fails when there exists no value larger than or
// equal to value. If so, it returns -1 as index.
func (d *Dict) NextGEQ(value uint) (i int, next uint) {
	if d.maxValue < value {
		return -1, d.maxValue
	}

	hValue := int(value >> d.sizeLValue)
	// pos denotes the end of bucket hValue-1
	pos := d.select0(hValue)
	// k is the number of integers whose high bits are less than hValue.
	// So, k is the key of the first value we need to compare to.
	k, pos63, pos64 := pos-hValue, pos&63, pos&^63 // k is wrong!!!
	buf := d.b[pos64>>6] &^ (1<<pos63 - 1)
	// This does not mean that the p+1-th value is larger than or equal to value!
	// The p+1-th value has either the same hValue or a larger one. In case of a
	// larger hValue, this value is the result we are seeking for. In case the
	// hValue is identical, we need to scan the bucket to be sure that our value
	// is larger. We can do this with select1, combined with search or scanning.
	// We can also do this with scanning without select1. Most probably the
	// fastest way for regular cases.

	// scan potential solutions
	lValue := value & d.lMask
	for {
		for buf == 0 {
			pos64 += 64
			buf = d.b[pos64>>6]
		}
		pos63 = bits.TrailingZeros(buf) // ==> potential solution with key k

		// check potential solution
		hVal := pos64 + pos63 - k
		if hValue < hVal || lValue <= d.lValue(k) {
			break
		}

		// next potential solution
		buf &= buf - 1
		k++
	}
	return k, uint(pos64+pos63-k)<<d.sizeLValue | d.lValue(k)
}
