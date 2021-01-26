package ef

import (
	"math/bits"

	"github.com/gevg/cssel" // (rs wordt geïmporteerd via cssel)
	"github.com/gevg/rs"
)

// EliasFano codec structure
type EliasFano struct {
	b         []uint64     // bit array containing higher bits, followed by lower bits
	idx       *cssel.Index // Select index on b (moet dit een pointer zijn???)
	m         int          // number of values in the data structure
	mask      uint32       // mask to isolate the lower Bits
	nLowBits  uint32       // number of lower bits per value
	lHighBits uint32       // bit length of the higher bits array
}

// From creates an EliasFano data structure from the input values.
func From(values []uint32) EliasFano {
	nLowBits, lHighBits, lBits := bitSizes(values)

	ef := EliasFano{
		b:         make([]uint64, (lBits+63)/64),
		m:         len(values),
		mask:      1<<nLowBits - 1,
		nLowBits:  nLowBits,
		lHighBits: lHighBits,
	}

	ef.build(values)

	return ef
}

// New creates an empty EliasFano data structure.
func New(universe int32, nValues int) EliasFano {
	// Je kan ook nog nLowBits, mask en lHighBits berekenen!
	return EliasFano{m: nValues}
}

// Add adds a value to the EliasFano data structure.
func (ef *EliasFano) Add(value int32) bool {
	//TODO
	return true
}

// universe returns the highest value.
func (ef *EliasFano) universe() uint32 {
	return ef.Access(ef.m - 1)
}

// bitSizes computes the bit size of the lower part of an Elias-Fano number,
// the lengths of the high bits and the all bits array.
func bitSizes(values []uint32) (uint32, uint32, uint32) {
	m := uint32(len(values))
	nLow := max0(bits.Len32(values[m-1]) - bits.Len32(m))
	if nLow > 0 && m > values[m-1]>>nLow {
		nLow--
	}

	// lHigh := m + values[m-1]>>nLow + 1 	// Folly telt de 1 niet! Dus geen 0 na de laatste 1-bit!
	lHigh := m + values[m-1]>>nLow

	return nLow, lHigh, lHigh + m*nLow
}

// max0 returns the maximum of parameter a and 0.
func max0(a int) uint32 {
	if a > 0 {
		return uint32(a)
	}
	return 0
}

// In case the sequence x0, x1, ..., xn-1 to be represented is !strictly!
// monotone, it is possible to reduce the space usage by storing the
// sequence xi - i using the upper bound u - n. Retrieval happens in the same
// way — one just has to adjust the retrieved value for the i-th element by
// adding i.

// build encodes a monotone increasing array of positive integers into the
// Elias-Fano data structure. build will fail when the array is not increasing
// monotonically.
func (ef *EliasFano) build(values []uint32) { // Als values []uint64 is, dan moet ik minder casten!!!
	offset := ef.lHighBits

	for i, v := range values {
		// higher bits
		posHigh := v>>ef.nLowBits + uint32(i)
		ef.b[posHigh/64] |= 1 << (posHigh & 63)

		// lower bits
		lowBits := uint64(v & ef.mask)
		ef.b[offset/64] |= lowBits << (offset & 63)
		if bs := lowBits >> (64 - offset&63); bs != 0 { // Dit kan sneller met unsafe!!!
			ef.b[offset/64+1] = bs
		}
		offset += ef.nLowBits
	}
}

// Access returns the i-th value.
// access(i, selector, l) = (selector.select(i) - i) L[il, (i+1)l);
func (ef *EliasFano) Access(i int) uint32 { // int32 van maken zodat ik -1 kan teruggeven wanneer de index i niet bestaat.
	if i < 0 || i >= ef.m {
		return 0 // Nog foutmelding toevoegen! -1 indien mogelijk!
	}
	return ef.higherBits(i)<<ef.nLowBits | ef.lowerBits(i)
}

// AccessRange returns a contiguous range of values, starting at index start
// until index end non-inclusive.
func (ef *EliasFano) AccessRange(start, end int) []uint64 {
	// TODO
	return nil
}

// We kunnen een batchSelect implementeren die de select berekent van een
// gesorteerde lijst ranks. Dit is nuttig, aangezien we dikwijls meerdere
// selects moeten kennen, bvb om een ondergrens en een bovengrens te bepalen.

