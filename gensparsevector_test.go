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
