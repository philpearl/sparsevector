package sparsevector

import (
	"math"
	"sort"
)

type GenSparseVector struct {
	index    VectorIndex
	values   []Value
	mag      Value
	magClean bool
}

func NewGenSparseVector(index VectorIndex, values []Value) *GenSparseVector {
	// Want to sort the index and values at the same time
	v := &GenSparseVector{index: index, values: values}
	sort.Sort(v)
	return v
}

func (v *GenSparseVector) Swap(i, j int) {
	v.index.Swap(i, j)
	v.values[i], v.values[j] = v.values[j], v.values[i]
}

func (v *GenSparseVector) Len() int { return len(v.values) }

func (v *GenSparseVector) Less(i, j int) bool { return v.index.Less(i, j) }

func (sv1 *GenSparseVector) Dot(svi2 Vector) Value {
	sv2 := svi2.(*GenSparseVector)

	var i1, i2 int
	var dp Value
	sv1l := sv1.index.Len()
	sv2l := sv2.index.Len()
	for i1 < sv1l && i2 < sv2l {
		if sv1.index.LessThanOther(i1, sv2.index, i2) {
			i1 += 1
		} else if sv2.index.LessThanOther(i2, sv1.index, i1) {
			i2 += 1
		} else {
			dp += sv1.values[i1] * sv2.values[i2]
			i1 += 1
			i2 += 1
		}
	}
	return dp
}

func (v *GenSparseVector) Mag() Value {
	if !v.magClean {
		// Could use v1.Dot(v2), but this is more efficient
		var magsq Value
		for _, val := range v.values {
			magsq += val * val
		}
		v.mag = Value(math.Sqrt(float64(magsq)))

		v.magClean = true
	}
	return v.mag
}

func (sv1 *GenSparseVector) Cos(svi2 Vector) Value {
	sv2 := svi2.(*GenSparseVector)
	return sv1.Dot(sv2) / (sv1.Mag() * sv2.Mag())
}

// Mean() Calculates the mean element value (mean of values that are present)
func (sv *GenSparseVector) Mean() Value {
	var total Value
	for _, v := range sv.values {
		total += v
	}
	return total / Value(len(sv.values))
}

// AddConst adds a constant value to each of the present values in the sparse
// vector
func (sv *GenSparseVector) AddConst(toAdd Value) {
	for i, v := range sv.values {
		sv.values[i] = v + toAdd
	}
	sv.magClean = false
}

// SubConst subtracts a constant value to each of the present values in the sparse
// vector.
func (sv *GenSparseVector) SubConst(toSub Value) {
	sv.AddConst(-toSub)
}

// Iter lets you iterate over the members of the sparse vector
func (sv *GenSparseVector) Iter(f func(index interface{}, value Value)) {
	for i, value := range sv.values {
		index := sv.index.GetAtLocation(i)
		f(index, value)
	}
}

// Iter lets you iterate over the members of the sparse vector
func (sv *GenSparseVector) IterUpdate(f func(index interface{}, value Value) Value) {
	for i, value := range sv.values {
		index := sv.index.GetAtLocation(i)
		sv.values[i] = f(index, value)
	}
	sv.magClean = false
}

// GetIndices returns the index values of the vector
func (sv *GenSparseVector) GetIndex() VectorIndex { return sv.index }

// Assert GenSparseVector implements Vector
var _ Vector = (*GenSparseVector)(nil)
