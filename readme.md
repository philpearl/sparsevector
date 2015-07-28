# SparseVector

Efficient sparse vector implemenations in Go.

[![Build Status](https://travis-ci.org/philpearl/sparsevector.svg)](https://travis-ci.org/philpearl/sparsevector) [![GoDoc](https://godoc.org/github.com/philpearl/sparsevector?status.svg)](https://godoc.org/github.com/philpearl/sparsevector)

## What's included?

| SparseVectorUint32 | Sparse Vector with uint32 indices and Value values implemented by parallel ordered lists of indices and values |
| GenSparseVector | Sparse Vector with generic indices and a parallel ordered list of values. The index must implement the VectorIndex interface, and hence be sortable. |
| MapSparseVector | A Sparse Vector with uint32 indices and Value values implemented using a map |
| Uint32Index | a GenSparseVector index for uint32 |
| IntIndex | A GenSparseVector index for int |
| StringIndex | A GenSparseVector index for strings |

## Performance

Benchmarks are included for uint32 versions of all the Sparse Vector implementations. In these tests SparseVectorUint32 is by far the fastest, GenSparseVector takes about 5 times as long, and MapSparseVector takes about 1.6 times more again.

## What can you do with it?

Vector stuff.  You know. Like with vectors.