// NextLEQ returns the index of the largest number in EliasFano lesser than or
// equal to value.
// Ik kan misschien ook het verschil teruggeven als een tweede parameter!!!
func (ef *EliasFano) NextLEQ(value uint32) int {
	if value >= ef.Access(ef.m-1) {
		return ef.m - 1
	}

	bucket := value >> ef.nLowBits
	_ = rs.Select0(ef.b, bucket) - int(bucket) // p is bitpositie van de 0 die de start van de juiste bucket aangeeft.
	// Ik kan uit p het aantal verbruikte enen berekenen en dus kan ik ook de index berekenen van de lowbits.
	// Als het eerste cijfer in de bucket direct te groot is, dan vorige index meegeven.
	// Het kan zelfs zijn dat de bucket geen cijfers bevat.
	// Als alle cijfers in de bucket te klein zijn, dan moet je de index van het laatste cijfer in de bucket teruggeven.
	// De volgende bucket is immers steeds te groot (als hij bestaat).

	var zeros, n, sum int
	for _, b := range ef.b {
		sum = bits.OnesCount64(^b)
		if zeros+sum >= int(bucket) {
			// n := 64*i - zeros + pt3.Select(^b, uint8(bucket)-uint8(zeros)) - (bucket - zeros)
			// n := 64*i - zeros + rs.Select64(^b, bucket-zeros) - (bucket - zeros)
			break
		}
		zeros += sum
	}

	for i := 64*n - zeros; i < ef.m; i++ {
		if ef.Access(i) > value {
			return i - 1
		}
	}
	return -1
}

// NextGEQ returns the first value in the Elias-Fano data structure, larger than or
// equal to value.
// func (ef *EliasFano) NextGEQ(value uint32) uint32 {
// 	y := int(value >> ef.nLowBits)
// 	value &= ef.mask
// 	for i := range ef.b {
// 		zeros := bits.OnesCount64(^ef.b[i])
// 		if y-zeros <= 0 {
// 			pos := i<<6 - y + int(pt3.Select(^ef.b[i], uint8(y))) - 1 // Wat als y == 0???
// 			for ef.lowerBits(pos) < value {
// 				pos++
// 			}
// 			return ef.Access(pos)
// 		}
// 		y -= zeros
// 	}
// 	return 0
// }

// NextGEQ returns the first value in the Elias-Fano data structure, equal or
// larger than x.
// TODO: Proberen de lus met Select te vervangen door b&-b of b<<r
// func (ef *EliasFano) NextGEQ(x uint32) uint32 {
// 	y := int(x >> ef.nLowBits)
// 	for i := range ef.b {
// 		zeros := bits.OnesCount64(^ef.b[i])
// 		if y-zeros <= 0 {
// 			idx := i<<6 - y + int(pt3.Select(^ef.b[i], uint8(y))) - 1
// 			lowBits := ef.lowerBits(idx)
// 			for lowBits < x&ef.mask {
// 				idx++
// 				lowBits = ef.lowerBits(idx)
// 			}
// 			return x&^ef.mask | lowBits
// 		}
// 		y -= zeros
// 	}
// 	return 0
// }

// // NextGEQ returns the first value in the Elias-Fano data structure, equal or
// // larger than x. NextGEQ fails when there exists no value, larger than or
// // equal to x.
// // TODO: add a comparison to elems[m-1]? I do not store elems[m-1]. I could use
// // ef.lHighBits although not very efficient.
// func (ef *EliasFano) NextGEQ(x uint32) uint32 { // NextGEQ(x) zit niet altijd in dezelfde bucket als x!
// 	bkt := int(x >> ef.nLowBits)
// 	for i := range ef.b {
// 		zeros := bits.OnesCount64(^ef.b[i])
// 		if bkt-zeros <= 0 {
// 			idx := i<<6 - bkt + int(pt3.Select(^ef.b[i], uint8(bkt))) - 1
// 			for {
// 				lowBits := ef.lowerBits(idx)
// 				if lowBits >= x&ef.mask {
// 					return x&^ef.mask | lowBits
// 				}
// 				idx++
// 			}
// 		}
// 		bkt -= zeros
// 	}
// 	return 0
// }

// // NextGEQ returns the first value in the Elias-Fano data structure, equal or
// // larger than x. x does not need to be an element of the data structure.
// func (ef *EliasFano) NextGEQ(x uint32) uint32 { // uint32 moet int32 worden aangezien psi int32 is!
// 	if x > ef.Max() { // Kan nog sneller door eerst highBits te vergelijken en da evtl de lowBits!
// 		return 0 // Nog error boodschap toevoegen
// 	}

// 	bucket := int(x >> ef.nLowBits)
// 	nZeros := bucket
// 	for i, b := range ef.b {
// 		b = ^b
// 		zeros := bits.OnesCount64(b)
// 		if bucket-zeros <= 0 {
// 			p := pt3.Select(b, uint8(bucket))
// 			idx := i<<6 + int(p) - nZeros
// 			b &^= (1 << (p + 1)) - 1

