package ef

import "math/bits"

// Iterator enables concurrent iteration over the dictionary.
type Iterator struct {
	d   *Dict // pointer to dictionary
	k   int   // last index read
	p64 int   // index over uint64 blocks of the H bit sequence.
	buf uint  // read buffer
}

// Iterator returns an iterator that iterates over the values in the dictionary.
// Iteration should not start before the static dictionary has been built.
// Iterators are required for concurrent access.
func (d *Dict) Iterator() Iterator {
	return Iterator{d: d, k: -1}
}

// Value returns the value at index k in the dictionary. It also sets the
// iterator state. Subsequent calls to Next will return the values with index
// k+1, k+2, ..., n-1. Iteration beyond n-1 will return 0, false.
func (it *Iterator) Value(k int) (uint, bool) {
	d := it.d
	if k < 0 || d.n <= k {
		return 0, false
	}
	it.k = k
	pos := d.select1(k)
	it.p64 = pos >> 6
	it.buf = d.b[it.p64] &^ (1<<(pos&63) - 1)

	return uint(pos-k)<<d.sizeLValue | d.lValue(k), true
}

// Next returns the next iteration value. If the iterator is not initialized
// through a call to Value, the first call to Next will return the first value
// in the dictionary.
func (it *Iterator) Next() (uint, bool) {
	if it.k == -1 {
		return it.Value(0)
	}

	d := it.d
	if it.k++; d.n <= it.k {
		it.k--
		return 0, false
	}

	it.buf &= it.buf - 1
	for it.buf == 0 {
		it.p64++
		it.buf = d.b[it.p64]
	}

	p63 := bits.TrailingZeros(it.buf)

	return uint(it.p64<<6+p63-it.k)<<d.sizeLValue | d.lValue(it.k), true
}
