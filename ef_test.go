package ef

import (
	"fmt"
	"math/bits"
	"testing"

	"github.com/gevg/rs"
)

var testCases = []struct {
	in  []uint32
	out []uint64
	l   uint32
	n   int
}{
	{ //
		[]uint32{0, 8, 32},
		[]uint64{0x498}, // 0b 0001 1001 0010    (6:0110 7:0111 10:1010)  0b 0100 1001 1000
		3,               // 0 low bits
		3,
	},
	{ //
		[]uint32{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29, 31, 33, 35, 37, 39, 41, 43, 45, 47, 49},
		[]uint64{0x498}, // 0b 0001 1001 0010    (6:0110 7:0111 10:1010)  0b 0100 1001 1000
		0,               // 0 low bits
		25,
	},
	{ //
		[]uint32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]uint64{0x498}, // 0b 0001 1001 0010    (6:0110 7:0111 10:1010)  0b 0100 1001 1000
		0,               // 0 low bits
		36,
	},
	{ //
		[]uint32{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		[]uint64{0x498}, // 0b 0001 1001 0010    (6:0110 7:0111 10:1010)  0b 0100 1001 1000
		0,               // 0 low bits
		35,
	},
	{ //
		[]uint32{3, 9, 130},
		[]uint64{0x498}, // 0b 0001 1001 0010    (6:0110 7:0111 10:1010)  0b 0100 1001 1000
		5,               // 1 low bits
		3,
	},
	{ // CSA++
		[]uint32{6, 7, 10}, // 3 numbers
		[]uint64{0x498},    // 0b 0001 1001 0010    (6:0110 7:0111 10:1010)  0b 0100 1001 1000
		1,                  // 1 low bits
		12,
	},
	{ // Index Construction - Compression of Postings
		[]uint32{1, 4, 7, 18, 24, 26, 30, 31}, // 8 numbers
		[]uint64{0x30505A},                    // 0b
		1,                                     // 2 low bits
		32,
	},
	{ // (2017) Compressed Suffix Arrays to Self-Indexes based on Partitioned Elias-Fano
		[]uint32{5, 8, 8, 15, 32}, // 5 numbers
		[]uint64{0x30505A},        // 0b 11 0000 0101 0000 0101 1010
		2,                         // 2 low bits
		24,
	},
	{
		[]uint32{5, 8, 9, 10, 14, 32}, // 6
		[]uint64{0x148A0BA},           // 0b 1 0100 1000 1010 0000 1011 1010
		2,                             // 2 low bits
		27,
	},
	{
		[]uint32{3, 4, 7, 13, 14, 15, 21, 25, 36, 38, 54, 62}, // 12 numbers
		[]uint64{0xED42014091A2A},                             // 0b 1110 1101 0100 0010 0000 0001 0100 0000 1001 0001 1010 0010 1010
		2,                                                     // 2 lowbits
		56,
	},
	{
		[]uint32{2, 3, 5, 7, 11, 13, 24}, // 7 numbers
		[]uint64{0x3E40A56},              // 0b 11 1110 0100 0000 1010 0101 0110
		1,                                // 2 lowbits
		27,
	},
	{ // Steve_Thomas (er is geen laatste stopbit in de highbits)
		[]uint32{4, 5, 6, 13, 22, 25}, // 6 numbers
		[]uint64{0x152822C},           // 0b 1 0101 0010 1000 0010 0010 1100
		2,                             // 2 lowbits
		25,
	},
	{ // (2016) Elias-Fano encoding
		[]uint32{3, 4, 7, 13, 14, 15, 21, 43}, // 8 numbers
		[]uint64{0x6F39A09CD},                 // 0b 110 1111 0011 1001 1010 0000 1001 1100 1101
		2,                                     // 3 lowbits
		35,
	},
	{ //
		[]uint32{3, 4, 13, 15, 24, 26, 27, 29}, // 8 numbers
		[]uint64{0x66AD050A},                   // 0b 110 0110 1010 1101 0000 0101 0000 1010
		1,                                      // 3 lowbits
		31,
	},
	{ // sds-tutorial
		[]uint32{8, 9, 11, 13, 16, 18, 23}, // 7 numbers
		[]uint64{0x27252B0},                // 0b 10 0111 0010 0101 0010 1011 0000
		1,                                  // 2 lowbits
		26,
	},
	{ // Text Indexing
		[]uint32{4, 13, 15, 24, 26, 27, 29}, // 7 numbers
		[]uint64{0x19968284},                // 0b 1 1001 1001 0110 1000 0010 1000 0100
		2,                                   // 3 lowbits
		29,
	},
}

func TestFilter(t *testing.T) {
	var b uint64 = 10 // 0b 1010
	p := rs.Select64(b, 0) + 1
	b &^= (1 << uint(p)) - 1
	fmt.Printf("new b: %b, p: %d\n", b, p)
}

func TestHigherBits(t *testing.T) {
	for i, tc := range testCases {
		ef := From(tc.in)
		for j, in := range tc.in {
			want := in >> ef.nLowBits
			if got := ef.higherBits(j); got != want {
				t.Errorf("tc: %d, i: %d, got: %d, want: %d\n", i, j, got, want)
			}
		}
	}
}

func TestFrom(t *testing.T) {
	for _, tc := range testCases {
		ef := From(tc.in)
		got := ef.nLowBits
		want := tc.l
		if got != want {
			t.Errorf(
				"nLow - got: %d, want: %d, total: %d\n",
				ef.nLowBits, tc.l, ef.lHighBits+uint32(ef.m)*ef.nLowBits,
			)
		}
	}
}