// 			// if c := b & (b - 1); c != 0 {
// 			// 	tzc := bits.TrailingZeros64(c)
// 			// tzb := bits.TrailingZeros64(b)
// 			if b != 0 {
// 				tzb := bits.TrailingZeros64(b)
// 				j := 0
// 				for j < tzb-int(p) {
// 					if x&ef.mask <= ef.lowerBits(idx+j) { // Dit kan nog sneller!
// 						return ef.Access(idx + j) // Dit kan nog sneller! Moet dit wel idx zijn???
// 					}
// 					j++
// 				}
// 				if uint32(tzb) < ef.lHighBits-uint32(i)<<6 {
// 					return ef.Access(idx + j) // Dit kan nog sneller!
// 				}
// 				// Volgende ef.b[] opvragen
// 			}
// 			// b ^= b & -b
// 			// Search in lowBits, belonging to bucket 'bkt'
// 			// for j := idx0; j < idx1; j++ {
// 			// 	if low := ef.lowerBits(j); low >= uint32(nZeros) {
// 			// 		return uint32(nZeros)<<ef.nLowBits | low
// 			// 	}
// 			// }
// 			// Als we hier gekomen zijn, is er geen groter getal in de huidige bucket.
// 			// We moeten dus het eerste getal in de volgende bucket returnen, indien
// 			// dit bestaat.
// 			// return ...

// 			// Als bkt-zeros < 0 of als b de laatste b is in ef.b, dan weten we
// 			// zeker dat de gezochte y in deze b zit. Als bkt-zeros == 0, dan
// 			// is het niet zeker dat de gezochte y in deze b zit en dan moet in
// 			// het ongunstige geval een extra b('s) gefetched worden!
// 			if bucket-zeros < 0 || len(ef.b) == i+1 { // bkt-zeros < 0 is enkel juist als achter de laatste bucket geen 0 meer komt!!!
// 				for { // Wat als x te groot is? Komen we dan in een oneindige lus?
// 					lowBits := ef.lowerBits(idx)
// 					if lowBits >= x&ef.mask {
// 						return x&^ef.mask | lowBits
// 					}
// 					idx++
// 				}
// 			}
// 			// Hier de andere gevallen.
// 			// We moeten controleren dat we tijdens de vergelijking van de
// 			// lowBits niet in een nieuwe bucket belanden zonder dat we het
// 			// weten. Dit gebeurt best door de volgende nul te zoeken?

// 			b ^= b & -b
// 		}
// 		bucket -= zeros
// 	}

// 	return 0 // hier zou ik een error moeten returnen!
// }

// NextGEQ returns the first value in the Elias-Fano data structure, equal or
// larger than x. x does not need to be an element of the data structure.
// func (ef *EliasFano) NextGEQ(x uint32) uint32 { // uint32 moet int32 worden aangezien psi int32 is!
// // Kan sneller door eerst highBits te vergelijken en dan evtl de lowBits!
// if x > ef.Max() {
// 	// Nog error boodschap toevoegen
// 	return 0
// }

// var start uint64
// b := ef.b[0]
// bucket := uint64(x >> ef.nLowBits)
// if bucket != 0 {
// 	start = pt3.Select(^b, uint8(bucket-1)) - (bucket - 1)
// 	b &^= 1<<start - 1
// }

// // check higherBits of first number that is in bucket or higher
// if high := pt3.Select(b, uint8(start)) - start; high > bucket {
// 	return uint32(high)<<ef.nLowBits | ef.lowerBits(int(start))
// }

// // check the x bucket
// // for {

// // }
// return 0

// nZeros := bucket
// for i, b := range ef.b {
// 	b = ^b
// 	zeros := bits.OnesCount64(b)
// 	if bucket-zeros < 0 {
// 		p := pt3.Select(b, uint8(bucket))
// 		idx := i<<6 + int(p) - nZeros
// 		b &^= (1 << (p + 1)) - 1

// 		// if c := b & (b - 1); c != 0 {
// 		// 	tzc := bits.TrailingZeros64(c)
// 		// tzb := bits.TrailingZeros64(b)
// 		if b != 0 {
// 			tzb := bits.TrailingZeros64(b)
// 			j := 0
// 			for j < tzb-int(p) {
// 				if x&ef.mask <= ef.lowerBits(idx+j) { // Dit kan nog sneller!
// 					return ef.Access(idx + j) // Dit kan nog sneller! Moet dit wel idx zijn???
// 				}
// 				j++
// 			}
// 			if uint32(tzb) < ef.lHighBits-uint32(i)<<6 {
// 				return ef.Access(idx + j) // Dit kan nog sneller!
// 			}
// 			// Volgende ef.b[] opvragen
// 		}
// 		// b ^= b & -b
// 		// Search in lowBits, belonging to bucket 'bkt'
// 		// for j := idx0; j < idx1; j++ {
// 		// 	if low := ef.lowerBits(j); low >= uint32(nZeros) {
// 		// 		return uint32(nZeros)<<ef.nLowBits | low
// 		// 	}
// 		// }
// 		// Als we hier gekomen zijn, is er geen groter getal in de huidige bucket.
// 		// We moeten dus het eerste getal in de volgende bucket returnen, indien
// 		// dit bestaat.
// 		// return ...

