package sparsevector

import (
	"math"
	"sort"
)

// SparseVector is a sparse representation of a vector. It stores only non-zero
// values, and stores them in parallel arrays of indices and values. When the
// vector is created the arrays are sorted by index.
type SparseVector struct {
	indices []uint32
	values  []Value
	mag     Value
}

// NewSparseVector creates a new sparse vector. Pass in parallel arrays of the
// indices and values for non-zero entries. NewSparseVector will sort these
// entries by index to enable faster calculations later. It also calculates the
// vector magnitude.
func NewSparseVector(indices []uint32, values []Value) *SparseVector {
	sv := &SparseVector{
		indices: indices,
		values:  values,
	}
	sort.Sort(sv)
	sv.mag = sv.calcMag()
	return sv
}

// We implement a sort interface to order the elements by increasing index
func (sv *SparseVector) Len() int           { return len(sv.indices) }
func (sv *SparseVector) Less(i, j int) bool { return sv.indices[i] < sv.indices[j] }
func (sv *SparseVector) Swap(i, j int) {
	sv.indices[i], sv.indices[j] = sv.indices[j], sv.indices[i]
	sv.values[i], sv.values[j] = sv.values[j], sv.values[i]
}

// Dot calculates the dot product of this vector and another sparse vector.
func (sv1 *SparseVector) Dot(sv2 *SparseVector) Value {

	var i1, i2 int
	var dp Value
	sv1l := len(sv1.indices)
	sv2l := len(sv2.indices)
	for i1 < sv1l && i2 < sv2l {
		if sv1.indices[i1] < sv2.indices[i2] {
			i1 += 1
		} else if sv2.indices[i2] < sv1.indices[i1] {
			i2 += 1
		} else {
			dp += sv1.values[i1] * sv2.values[i2]
			i1 += 1
			i2 += 1
		}
	}
	return dp
}

// Cos calculates the cosine of the angle between this sparse vector
// and another.
func (sv1 *SparseVector) Cos(sv2 *SparseVector) Value {
	return sv1.Dot(sv2) / (sv1.mag * sv2.mag)
}

// Mean() Calculates the mean element value (mean of values that are present)
func (sv *SparseVector) Mean() Value {
	var total Value
	for _, v := range sv.values {
		total += v
	}
	return total / Value(len(sv.values))
}

// AddConst adds a constant value to each of the present values in the sparse
// vector
func (sv *SparseVector) AddConst(toAdd Value) {
	for i, v := range sv.values {
		sv.values[i] = v + toAdd
	}
	sv.mag = sv.calcMag()
}

// SubConst subtracts a constant value to each of the present values in the sparse
// vector.
func (sv *SparseVector) SubConst(toSub Value) {
	sv.AddConst(-toSub)
}

// Iter lets you iterate over the members of the sparse vector
func (sv *SparseVector) Iter(f func(index uint32, value Value)) {
	for i, index := range sv.indices {
		value := sv.values[i]
		f(index, value)
	}
}

// Iter lets you iterate over the members of the sparse vector
func (sv *SparseVector) IterUpdate(f func(index uint32, value Value) Value) {
	for i, index := range sv.indices {
		value := sv.values[i]
		sv.values[i] = f(index, value)
	}
	sv.mag = sv.calcMag()
}

func (sv *SparseVector) calcMag() Value {
	var magsq Value
	for _, val := range sv.values {
		magsq += val * val
	}
	return Value(math.Sqrt(float64(magsq)))
}

// GetIndices returns the array of indices of non-zero values in the vector
func (sv *SparseVector) GetIndices() []uint32 { return sv.indices }

// MapIndices applies a map to the indices of non-zero values in the vector.
// The idea is to allow you to shuffle non-zero values to the start of a set
// of vectors, potentially reducing the lengths of vectors in the set
func (sv *SparseVector) MapIndices(im map[uint32]uint32) {
	for i, idx := range sv.indices {
		sv.indices[i] = im[idx]
	}
}