func TestAccess(t *testing.T) {
	for i, tc := range testCases {
		ef := From(tc.in)
		for j, want := range tc.in {
			if got := ef.Access(j); got != want {
				t.Errorf("tc: %d, t: %d, got: %d, want: %d\n", i, j, got, want)
			}
		}
	}
}

func TestNextLEQ(t *testing.T) {
	in := []uint32{3, 4, 7, 13, 14, 15, 21, 25, 36, 38, 54, 62}
	values := []uint32{5, 60, 70}
	want := []int{1, 10, 11}
	// values := []uint32{1, 5, 60}
	// want := []uint32{3, 7, 62}
	ef := From(in)
	for i, v := range values {
		if got := ef.NextLEQ(v); got != want[i] {
			t.Errorf("x: %d, got: %d, want: %d\n", v, got, want[i])
		}
	}
}

// func TestNextGEQ(t *testing.T) { // Er zit hier een fout in!!! Checken!!!
// 	in := []uint32{3, 4, 7, 13, 14, 15, 21, 25, 36, 38, 54, 62}
// 	xs := []uint32{5, 60}
// 	want := []uint32{7, 62}
// 	// xs := []uint32{1, 5, 60}
// 	// want := []uint32{3, 7, 62}
// 	ef := From(in)
// 	for i, x := range xs {
// 		if got := ef.NextGEQ(x); got != want[i] {
// 			t.Errorf("x: %d, got: %d, want: %d\n", x, got, want[i])
// 		}
// 	}
// }

func TestNumLow(t *testing.T) {
	for _, tc := range testCases {
		ef := From(tc.in)

		// nLowSlow := floor(log(upperBound / size))
		m := uint32(len(tc.in))
		nLowSlow := uint32(bits.Len32(tc.in[m-1]/m) - 1)

		if got := ef.nLowBits; got != nLowSlow {
			t.Errorf("nLow - got: %d, want: %d\n", got, nLowSlow)
		}
	}
}

func TestMax(t *testing.T) {
	for i, tc := range testCases {
		ef := From(tc.in)
		want := tc.in[len(tc.in)-1]
		got := ef.Max()
		if got != want {
			t.Errorf("tc: %d, got: %d, want: %d\n", i, got, want)
		}
	}
}

func TestMin(t *testing.T) {
	for i, tc := range testCases {
		ef := From(tc.in)
		want := tc.in[0]
		got := ef.Min()
		if got != want {
			t.Errorf("tc: %d, got: %d, want: %d\n", i, got, want)
		}
	}
}

func TestDecode(t *testing.T) {
	ef := From([]uint32{4, 8, 12})
	ef.b[0] = 8
	t.Errorf("in: %b, got: %v, want: %v\n", ef.b[0], ef.decode(), []uint64{3})
}

func TestValues(t *testing.T) {
	for k, tc := range testCases {
		ef := From(tc.in)
		for i, el := range ef.Values() {
			if el != tc.in[i] {
				t.Errorf("Testcase %d renders faulty extractions\n", k)
				break
			}
		}
	}
}

func BenchmarkNew(b *testing.B) { // ~10 ns/op/number  8 B/op  1 allocs/op
	for k, tc := range testCases {
		b.Run(fmt.Sprintf("New(%d)", k), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				From(tc.in)
			}
		})
	}
}

func BenchmarkNew2(b *testing.B) { // 72.5 ns/op  8 B/op  1 allocs/op
	arr := []uint32{3, 4, 7, 13, 14, 15, 21, 25, 36, 38, 54, 62}
	for i := 0; i < b.N; i++ {
		From(arr)
	}
}

func BenchmarkAccess(b *testing.B) {
	for idx, tc := range testCases {
		ef := From(tc.in)

		b.Run(fmt.Sprintf("Access(%d)", idx), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := range tc.in {
					ef.Access(j)
				}
			}
		})
	}
}

func BenchmarkAccess2(b *testing.B) { // ~10.0 ns/op  0 B/op  0 allocs/op
	arr := []uint32{3, 4, 7, 13, 14, 15, 21, 25, 36, 38, 54, 62}
	ef := From(arr)

	for idx := range arr {
		b.Run(fmt.Sprintf("Access(%d)", idx), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ef.Access(idx)
			}
		})
	}
}

// func BenchmarkNextGEQ(b *testing.B) { // ~8.96 ns  0 B/op  0 allocs/op
// 	in := []uint32{3, 4, 7, 13, 14, 15, 21, 25, 36, 38, 54, 62}
// 	xs := []uint32{1, 5, 60}

// 	ef := From(in)

// 	b.ReportAllocs()
// 	b.ResetTimer()

// 	for _, x := range xs {
// 		b.Run(fmt.Sprintf("NextGEQ(%d)", x), func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				ef.NextGEQ(x)
// 			}
// 		})
// 	}
// }

func BenchmarkValues(b *testing.B) {
	for k, tc := range testCases {
		b.Run(fmt.Sprintf("Values(%d)", k), func(b *testing.B) {
			ef := From(tc.in)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				ef.Values()
			}
		})
	}
}

func BenchmarkMax(b *testing.B) { // 0.80 ns/op  0 B/op  0 allocs/op
	for k, tc := range testCases {
		b.Run(fmt.Sprintf("Max(%d)", k), func(b *testing.B) {
			ef := From(tc.in)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				ef.Max()
			}
		})
	}
}

func BenchmarkMin(b *testing.B) { // 0.73 ns/op  0 B/op  0 allocs/op
	for k, tc := range testCases {
		b.Run(fmt.Sprintf("Min(%d)", k), func(b *testing.B) {
			ef := From(tc.in)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				ef.Min()
			}
		})
	}
}
