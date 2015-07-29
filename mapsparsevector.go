package sparsevector

import (
	"math"
)

// MapSparseVector is a sparse vector implemented using a map
type MapSparseVector struct {
	values   map[uint32]Value
	mag      Value
	magClean bool
}

var _ Vector = (*MapSparseVector)(nil)

// NewMapSparseVector creates a new MapSparseVector.
// Pass in parallel arrays of the indices and values of non-zero entries in the vector
func NewMapSparseVector(indices []uint32, values []Value) *MapSparseVector {
	m := make(map[uint32]Value, len(indices))
	for i, index := range indices {
		m[index] = values[i]
	}

	return &MapSparseVector{
		values: m,
	}
}

// Mag returns the magnitude of the vector
func (m *MapSparseVector) Mag() Value {
	if !m.magClean {
		var mag Value
		for _, v := range m.values {
			mag += v * v
		}
		m.mag = Value(math.Sqrt(float64(mag)))
		m.magClean = true
	}
	return m.mag
}

// Dot calculates the dot product of this sparse vector and another.
// Both vectors must be MapSparseVectors
func (m1 *MapSparseVector) Dot(v2 Vector) Value {
	m2 := v2.(*MapSparseVector)
	if len(m2.values) < len(m1.values) {
		m1, m2 = m2, m1
	}

	var dp Value
	for i, v := range m1.values {
		v2, ok := m2.values[i]
		if ok {
			dp += v * v2
		}
	}
	return dp
}

// Cos calculates the cosine of this vector and another.
// Both vectors must be MapSparseVectors
func (m1 *MapSparseVector) Cos(v2 Vector) Value {
	return m1.Dot(v2) / (m1.Mag() * v2.Mag())
}

func (m1 *MapSparseVector) Add(v2 Vector) Vector {
	return m1.runOp(v2, AddOp)
}

func (m1 *MapSparseVector) Sub(v2 Vector) Vector {
	return m1.runOp(v2, SubOp)
}

func (m1 *MapSparseVector) runOp(v2 Vector, op ValueOp) Vector {
	m2 := v2.(*MapSparseVector)

	// Build a map to back the output vector.
	// It should be at least as big as the biggest of our
	// two vectors
	l := len(m1.values)
	if l < len(m2.values) {
		l = len(m2.values)
	}
	om := make(map[uint32]Value, l)

	// Copy all the values from m1
	for k, v := range m1.values {
		om[k] = v
	}

	// Add values from m2
	for k, v := range m2.values {
		ev, ok := m1.values[k]
		if ok {
			om[k] = op(ev, v)
		} else {
			om[k] = op(0, v)
		}
	}

	return &MapSparseVector{
		values: om,
	}
}
