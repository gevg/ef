# Dictionary with Elias-Fano compression

<p align="center">
  <a href="https://github.com/gevg/ef/actions"><img src="https://github.com/gevg/ef/workflows/ci/badge.svg" alt="Build Status" /></a>
  <a href="https://pkg.go.dev/github.com/gevg/ef"><img src="https://pkg.go.dev/badge/github.com/gevg/ef.svg" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/gevg/ef"><img src="https://goreportcard.com/badge/github.com/gevg/ef?style=flat-square" alt="Go Report Card" /></a>
  <a href="https://codecov.io/gh/gevg/ef"><img src="https://codecov.io/gh/gevg/ef/branch/main/graph/badge.svg" alt="codecov" /></a>
</p>

The static Elias-Fano representation encodes a sequence of positive, monotonically increasing integers.
The representation is quasi-succinct, meaning it is less than half a bit per element away from the 0-th order empirical entropy. It supports random access of an integer in O(1) worst-case time, despite the compression.
However, for random access to be fast in practice, it requires the integration of a Select index. The representation also facilitates efficient search capabilities, such as `NextGEQ()` and `NextLEQ()`.

The Elias-Fano representation is built into a dictionary. Its purpose is to keep large volumes of integers in RAM in a compressed format, while still being efficient to use. The fastest access method is through iteration, but direct access is fast as well.

This is an implementation of a static Elias-Fano representation. So, stored values cannot be updated. If this is essential to your use case, please have a look at [2].

The implementation is to a large extent based on the paper of S. Vigna [1]. Forward and skip pointers are not implemented. This implementation relies on a fast Select index instead. Storage of unsorted and negative integers is on the roadmap.

The API is NOT stable yet. You can send feedback about the current API by reporting an issue.

## Installation

Install `ef` with `go get`:

```bash
$ go get -u github.com/gevg/ef
```

## Usage

Construct the dictionary, while initializing it with a number array, and iterate over the values.

```go
// Initialize a slice of monotonically non-decreasing numbers and construct the compressed dictionary.
values := []uint{1, 3, 5, 7, 9, 11, 11, 13, 15, 17}
dict, err := From(values)
```

When not all values are available at construction, or when the construction of a huge values slice is not desirable, one can construct a dictionary using `New()` and add each entry individually with `Append()`. The dictionary is ready to use after each Append. There is no need to completely fill the dictionary. However, a partially filled dictionary will show inferior compression.

```go
// Construct the dictionary with New(). Supply the number of values and
// the maximum value as configuration parameters. Add the entries.
nValues := len(values)
maxValue := values[nValues-1]

dict, err := New(nValues, maxValue)
for _, v := range values {
    if err := dict.Append(v); err != nil {
        ...
    }
}
```

```go
// Cap() returns the capacity (max. number of entries) of the dictionary.
// Len() returns the current number of input values in the dictionary.
cap := Cap(dict)
length := Len(dict)
```

There are several ways to read back data. One can directly read one or more values, starting at a given index.

```go
// Read the value at index 2.
v, ok := dict.Value(2)

// Read 5 values starting at index start := 4.
vals := make([]uint, 5)
vals, ok := v.RangeValues(start, vals)

// Read all values from the dictionary.
for k, v := range dict.Values() {
    ...
}
```

Alternatively, one can use an iterator. This allows for concurrent access.

```go
// Create an iterator.
it := dict.Iterator()
vals := make([]uint, 0, m)

// Retrieve the first m values of the dictionary.
for i := 0; i < m; i++ {
    v, _ := it.Next()
    vals = append(vals, v)
}
```

The iterator does not have to start reading at index 0. A read with `it.Value(k)` returns the k-th value and sets the iterator state on index `k+1`. Subsequent reads with `Next()` allows to read the other entries.

```go
// Retrieval of the values in index range [k, l). Read the k-th value and add it to the vals slice.
vals := make([]uint, 0, l-k)
v, ok := it.Value(k)
if !ok {
    ...
}
vals = append(vals, v)

// Read the other values by iterating with Next().
for i:=1; i < k-l; i++ {
    if v, ok = it.Next(); !ok {
        ...
    }
    vals = append(vals, v)
}
```

## Benchmarks

TODO

## To Do

Implementation of the following features:

- [ ] Integration of a fast Select() algorithm
- [ ] Support of int values
- [ ] Support of strict monotonicity.
- [ ] Support of non-monotonic series
- [ ] Support of search: NextLEQ() and NextGEQ()
- [ ] Benchmarks and further performance improvements

## License

`ef` is available under the [BSD 3-Clause License](LICENSE).

## References

[1] S. Vigna: Quasi-succinct Indices, WSDM '13: Proceedings of the sixth ACM
international conference on Web search and data mining, 2013, Pages 83–92

[2] G.E. Pibiri, R. Venturini: Dynamic Elias-Fano Representation, 28th Annual
Symposium on Combinatorial Pattern Matching (CPM 2017), Warsaw, Poland, 2017
