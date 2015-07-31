package sparsevector

import (
	"math/rand"
	"reflect"
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

func TestAddMapSparseVector(t *testing.T) {

	tests := []struct {
		v1  *MapSparseVector
		v2  *MapSparseVector
		exp *MapSparseVector
	}{{
		v1:  NewMapSparseVector([]uint32{1, 2, 3}, []Value{4, 5, 6}),
		v2:  NewMapSparseVector([]uint32{1, 3, 4}, []Value{4, 5, 6}),
		exp: NewMapSparseVector([]uint32{1, 2, 3, 4}, []Value{8, 5, 11, 6}),
	}, {
		v1:  NewMapSparseVector([]uint32{}, []Value{}),
		v2:  NewMapSparseVector([]uint32{1, 3, 4}, []Value{4, 5, 6}),
		exp: NewMapSparseVector([]uint32{1, 3, 4}, []Value{4, 5, 6}),
	}, {
		v1:  NewMapSparseVector([]uint32{2}, []Value{7}),
		v2:  NewMapSparseVector([]uint32{1, 3, 4}, []Value{4, 5, 6}),
		exp: NewMapSparseVector([]uint32{1, 2, 3, 4}, []Value{4, 7, 5, 6}),
	}, {
		v1:  NewMapSparseVector([]uint32{1, 3, 4}, []Value{4, 5, 6}),
		v2:  NewMapSparseVector([]uint32{1, 3, 4}, []Value{4, 5, 6}),
		exp: NewMapSparseVector([]uint32{1, 3, 4}, []Value{8, 10, 12}),
	}, {
		v1:  NewMapSparseVector([]uint32{1, 3, 4}, []Value{4, 5, 6}),
		v2:  NewMapSparseVector([]uint32{1, 3}, []Value{4, 5}),
		exp: NewMapSparseVector([]uint32{1, 3, 4}, []Value{8, 10, 6}),
	}, {
		v1:  NewMapSparseVector([]uint32{1, 3, 4}, []Value{4, 5, 6}),
		v2:  NewMapSparseVector([]uint32{}, []Value{}),
		exp: NewMapSparseVector([]uint32{1, 3, 4}, []Value{4, 5, 6}),
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

func TestMultMapSparseVector(t *testing.T) {
	tests := []struct {
		v   *MapSparseVector
		l   Value
		exp *MapSparseVector
	}{{
		v:   NewMapSparseVector([]uint32{1, 2, 3}, []Value{4, 5, 6}),
		l:   1,
		exp: NewMapSparseVector([]uint32{1, 2, 3}, []Value{4, 5, 6}),
	}, {
		v:   NewMapSparseVector([]uint32{}, []Value{}),
		l:   1,
		exp: NewMapSparseVector([]uint32{}, []Value{}),
	}, {
		v:   NewMapSparseVector([]uint32{2}, []Value{7}),
		l:   2,
		exp: NewMapSparseVector([]uint32{2}, []Value{14}),
	}, {
		v:   NewMapSparseVector([]uint32{1, 3, 4}, []Value{4, 5, 6}),
		l:   0,
		exp: NewMapSparseVector([]uint32{1, 3, 4}, []Value{0, 0, 0}),
	}, {
		v:   NewMapSparseVector([]uint32{1, 3, 4}, []Value{4, 5, 6}),
		l:   -1,
		exp: NewMapSparseVector([]uint32{1, 3, 4}, []Value{-4, -5, -6}),
	}, {
		v:   NewMapSparseVector([]uint32{1, 3, 4}, []Value{4, 5, 6}),
		l:   137.4,
		exp: NewMapSparseVector([]uint32{1, 3, 4}, []Value{549.6, 687, 824.39996}),
	},
	}

	for i, test := range tests {
		test.v.Mult(test.l)

		if !reflect.DeepEqual(test.exp, test.v) {
			t.Errorf("Test %d. Sum not as expected. Have %v", i, test.v)
		}
	}
}