// 		// Als bkt-zeros < 0 of als b de laatste b is in ef.b, dan weten we
// 		// zeker dat de gezochte y in deze b zit. Als bkt-zeros == 0, dan
// 		// is het niet zeker dat de gezochte y in deze b zit en dan moet in
// 		// het ongunstige geval een extra b('s) gefetched worden!
// 		if bucket-zeros < 0 || len(ef.b) == i+1 { // bkt-zeros < 0 is enkel juist als achter de laatste bucket geen 0 meer komt!!!
// 			for { // Wat als x te groot is? Komen we dan in een oneindige lus?
// 				lowBits := ef.lowerBits(idx)
// 				if lowBits >= x&ef.mask {
// 					return x&^ef.mask | lowBits
// 				}
// 				idx++
// 			}
// 		}
// 		// Hier de andere gevallen.
// 		// We moeten controleren dat we tijdens de vergelijking van de
// 		// lowBits niet in een nieuwe bucket belanden zonder dat we het
// 		// weten. Dit gebeurt best door de volgende nul te zoeken?

// 		b ^= b & -b
// 	}
// 	bucket -= zeros
// }
// }

// Max returns the largest element in the Elias-Fano data structure.
func (ef *EliasFano) Max() uint32 {
	return ef.Access(ef.m - 1)
}

// Min returns the smallest element in the Elias-Fano data structure.
func (ef *EliasFano) Min() uint32 {
	return ef.Access(0)
}

// func (ef *EliasFano) nextHigherBits

// decode returns a list with positions of 1-bits.
// TODO: Als ik m ken, moet ik de list niet laten groeien!!!
func (ef *EliasFano) decode() (list []int) {
	for i, b := range ef.b {
		for b != 0 {
			list = append(list, 64*i+bits.TrailingZeros64(b))
			b &= b - 1
		}
	}
	return
}

// Values exports the numbers from the elias-fano data structure.
func (ef *EliasFano) Values() []uint32 {
	list := make([]uint32, ef.m)
	offset := ef.lHighBits
	var idx uint32
	for k, b := range ef.b {
		b &= 1<<(ef.lHighBits-uint32(k<<6)) - 1 // TODO: Deze lijn testen.
		// Het verwondert mij dat er geen overflow wordt gemeld!!!
		// Anders shift argument beperken met &63!
		for b != 0 {
			hval := uint32(64*k+bits.TrailingZeros64(b)) - idx
			lval := uint32(ef.b[offset>>6]>>(offset&63)) & ef.mask
			if (offset&63)+ef.nLowBits > 64 {
				lval |= uint32(ef.b[offset/64+1]<<(64-(offset&63))) & ef.mask
			}
			list[idx] = hval<<ef.nLowBits | lval

			b &= b - 1
			offset += ef.nLowBits
			idx++
		}
	}
	return list
}

// nextHigherBits ...onderdeel van een iterator!
// func (ef *EliasFano) nextHigherBits() int64 {
// 	for word == 0 {
// 		curr++
// 		word = upperBits[curr]
// 	}

// 	upperBits := curr<<6 + builtin_ctzll(word) - index
// 	index++
// 	word &= word - 1

// 	return upperBits
// }

// Ook bulk methode voorzien waar je meerdere opeenvolgende reads in een tijd kan doen!
// Probeer read en write te gebruiken ipv get and set!!! Vergelijk met stdlib

// Search ...
func (ef *EliasFano) Search(value uint64) bool {
	// TODO: nuttig?
	return false
}

// higherBits returns the higher bit value (bucket value) of the i-th number.
// Dit is niet efficient voor grote data structures!!! TODO index en select0 gebruiken
func (ef *EliasFano) higherBits(i int) uint32 {
	nOnes := i
	for j, bb := range ef.b {
		ones := bits.OnesCount64(bb)
		if i-ones <= 0 {
			return uint32(rs.Select64(bb, uint32(i))) + uint32(j<<6-nOnes)
		}
		i -= ones
	}
	return 0 // Geen fout teruggeven!!! Dit wordt op een hoger niveau geregeld.
}

// lowerBits returns the lower bits of the i-th number.
func (ef *EliasFano) lowerBits(i int) uint32 {
	offset := ef.lHighBits + uint32(i)*ef.nLowBits
	val := uint32(ef.b[offset>>6]>>(offset&63)) & ef.mask
	if (offset&63)+ef.nLowBits > 64 { // Dit kan efficienter met unsafe!!!
		val |= uint32(ef.b[offset>>6+1]<<(64-(offset&63))) & ef.mask
	}
	return val
}
