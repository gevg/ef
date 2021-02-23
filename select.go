package ef

import "math/bits"

// hValue returns the higher bit value (bucket value) of the i-th number.
func (d *Dict) hValue(i int) uint {
	nOnes := i
	for j, bb := range d.b {
		ones := bits.OnesCount(bb)
		if i-ones < 0 {
			return uint(select64(bb, uint(i)) + j<<6 - nOnes)
		}
		i -= ones
	}
	return 0
}

// select1 ...
func (d *Dict) select1(i int) int {
	for j, bb := range d.b {
		ones := bits.OnesCount(bb)
		if i-ones < 0 {
			return select64(bb, uint(i)) + j<<6
		}
		i -= ones
	}
	return 0
}

// select0 ...
func (d *Dict) select0(i int) int {
	// i--
	// if 0 <= i {
	for j, bb := range d.b {
		zeros := bits.OnesCount(^bb)
		if i-zeros < 0 {
			return select64(^bb, uint(i)) + j<<6
		}
		i -= zeros
	}
	// }
	return 0
}

// select64 ...
func select64(x0, k uint) int {
	x1 := (x0 & 0x5555555555555555) + (x0 >> 1 & 0x5555555555555555)
	x2 := (x1 & 0x3333333333333333) + (x1 >> 2 & 0x3333333333333333)
	x3 := (x2 & 0x0F0F0F0F0F0F0F0F) + (x2 >> 4 & 0x0F0F0F0F0F0F0F0F)
	x4 := (x3 & 0x00FF00FF00FF00FF) + (x3 >> 8 & 0x00FF00FF00FF00FF)
	x5 := (x4 & 0x0000FFFF0000FFFF) + (x4 >> 16 & 0x0000FFFF0000FFFF)

	var ret uint = 0
	t := x5 >> ret & 0xFFFFFFFF
	if t <= k {
		k -= t
		ret += 32
	}
	t = x4 >> ret & 0xFFFF
	if t <= k {
		k -= t
		ret += 16
	}
	t = x3 >> ret & 0xFF
	if t <= k {
		k -= t
		ret += 8
	}
	t = x2 >> ret & 0xF
	if t <= k {
		k -= t
		ret += 4
	}
	t = x1 >> ret & 0x3
	if t <= k {
		k -= t
		ret += 2
	}
	t = x0 >> ret & 0x1
	if t <= k {
		k -= t
		ret++
	}
	return int(ret)
}

// NextLEQ ...
func (d *Dict) NextLEQ(value uint) (i int, previous uint) {

	return
}
