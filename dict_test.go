package ef

import (
	"fmt"
	"math/rand"
	"testing"
)

var testCases = []struct {
	in         []uint
	out        []uint
	sizeLValue int
	sizeT      int
}{
	// {
	// 	[]uint{},
	// 	[]uint{},
	// 	0,
	// 	1,
	// },
	{
		[]uint{0}, // H|L: 1 + 0, u|l: 1|0, H: 1+0
		[]uint{0b1},
		0,
		1,
	},
	{
		[]uint{32}, // H|L: 2 + 5, u|l: 6|5, H: 1+1
		[]uint{0b000_0010},
		5,
		7,
	},
	{
		[]uint{2, 4, 8, 30, 32}, // H|L: 13 + 10, u|l: 6|2, H: 5+8
		[]uint{0b001_0000_0101_0100_0001_0101},
		2,
		23,
	},
	{
		[]uint{5, 5, 5, 5, 5, 5, 5, 5, 5, 31}, // H|L: 25 + 10, u|l: 5|1, H: 10+15
		[]uint{0b111_1111_1111_0000_0000_0000_0111_1111_1100},
		1,
		35,
	},
	{
		[]uint{1, 4, 7, 18, 24, 26, 30, 31}, // H|L: 23 + 8, u|l: 5|1, H: 8+15
		[]uint{0b100_0010_1110_0101_0001_0000_0010_1001},
		1,
		31,
	},
	{
		[]uint{5, 8, 8, 15, 32}, // H|L: 13 + 10, u|l: 6|2, H: 5+8
		[]uint{0b001_1000_0011_0000_0101_1010},
		2,
		23,
	},
	{
		[]uint{0, 8, 32}, // H|L: 7 + 9, u|l: 6|3, H: 3+4
		[]uint{0b_0000_0000_0100_0101},
		3,
		16,
	},
	{
		[]uint{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29, 31, 33, 35, 37, 39, 41, 43, 45, 47, 49},
		[]uint{
			0b_0010_0100_1001_0010_0100_1001_0010_0100_1001_0010_0100_1001_0010_0100_1001_0010,
			0b10_0100_1001,
		}, // H|L: 74 + 0, u|l: 6|0, H: 25+49
		0,
		74,
	},
	{
		[]uint{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]uint{0b_1111_1111_1111_1111_1111_1111_1111_1111_1111}, // H|L: 36 + 0, u|l: 1|0, H: 36+0
		0,
		36,
	},
	{
		[]uint{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		[]uint{0b_1111_1111_1111_1111_1111_1111_1111_1111_1110}, // H|L: 36 + 0, u|l: 1|0, H: 35+1
		0,
		36,
	},
	{
		[]uint{3, 9, 130}, // H|L: 7 + 15, u|l: 8|5, H: 3+4
		[]uint{0b_00_0100_1001_0001_1100_0011},
		5,
		22,
	},
	{
		[]uint{6, 7, 10}, // H|L: 8 + 3, u|l: 4|1, H: 3+5
		[]uint{0b010_1001_1000},
		1,
		11,
	},
	{
		[]uint{1, 4, 7, 18, 24, 26, 30, 31}, // H|L: 23 + 8, u|l: 5|1, H: 8+15
		[]uint{0b100_0010_1110_0101_0001_0000_0010_1001},
		1,
		31,
	},
	{
		[]uint{5, 8, 8, 15, 32}, // H|L: 13 + 10, u|l: 6|2, H: 5+8
		[]uint{0b001_1000_0011_0000_0101_1010},
		2,
		23,
	},
	{
		[]uint{5, 8, 9, 10, 14, 32}, // H|L: 14 + 12, u|l: 6|2, H: 6+8
		[]uint{0b00_1010_0100_0110_0000_1011_1010},
		2,
		26,
	},
	{
		[]uint{3, 4, 7, 13, 14, 15, 21, 25, 36, 38, 54, 62}, // H|L: 27|24, u|l: 6|2, H: 12+15
		[]uint{0b101_0100_0010_1111_0011_1001_1100_1000_0110_0010_1001_1100_1101},
		2,
		51,
	},
	{
		[]uint{2, 3, 5, 7, 11, 13, 24}, // H|L: 19|7, u|l: 5|1, H: 7+12
		[]uint{0b01_1111_0100_0000_1010_0101_0110},
		1,
		26,
	},
	{
		[]uint{4, 5, 6, 13, 22, 25},
		[]uint{0b_0110_0110_0100_1010_0100_1110}, // H|L: 12|12, u|l: 5|2, H: 6+6
		2,
		24,
	},
	{
		[]uint{3, 4, 7, 13, 14, 15, 21, 43}, // H|L: 18|16, u|l: 6|2, H: 8+10
		[]uint{0b11_0111_1001_1100_1110_0000_1001_1100_1101},
		2,
		34,
	},
	{
		[]uint{3, 4, 13, 15, 24, 26, 27, 29}, // H|L: 22|8, u|l: 5|1, H: 8+14
		[]uint{0b11_0011_0110_1101_0000_0101_0000_1010},
		1,
		30,
	},
	{
		[]uint{8, 9, 11, 13, 16, 18, 23}, // H|L: 18|7, u|l: 5|1, H: 7+11
		[]uint{0b1_0011_1010_0101_0010_1011_0000},
		1,
		25,
	},
	{
		[]uint{4, 13, 15, 24, 26, 27, 29}, // H|L: 14|14, u|l: 5|2, H: 7+7
		[]uint{0b_0111_1000_1101_0010_1110_0011_0010},
		2,
		28,
	},
}

func TestFrom(t *testing.T) {
	for i, tc := range testCases {
		d, err := From(tc.in)
		if err != nil {
			t.Error(err)
			return
		}

		gotSizeLValue := d.sizeLValue
		wantSizeLValue := tc.sizeLValue
		gotSizeT := d.sizeH + d.n*d.sizeLValue
		wantSizeT := tc.sizeT
		if gotSizeLValue != wantSizeLValue || gotSizeT != wantSizeT {
			t.Errorf(
				"tc: %d, - sizeLValue - got: %d, want: %d, sizeT: - got: %d, want: %d\n",
				i, gotSizeLValue, wantSizeLValue, gotSizeT, wantSizeT,
			)
		}

		for j, want := range tc.out {
			got := d.b[j]
			if got != want {
				t.Errorf("tc: %d, - b[%d]\ngot:  %064b,\nwant: %064b\n", i, j, got, want)
			}
		}
	}
}

func TestValue(t *testing.T) {
	for i, tc := range testCases {
		d, err := From(tc.in)
		if err != nil {
			t.Error(err)
		}

		for j, want := range tc.in {
			if got, _ := d.Value(j); got != want {
				t.Errorf("tc: %d, t: %d, got: %d, want: %d\n", i, j, got, want)
			}
		}
	}
}

func TestNewAndAppend(t *testing.T) {
	for i, tc := range testCases {
		d, err := New(len(tc.in), tc.in[len(tc.in)-1])
		if err != nil {
			t.Error(err)
		}

		for j, v := range tc.in {
			if k := d.Append(v); k == -1 {
				t.Errorf("tc: %d, t: %d: error appending %d\n", i, j, v)
			}
		}

		for j, want := range tc.in {
			if got, _ := d.Value(j); got != want {
				t.Errorf("tc: %d, t: %d, got: %d, want: %d\n", i, j, got, want)
			}
		}
	}
}

func TestHValue(t *testing.T) {
	for i, tc := range testCases {
		d, err := From(tc.in)
		if err != nil {
			t.Error(err)
		}

		for j, in := range tc.in {
			got := d.hValue(j)
			want := in >> d.sizeLValue
			if got != want {
				t.Errorf("tc: %d, i: %d, got: %d, want: %d\n", i, j, got, want)
			}
		}
	}
}

func TestLValue(t *testing.T) {
	for i, tc := range testCases {
		d, err := From(tc.in)
		if err != nil {
			t.Error(err)
		}

		for j, in := range tc.in {
			got := d.lValue(j)
			want := in & (1<<d.sizeLValue - 1)
			if got != want {
				t.Errorf("tc: %d, i: %d, got: %d, want: %d\n", i, j, got, want)
			}
		}
	}
}

func TestValues(t *testing.T) {
	for k, tc := range testCases {
		d, err := From(tc.in)
		if err != nil {
			t.Error(err)
			return
		}

		for i, got := range d.Values() {
			want := tc.in[i]
			if got != want {
				t.Errorf("Testcase: %d, i: %d, value got: %d != %d want\n", k, i, got, want)
				break
			}
		}
	}
}

func TestValuesLong(t *testing.T) {
	const (
		n          = 10_000
		max        = 100
		iterations = 10_000
	)

	in := make([]uint, n)
	rand.Seed(18)

	for k := 0; k < iterations; k++ {
		in = in[:1+rand.Intn(n)]
		var prev uint
		for i := range in {
			prev += uint(rand.Intn(max))
			in[i] = prev
		}

		t.Run(fmt.Sprintf("From([%d])", n), func(t *testing.T) {
			d, err := From(in)
			if err != nil {
				t.Error(err)
				return
			}

			for i, got := range d.Values() {
				want := in[i]
				if got != want {
					t.Errorf("Iteration: %d, i: %d, value got: %d != %d want\n", k, i, got, want)
					break
				}
			}
		})

		t.Run(fmt.Sprintf("Append([%d])", n), func(t *testing.T) {
			d, err := New(len(in), in[len(in)-1])
			if err != nil {
				t.Error(err)
				return
			}

			for i, v := range in {
				if j := d.Append(v); j == -1 {
					t.Errorf("Iteration: %d, i: %d, Append failed\n", k, i)
					break
				}
			}

			for i, got := range d.Values() {
				want := in[i]
				if got != want {
					t.Errorf("Iteration: %d, i: %d, value got: %d != %d want\n", k, i, got, want)
					break
				}
			}
		})
	}
}

// func TestNextGEQ(t *testing.T) {
// 	for i, tc := range testCases {
// 		d, err := From(tc.in)
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		for j, want := range tc.in {
// 			k, got := d.NextGEQ(want)
// 			if j != k || got != want {
// 				t.Errorf("tc: %d, (k, got): (%d, %d), (k, want): (%d, %d)\n", i, k, got, j, want)
// 			}
// 		}
// 	}
// }

func BenchmarkConstruction(b *testing.B) {
	const (
		n   = 1_000_000
		max = 100
	)

	in := make([]uint, n)

	rand.Seed(18)
	var prev uint
	for i := range in {
		prev += uint(rand.Intn(max))
		in[i] = prev
	}

	b.Run(fmt.Sprintf("From([%d])", n), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			From(in)
		}
	})

	b.Run(fmt.Sprintf("Append([%d])", n), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d, _ := New(len(in), in[len(in)-1])

			for _, v := range in {
				d.Append(v)
			}
		}
	})
}

