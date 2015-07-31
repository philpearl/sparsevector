package sparsevector

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
)

func BenchmarkGenSparseVectorInt1000(b *testing.B) {
	benchmarkGenSparseVectorIntM(b, 1000)
}

func BenchmarkGenSparseVectorInt10000(b *testing.B) {
	benchmarkGenSparseVectorIntM(b, 10000)
}

func benchmarkGenSparseVectorIntM(b *testing.B, m int) {
	// Generate 2 sparse vectors using 75% of numbers 1-m
	v1 := genRandomGenSparseVectorInt(m)
	v2 := genRandomGenSparseVectorInt(m)

	b.ReportAllocs()
	b.ResetTimer()

	var total Value
	for i := 0; i < b.N; i++ {
		result := v1.Cos(v2)
		total += result
	}
}

func genRandomGenSparseVectorInt(m int) *GenSparseVector {
	l := m * 3 / 4
	values := make([]Value, l)
	for i, v := range rand.Perm(m)[:l] {
		values[i] = Value(v)
	}
	return NewGenSparseVector(IntIndex(rand.Perm(m)[:l]), values)
}

func BenchmarkGenSparseVectorUint321000(b *testing.B) {
	benchmarkGenSparseVectorUint32M(b, 1000)
}

func BenchmarkGenSparseVectorUint3210000(b *testing.B) {
	benchmarkGenSparseVectorUint32M(b, 10000)
}

func benchmarkGenSparseVectorUint32M(b *testing.B, m int) {
	// Generate 2 sparse vectors using 75% of numbers 1-m
	v1 := genRandomGenSparseVectorUint32(m)
	v2 := genRandomGenSparseVectorUint32(m)

	b.ReportAllocs()
	b.ResetTimer()

	var total Value
	for i := 0; i < b.N; i++ {
		result := v1.Cos(v2)
		total += result
	}
}

func genRandomGenSparseVectorUint32(m int) *GenSparseVector {
	l := m * 3 / 4
	values := make([]Value, l)
	for i, v := range rand.Perm(m)[:l] {
		values[i] = Value(v)
	}
	index := make(Uint32Index, l)
	for i, v := range rand.Perm(m)[:l] {
		index[i] = uint32(v)
	}

	return NewGenSparseVector(index, values)
}

func TestGenSparseVector(t *testing.T) {

	si1 := StringIndex{"a", "b", "c"}
	v1 := []Value{1, 2, 3}

	sv1 := NewGenSparseVector(si1, v1)

	if sv1.Mag() != Value(math.Sqrt(1+4+9)) {
		t.Fatalf("Mag wrong - have %f", sv1.Mag())
	}

	si2 := StringIndex{"b", "d", "a"}
	v2 := []Value{1, 2, 7}
	sv2 := NewGenSparseVector(si2, v2)

	dot := sv1.Dot(sv2)
	if dot != 7+2 {
		t.Fatalf("Dot not as expected. Have %f", dot)
	}

	cos := sv1.Cos(sv2)
	if cos != Value((7+2)/(math.Sqrt(1+4+9)*math.Sqrt(1+4+49))) {
		t.Fatalf("cos wrong. Have %f", cos)
	}

	mean := sv1.Mean()
	if mean != 2 {
		t.Fatalf("Mean not as expected, have %f", mean)
	}

	sv1.AddConst(2)
	if sv1.Mag() != 7.071068 {
		t.Fatalf("Mag (2) not as expected, have %f", sv1.Mag())
	}
	mean = sv1.Mean()
	if mean != 4 {
		t.Fatalf("Mean (2) not as expected. Have %f", mean)
	}

	sv1.SubConst(2)
	mean = sv1.Mean()
	if mean != 2 {
		t.Fatalf("Mean (3) not as expected. Have %f", mean)
	}

	vals := make([]Value, 0, 3)
	sv1.Iter(func(index interface{}, value Value) {
		vals = append(vals, value)
	})
	if !reflect.DeepEqual(vals, v1) {
		t.Fatalf("iter values not as expected. Have %v", vals)
	}

	sv1.IterUpdate(func(index interface{}, value Value) Value {
		return 3
	})
	if sv1.Mean() != 3 {
		t.Fatalf("that didn't work")
	}
}

