package sparsevector

import (
	"math/rand"
	"reflect"
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

func TestAddSparseVectorUint32(t *testing.T) {

	tests := []struct {
		v1  *SparseVectorUint32
		v2  *SparseVectorUint32
		exp *SparseVectorUint32
	}{{
		v1:  NewSparseVectorUint32([]uint32{1, 2, 3}, []Value{4, 5, 6}),
		v2:  NewSparseVectorUint32([]uint32{1, 3, 4}, []Value{4, 5, 6}),
		exp: NewSparseVectorUint32([]uint32{1, 2, 3, 4}, []Value{8, 5, 11, 6}),
	}, {
		v1:  NewSparseVectorUint32([]uint32{}, []Value{}),
		v2:  NewSparseVectorUint32([]uint32{1, 3, 4}, []Value{4, 5, 6}),
		exp: NewSparseVectorUint32([]uint32{1, 3, 4}, []Value{4, 5, 6}),
	}, {
		v1:  NewSparseVectorUint32([]uint32{2}, []Value{7}),
		v2:  NewSparseVectorUint32([]uint32{1, 3, 4}, []Value{4, 5, 6}),
		exp: NewSparseVectorUint32([]uint32{1, 2, 3, 4}, []Value{4, 7, 5, 6}),
	}, {
		v1:  NewSparseVectorUint32([]uint32{1, 3, 4}, []Value{4, 5, 6}),
		v2:  NewSparseVectorUint32([]uint32{1, 3, 4}, []Value{4, 5, 6}),
		exp: NewSparseVectorUint32([]uint32{1, 3, 4}, []Value{8, 10, 12}),
	}, {
		v1:  NewSparseVectorUint32([]uint32{1, 3, 4}, []Value{4, 5, 6}),
		v2:  NewSparseVectorUint32([]uint32{1, 3}, []Value{4, 5}),
		exp: NewSparseVectorUint32([]uint32{1, 3, 4}, []Value{8, 10, 6}),
	}, {
		v1:  NewSparseVectorUint32([]uint32{1, 3, 4}, []Value{4, 5, 6}),
		v2:  NewSparseVectorUint32([]uint32{}, []Value{}),
		exp: NewSparseVectorUint32([]uint32{1, 3, 4}, []Value{4, 5, 6}),
	},
	}

	for i, test := range tests {
		v3 := test.v1.Add(test.v2)
		v4 := test.v2.Add(test.v1)

		if !reflect.DeepEqual(test.exp, v3) {
			t.Errorf("Test %d. Sum not as expected. Have %v", i, v3)
		}
		if !reflect.DeepEqual(test.exp, v4) {
			t.Errorf("Test %d. Reverse Sum not as expected. Have %v", i, v4)
		}

	}
}
