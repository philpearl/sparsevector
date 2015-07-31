# SparseVector

Efficient sparse vector implementations in Go.

[![Build Status](https://travis-ci.org/philpearl/sparsevector.svg)](https://travis-ci.org/philpearl/sparsevector) [![GoDoc](https://godoc.org/github.com/philpearl/sparsevector?status.svg)](https://godoc.org/github.com/philpearl/sparsevector)

## What's included?

|    |    |
| ---| ---|
| SparseVectorUint32 | Sparse Vector with uint32 indices and Value values implemented by parallel ordered lists of indices and values |
| GenSparseVector | Sparse Vector with generic indices and a parallel ordered list of values. The index must implement the VectorIndex interface, and hence be sortable. |
| MapSparseVector | A Sparse Vector with uint32 indices and Value values implemented using a map |
| Uint32Index | a GenSparseVector index for uint32 |
| IntIndex | A GenSparseVector index for int |
| StringIndex | A GenSparseVector index for strings |

## Performance

Benchmarks are included for uint32 versions of all the Sparse Vector implementations. In these tests SparseVectorUint32 is by far the fastest, GenSparseVector takes about 5 times as long, and MapSparseVector takes about 1.6 times more again.

## What can you do with it?

I've focused on what I need for similarity calculations, so the vectors do cosine and dot-product. I've also included adding and subtracting vectors and constant values, and multiplying by constant values. You can discover the mean of the present values, and also iterate and perform operations on the elements present in the vectors.

## License

MIT license in LICENSE.txt

## Contributing

Drop a pull request if you'd like to contribute. 