func TestAddGenUint32SparseVector(t *testing.T) {

	tests := []struct {
		v1  *GenSparseVector
		v2  *GenSparseVector
		exp *GenSparseVector
	}{{
		v1:  NewGenSparseVector(Uint32Index{1, 2, 3}, []Value{4, 5, 6}),
		v2:  NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{4, 5, 6}),
		exp: NewGenSparseVector(Uint32Index{1, 2, 3, 4}, []Value{8, 5, 11, 6}),
	}, {
		v1:  NewGenSparseVector(Uint32Index{}, []Value{}),
		v2:  NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{4, 5, 6}),
		exp: NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{4, 5, 6}),
	}, {
		v1:  NewGenSparseVector(Uint32Index{2}, []Value{7}),
		v2:  NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{4, 5, 6}),
		exp: NewGenSparseVector(Uint32Index{1, 2, 3, 4}, []Value{4, 7, 5, 6}),
	}, {
		v1:  NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{4, 5, 6}),
		v2:  NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{4, 5, 6}),
		exp: NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{8, 10, 12}),
	}, {
		v1:  NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{4, 5, 6}),
		v2:  NewGenSparseVector(Uint32Index{1, 3}, []Value{4, 5}),
		exp: NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{8, 10, 6}),
	}, {
		v1:  NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{4, 5, 6}),
		v2:  NewGenSparseVector(Uint32Index{}, []Value{}),
		exp: NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{4, 5, 6}),
	},
	}

	for i, test := range tests {
		t.Logf("test %d", i)
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

func TestSubGenUint32SparseVector(t *testing.T) {

	tests := []struct {
		v1  *GenSparseVector
		v2  *GenSparseVector
		exp *GenSparseVector
	}{{
		v1:  NewGenSparseVector(Uint32Index{1, 2, 3}, []Value{4, 5, 6}),
		v2:  NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{4, 5, 6}),
		exp: NewGenSparseVector(Uint32Index{1, 2, 3, 4}, []Value{0, 5, 1, -6}),
	}, {
		v1:  NewGenSparseVector(Uint32Index{}, []Value{}),
		v2:  NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{4, 5, 6}),
		exp: NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{-4, -5, -6}),
	}, {
		v1:  NewGenSparseVector(Uint32Index{2}, []Value{7}),
		v2:  NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{4, 5, 6}),
		exp: NewGenSparseVector(Uint32Index{1, 2, 3, 4}, []Value{-4, 7, -5, -6}),
	}, {
		v1:  NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{4, 5, 6}),
		v2:  NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{4, 5, 6}),
		exp: NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{0, 0, 0}),
	}, {
		v1:  NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{4, 5, 6}),
		v2:  NewGenSparseVector(Uint32Index{1, 3}, []Value{4, 5}),
		exp: NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{0, 0, 6}),
	}, {
		v1:  NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{4, 5, 6}),
		v2:  NewGenSparseVector(Uint32Index{}, []Value{}),
		exp: NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{4, 5, 6}),
	},
	}

	for i, test := range tests {
		t.Logf("test %d", i)
		v3 := test.v1.Sub(test.v2)

		if !reflect.DeepEqual(test.exp, v3) {
			t.Errorf("Test %d. Sub not as expected. Have %v", i, v3)
		}
	}
}

func TestMultGenUint32SparseVector(t *testing.T) {
	tests := []struct {
		v   *GenSparseVector
		l   Value
		exp *GenSparseVector
	}{{
		v:   NewGenSparseVector(Uint32Index{1, 2, 3}, []Value{4, 5, 6}),
		l:   1,
		exp: NewGenSparseVector(Uint32Index{1, 2, 3}, []Value{4, 5, 6}),
	}, {
		v:   NewGenSparseVector(Uint32Index{}, []Value{}),
		l:   1,
		exp: NewGenSparseVector(Uint32Index{}, []Value{}),
	}, {
		v:   NewGenSparseVector(Uint32Index{2}, []Value{7}),
		l:   2,
		exp: NewGenSparseVector(Uint32Index{2}, []Value{14}),
	}, {
		v:   NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{4, 5, 6}),
		l:   0,
		exp: NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{0, 0, 0}),
	}, {
		v:   NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{4, 5, 6}),
		l:   -1,
		exp: NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{-4, -5, -6}),
	}, {
		v:   NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{4, 5, 6}),
		l:   137.4,
		exp: NewGenSparseVector(Uint32Index{1, 3, 4}, []Value{549.6, 687, 824.39996}),
	},
	}

	for i, test := range tests {
		test.v.Mult(test.l)

		if !reflect.DeepEqual(test.exp, test.v) {
			t.Errorf("Test %d. Sum not as expected. Have %v", i, test.v)
		}
	}
}
