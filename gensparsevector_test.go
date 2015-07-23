package sparsevector

import (
	"math"
	"testing"
)

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
}
