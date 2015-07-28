package sparsevector

import (
	"math/rand"
	"testing"
)

func BenchmarkMapSparseVector1000(b *testing.B) {
	benchmarkMapSparseVectorM(b, 1000)
}

func BenchmarkMapSparseVector10000(b *testing.B) {
	benchmarkMapSparseVectorM(b, 10000)
}

func benchmarkMapSparseVectorM(b *testing.B, m int) {
	// Generate 2 sparse vectors using 75% of numbers 1-m
	v1 := genRandomMapSparseVector(m)
	v2 := genRandomMapSparseVector(m)

	b.ReportAllocs()
	b.ResetTimer()

	var total Value
	for i := 0; i < b.N; i++ {
		result := v1.Cos(v2)
		total += result
	}
}

func genRandomMapSparseVector(m int) *MapSparseVector {
	l := m * 3 / 4
	values := make([]Value, l)
	for i, v := range rand.Perm(m)[:l] {
		values[i] = Value(v)
	}
	index := make([]uint32, l)
	for i, v := range rand.Perm(m)[:l] {
		index[i] = uint32(v)
	}
	return NewMapSparseVector(index, values)
}
