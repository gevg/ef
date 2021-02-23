# Comparisons

## dict vs []uint

Construction

- features: draw

- performance: []uint wins

Size: dict wins

Access performance

- sequential: []uint wins

- random: []uint wins

Search:

- features: draw (both need to be sorted)
- performance: ?

Summary: Smaller size is the only reason to choose Elias-Fano above []uint.

## dict vs map[int]uint or map[uint]bool

Construction

- features: draw

- performance: dict wins

Size: dict wins

Access performance

- sequential: dict wins

- random: dict wins

Search: draw in terms of features (if map is sorted), ? in terms of performance

Summary: There are many reasons to choose Elias-Fano above map[uint]uint.

## dict vs varInt (pkg encoding/binary)

Construction

- features: varint wins (no need to sort)

- performance: varint wins

Size: dict wins

Access performance

- sequential: ?

- random: dict wins (varint has no random access)

Search: dict wins (varint has no ranndom access)

Summary: There are many reasons to choose Elias-Fano above varint.

## dict vs bit encoders (Elias-Gamma, Delta, Omega, Golomb-Rice, Fibonacci, ...)
