package ef

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestIterator(t *testing.T) {
	for i, tc := range testCases {
		d, err := From(tc.in)
		if err != nil {
			t.Error(err)
		}

		it := d.Iterator()

		for j := range tc.in {
			got, gok := it.Next()
			want, wok := d.Value(j)

			if got != want || gok != wok {
				t.Errorf(
					"tc: %d, t: %d, got: %d, want: %d, gok: %v, wok: %v\n",
					i, j, got, want, gok, wok,
				)
			}
		}

		for j := range tc.in {
			got, gok := it.Value(j)
			want, wok := d.Value(j)

			if got != want || gok != wok {
				t.Errorf(
					"tc: %d, t: %d, got: %d, want: %d, gok: %v, wok: %v\n",
					i, j, got, want, gok, wok,
				)
			}

			for k := j + 1; k < len(tc.in); k++ {
				got, gok := it.Next()
				want, wok := d.Value(k)

				if got != want || gok != wok {
					t.Errorf(
						"tc: %d, t: %d, k: %d, got: %d, want: %d, gok: %v, wok: %v\n",
						i, j, k, got, want, gok, wok,
					)
				}
			}
		}
	}
}

func BenchmarkIterator(b *testing.B) {
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

	it := d.Iterator()

	b.Run(fmt.Sprintf("Next([%d])", n), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v, ok := it.Value(0)
			for i := 0; i < n; i++ {
				v, ok = it.Next()
			}
			_ = v
			_ = ok
		}
	})
}
