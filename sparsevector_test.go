package sparsevector

import (
	"math/rand"
	"testing"
)

func BenchmarkSparseVector1000(b *testing.B) {
	benchmarkSparseVectorM(b, 1000)
}

func BenchmarkSparseVector10000(b *testing.B) {
	benchmarkSparseVectorM(b, 10000)
}

func benchmarkSparseVectorM(b *testing.B, m int) {
	// Generate 2 sparse vectors using 75% of numbers 1-m
	v1 := genRandomSparseVector(m)
	v2 := genRandomSparseVector(m)

	b.ReportAllocs()
	b.ResetTimer()

	var total Value
	for i := 0; i < b.N; i++ {
		result := v1.Cos(v2)
		total += result
	}
}

func genRandomSparseVector(m int) *SparseVectorUint32 {
	l := m * 3 / 4
	values := make([]Value, l)
	for i, v := range rand.Perm(m)[:l] {
		values[i] = Value(v)
	}
	index := make([]uint32, l)
	for i, v := range rand.Perm(m)[:l] {
		index[i] = uint32(v)
	}
	return NewSparseVectorUint32(index, values)
}
