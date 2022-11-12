# a fast flat-file index for constant datasets

This package specifies a file format and Go implementation for indexing constant datasets.

*`compactindex` …*
- is an immutable file format;
- maps arbitrary keys into offsets in an external flat file;
- consumes a constant amount of space per entry
  - ~6-8 bytes, regardless of key size
  - 3 bytes per enty 
- while querying, requires `2 + O(log2(n))` lookups worst- & average-case (binary search);
- during construction, requires near-constant memory space and `O(n)` scratch space with regard to entries per file;
- during construction, features a constant >500k entry/s per-core write rate (2.5 GHz Intel laptop);
- works on any storage supporting random reads (regular files, HTTP range requests, on-chain, ...);
- is based on the "FKS method" which uses perfect (collision-free) hash functions in a two-level hashtable; [^1]
- is inspired by D. J. Bernstein's "constant database"; [^2]
- uses the xxHash64 non-cryptographic hash-function; [^3]

Refer to the Go documentation for the algorithms used and implementation details.

[![Go Reference](https://pkg.go.dev/badge/go.firedancer.io/radiance/pkg/compactindex.svg)](https://pkg.go.dev/go.firedancer.io/radiance/pkg/compactindex)

[^1]: Fredman, M. L., Komlós, J., & Szemerédi, E. (1984). Storing a Sparse Table with 0 (1) Worst Case Access Time. Journal of the ACM, 31(3), 538–544. https://doi.org/10.1145/828.1884
[^2]: cdb by D. J. Bernstein https://cr.yp.to/cdb.html
[^3]: Go implementation of xxHash by @cespare: https://github.com/cespare/xxhash/

## Interface

In programming terms:

```rs
fn lookup(key: &[byte]) -> Option<u64>
```

Given an arbitrary key, the index
- states whether the key exists in the index
- if it exists, maps the key to an integer (usually an offset into a file)

## Example

Here are some example scenarios where `compactindex` is useful:

- When working with immutable data structures
  - Example: Indexing [IPLD CAR files][3] carrying Merkle-DAGs of content-addressable data
- When working with archived/constant data
  - Example: Indexing files in `.tar` archives
- When dealing with immutable remote storage such as S3-like object storage
  - Example: Storing the index and target file in S3, then using [HTTP range requests][4] to efficiently query data

[3]: https://ipld.io/specs/transport/car/
[4]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Range_requests

Here are some things compactindex cannot do:

- Cannot add more entries to an existing index
  - Reason 1: indexes are tightly packed, so there is no space to insert new entries (though `fallocate(2)` with `FALLOC_FL_INSERT_RANGE` would technically work)
  - Reason 2: the second-level hashtable uses a perfect hash function ensuring collision-free indexing of a subset of entries;
    inserting new entries might cause a collision requiring 
  - Reason 3: adding too many entries will eventually create an imbalance in the first-level hashtable;
    fixing this imbalance effectively requires re-constructing the file from scratch
- Cannot iterate over keys
  - Reason: compactindex stores hashes, not the entries themselves.
    This saves space but also allows for efficient random reads used during binary search
