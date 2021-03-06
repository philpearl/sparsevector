package sparsevector

import (
	"math"
	"sort"
)

// GenSparseVector is a sparse vector whose rows can be identified by any type that can
// be ordered. For example, the rows could be usernames, URLs, document names.
//
// Note that using this with a simple uint32 index carries an approxiamate 5x performance
// penalty relative to using a SparseVector. However it is faster than MapSparseVector for
// the vector lengths we have benchmarked
type GenSparseVector struct {
	index    VectorIndex
	values   []Value
	mag      Value
	magClean bool
}

// NewGenSparseVector creates a new GenSparseVector. You should provide parallel arrays of
// indicies and values. Note that NewGenSparseVector will sort these in place.
func NewGenSparseVector(index VectorIndex, values []Value) *GenSparseVector {
	// Want to sort the index and values at the same time
	v := &GenSparseVector{index: index, values: values}

	gsv := genSparseVectorSort{v}
	sort.Sort(gsv)
	return v
}

type genSparseVectorSort struct {
	*GenSparseVector
}

func (v genSparseVectorSort) Swap(i, j int) {
	v.index.Swap(i, j)
	v.values[i], v.values[j] = v.values[j], v.values[i]
}

func (v genSparseVectorSort) Len() int { return len(v.values) }

func (v genSparseVectorSort) Less(i, j int) bool { return v.index.Less(i, j) }

// Dot calculates the dot-product of this vector and another.
// Both vectors should be GenSparseVectors
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

// Mag returns the magnitude of this vector
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

// Cos calculates the cosine between this vector and another.
// Both vectors must be GenSparseVectors
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

func (sv1 *GenSparseVector) Add(v2 Vector) Vector {
	return sv1.runOp(v2, AddOp)
}

func (sv1 *GenSparseVector) Sub(v2 Vector) Vector {
	return sv1.runOp(v2, SubOp)
}

func (sv1 *GenSparseVector) runOp(v2 Vector, op ValueOp) Vector {
	sv2 := v2.(*GenSparseVector)

	var i1, i2 int
	sv1l := sv1.index.Len()
	sv2l := sv2.index.Len()

	// Create index and value array to back the result
	l := sv1l
	if l < sv2l {
		l = sv2l
	}
	oi := sv1.index.New(l)
	ov := make([]Value, 0, l)

	sv1_exp := false
	sv2_exp := false

	for {
		if i1 >= sv1l {
			sv1_exp = true // sv1 is exhausted
		}
		if i2 >= sv2l {
			sv2_exp = true // sv2 is exhausted
			if sv1_exp {
				break
			}
		}
		if sv2_exp || (!sv1_exp && sv1.index.LessThanOther(i1, sv2.index, i2)) {
			oi = oi.Append(sv1.index.GetAtLocation(i1))
			ov = append(ov, op(sv1.values[i1], 0))
			i1 += 1
		} else if sv1_exp || (!sv2_exp && sv2.index.LessThanOther(i2, sv1.index, i1)) {
			oi = oi.Append(sv2.index.GetAtLocation(i2))
			ov = append(ov, op(0, sv2.values[i2]))
			i2 += 1
		} else {
			oi = oi.Append(sv1.index.GetAtLocation(i1))
			ov = append(ov, op(sv1.values[i1], sv2.values[i2]))
			i1 += 1
			i2 += 1
		}
	}

	return &GenSparseVector{
		index:  oi,
		values: ov,
	}
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

// Mult multiplies the vector by a constant.
func (sv *GenSparseVector) Mult(l Value) {
	for i, v := range sv.values {
		sv.values[i] = l * v
	}
	sv.magClean = false
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
var _ SparseVector = (*GenSparseVector)(nil)
var _ IndexSparseVector = (*GenSparseVector)(nil)