func BenchmarkSequentialAccess(b *testing.B) {
	const (
		n    = 1_000_000
		max  = 100
		step = 100
	)

	in := make([]uint, n)
	rand.Seed(18)
	var prev uint
	for i := range in {
		prev += uint(rand.Intn(max))
		in[i] = prev
	}

	d, err := From(in)
	if err != nil {
		b.Error(err)
	}

	b.Run(fmt.Sprintf("Value([%d])", len(in)), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := 0; j < n; j += step {
				d.Value(j)
			}
		}
	})
}

func BenchmarkRandomAccess(b *testing.B) {
	const (
		n    = 1_000_000
		max  = 100
		nIdx = 10_000
	)

	rand.Seed(18)
	idx := make([]int, nIdx)
	for i := range idx {
		idx[i] = rand.Intn(n)
	}

	var prev uint
	in := make([]uint, n)
	for i := range in {
		prev += uint(rand.Intn(max))
		in[i] = prev
	}

	d, err := From(in)
	if err != nil {
		b.Error(err)
	}

	b.Run(fmt.Sprintf("Value([%d])", len(in)), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, j := range idx {
				d.Value(j)
			}
		}
	})

	b.Run(fmt.Sprintf("Value2([%d])", len(in)), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, j := range idx {
				d.Value2(j)
			}
		}
	})
}

func BenchmarkValues(b *testing.B) {
	const (
		n   = 1_000_000
		max = 100
	)

	in := make([]uint, n)
	rand.Seed(18)
	var prev uint
	for i := range in {
		prev += uint(rand.Intn(max))
		in[i] = prev
	}

	d, err := From(in)
	if err != nil {
		b.Error(err)
	}

	b.Run(fmt.Sprintf("Values([%d])", len(in)), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d.Values() // values slice meegeven om relevante tijden te krijgen.
		}
	})
}

func BenchmarkMapAccess(b *testing.B) {
	const (
		n   = 1_000_000
		max = 100
	)

	m := make(map[int]uint, n)
	rand.Seed(18)
	var prev uint
	for i := 0; i < n; i++ {
		prev += uint(rand.Intn(max))
		m[i] = prev
	}

	b.Run(fmt.Sprintf("Access([%d])", len(m)), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var v uint
			for k := 0; k < n; k++ {
				v = m[k] // values slice meegeven om relevante tijden te krijgen.
			}
			_ = v
		}
	})

	b.Run(fmt.Sprintf("Iteration([%d])", len(m)), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var k int
			var v uint
			for k, v = range m {
			}
			_ = k
			_ = v
		}
	})
}
