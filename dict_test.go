package ef

import (
	"fmt"
	"math/rand"
	"testing"
)

var testCases = []struct {
	in         []uint
	out        []uint // TODO: Ga ik out nog ooit gebruiken???
	sizeLValue int
	sizeT      int
}{
	{
		[]uint{},
		[]uint{0x498},
		0,
		1,
	},
	{
		[]uint{0}, // H|L: 1 + 0, u|l: 1|0, H: 1+0
		[]uint{0x498},
		0,
		1,
	},
	{
		[]uint{32}, // H|L: 2 + 5, u|l: 6|5, H: 1+1
		[]uint{0x498},
		5,
		7,
	},
	{
		[]uint{2, 4, 8, 30, 32}, // H|L: 13 + 10, u|l: 6|2, H: 5+8
		[]uint{0x498},
		2,
		23,
	},
	{
		[]uint{5, 5, 5, 5, 5, 5, 5, 5, 5, 31}, // H|L: 25 + 10, u|l: 5|1, H: 10+15
		[]uint{0x498},
		1,
		35,
	},
	{
		[]uint{1, 4, 7, 18, 24, 26, 30, 31}, // H|L: 23 + 8, u|l: 5|1, H: 8+15
		[]uint{0x498},
		1,
		31,
	},
	{
		[]uint{5, 8, 8, 15, 32}, // H|L: 13 + 10, u|l: 6|2, H: 5+8
		[]uint{0x498},
		2,
		23,
	},
	{
		[]uint{0, 8, 32}, // H|L: 7 + 9, u|l: 6|3, H: 3+4
		[]uint{0x498},
		3,
		16,
	},
	{
		[]uint{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29, 31, 33, 35, 37, 39, 41, 43, 45, 47, 49},
		[]uint{0x498}, // H|L: 74 + 0, u|l: 6|0, H: 25+49
		0,
		74,
	},
	{
		[]uint{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]uint{0x498}, // H|L: 36 + 0, u|l: 1|0, H: 36+0
		0,
		36,
	},
	{
		[]uint{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		[]uint{0x498}, // H|L: 36 + 0, u|l: 1|0, H: 35+1
		0,
		36,
	},
	{
		[]uint{3, 9, 130},
		[]uint{0x498}, // H|L: 7 + 15, u|l: 8|5, H: 3+4
		5,
		22,
	},
	{
		[]uint{6, 7, 10},
		[]uint{0x498}, // H|L: 8 + 3, u|l: 4|1, H: 3+5
		1,
		11,
	},
	{
		[]uint{1, 4, 7, 18, 24, 26, 30, 31},
		[]uint{0x30505A}, // H|L: 23 + 8, u|l: 5|1, H: 8+15
		1,
		31,
	},
	{
		[]uint{5, 8, 8, 15, 32},
		[]uint{0x30505A}, // H|L: 13 + 10, u|l: 6|2, H: 5+8
		2,
		23,
	},
	{
		[]uint{5, 8, 9, 10, 14, 32},
		[]uint{0x148A0BA}, // H|L: 14 + 12, u|l: 6|2, H: 6+8
		2,
		26,
	},
	{
		[]uint{3, 4, 7, 13, 14, 15, 21, 25, 36, 38, 54, 62},
		[]uint{0xED42014091A2A}, // H|L: 27|24, u|l: 6|2, H: 12+15
		2,
		51,
	},
	{
		[]uint{2, 3, 5, 7, 11, 13, 24}, // H|L: 19|7, u|l: 5|1, H: 7+12
		[]uint{0b11_1110_0100_0000_1010_0101_0110},
		1,
		26,
	},
	{
		[]uint{4, 5, 6, 13, 22, 25},
		[]uint{0x152822C}, // H|L: 12|12, u|l: 5|2, H: 6+6
		2,
		24,
	},
	{
		[]uint{3, 4, 7, 13, 14, 15, 21, 43},
		[]uint{0x6F39A09CD}, // H|L: 18|16, u|l: 6|2, H: 8+10
		2,
		34,
	},
	{
		[]uint{3, 4, 13, 15, 24, 26, 27, 29},
		[]uint{0x66AD050A}, // H|L: 22|8, u|l: 5|1, H: 8+14
		1,
		30,
	},
	{
		[]uint{8, 9, 11, 13, 16, 18, 23},
		[]uint{0x27252B0}, // H|L: 18|7, u|l: 5|1, H: 7+11
		1,
		25,
	},
	{
		[]uint{4, 13, 15, 24, 26, 27, 29},
		[]uint{0x19968284}, // H|L: 14|14, u|l: 5|2, H: 7+7
		2,
		28,
	},
}

func TestFrom(t *testing.T) {
	for _, tc := range testCases {
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
				"sizeLValue - got: %d, want: %d, sizeT: - got: %d, want: %d\n",
				gotSizeLValue, wantSizeLValue, gotSizeT, wantSizeT,
			)
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

func TestValues2(t *testing.T) {
	for k, tc := range testCases {
		d, err := From(tc.in)
		if err != nil {
			t.Error(err)
			return
		}

		for i, got := range d.Values2() {
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

func TestValuesLong2(t *testing.T) {
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

			for i, got := range d.Values2() {
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

func TestNextGEQ(t *testing.T) {
	for i, tc := range testCases {
		d, err := From(tc.in)
		if err != nil {
			t.Error(err)
		}

		for j, want := range tc.in {
			k, got := d.NextGEQ(want)
			if j != k || got != want {
				t.Errorf("tc: %d, (k, got): (%d, %d), (k, want): (%d, %d)\n", i, k, got, j, want)
			}
		}
	}
}

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

	b.Run(fmt.Sprintf("Values2([%d])", len(in)), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d.Values2() // values slice meegeven om relevante tijden te krijgen.
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
