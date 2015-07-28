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
