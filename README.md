# Dictionary with Elias-Fano compression

(early) stage work-in-progress, not to be used for production

static

quasi succinct

direct access to any value

uses Select() to locate values in a fast way.

fast? bounds checks? allocations?

monotonically increasing values

## Installation

Install `ef` with `go get`:

```bash
$ go get -u github.com/gevg/ef
```

## Usage

Construct the dictionary, initializing it with a number array, and iterate over the values.

```go
// Initialize a monotonically increasing number array, construct the dictionary,
// and iterate over the values.
values := []int32{1, 3, 5, 7, 9, 11, 11, 13, 15, 17}
nValues := len(values)
dict, err := From(values)

dict, err := New(nValues, values[nValues-1]) // up to nValues, up to values[nValues-1]
dict.Append(values[0])
dict.Append(values[1], values[2])
dict.Append(values[3:]) // Wat als je nog niet aan n bent? Welke n moet je dan in d.n hebben?
// We willen geen commando dat de overgang maakt van write naar read!
length := Len(dict)

k := 4
v, ok := dict.Value(k)
fmt.Printf("value[%d] = %d\n", k, v)

for k, v := range dict.Values() {
    fmt.Printf("value[%d] = %d\n", k, v)
}

// retrieve the first m values
it := dict.Iterator()
vals := make([]uint, 0, m)
for i := 0; i < m; i++ {
    v, _ := it.Next()
    vals = append(vals, v)
}

// retrieve the values with index [k, l)
vals := make([]uint, 0, l-k)
v, ok := it.Value(k)
if !ok {
    ...
}
vals = append(vals, v)

for i:=1; i < k-l; i++ {
    if v, ok = it.Next(); !ok {
        ...
    }
    vals = append(vals, v)
}
```

## References

[1] P. Elias: On binary representation of monotone sequences, Proceedings
of the 6-th Princeton Conference on Information Sciences and Systems:
54-47, Princeton University, 1972. †2 and 25

[2] P. Elias: Efficient storage and retrieval by content and adress of static
files, J. ACM 21(2): 246-260, 1974. †2 and 25

[3] R. M. Fano: On the number of bits required to implement an associative
memory, Memorandum 61, Computer Structures Group,
Project MAC, Massachusetts Institute of Technology, 1971. † 2
and 25

[4] S. Vigna: Quasi-succinct Indices, WSDM '13: Proceedings of the sixth ACM
international conference on Web search and data mining, 2013, Pages 83–92

[5] G.E. Pibiri, R. Venturini: Dynamic Elias-Fano Representation, 28th Annual
Symposium on Combinatorial Pattern Matching (CPM 2017), Warsaw, Poland, 2017
