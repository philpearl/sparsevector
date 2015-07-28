// sparsevector contains a number of implementations of SparseVectors, with a focus on
// dot products and cosines.
//
// The fastest implementation is SparseVectorUint32. This works by having parallel arrays
// for the vector element indices and values.  The arrays are sorted in parallel when the
// vector is created to allow for a linear scan through the vectors when calculating Cos()
// and Dot(). The indices and values are separate to optimise scanning through the indices.
//
// Next in performance is GenSparseVector. This is similar to SparseVectorUint32 but the index can be any
// type for which you can implement VectorIndex. It works similarly to SparseVectorUint32 but
// is slowed down by accessing index values via the VectorIndex interface.
//
// The slowest implementation is MapSparseVector. This is implemented using a map[uint32]Value.
// I view this as the baseline implementation. Since performance depends very much on your
// data and what you're doing with it I'd advise testing with each of the implementations to
// find which is fastest for you.
package sparsevector

import (
	"math"
	"sort"
)

// SparseVector is a sparse representation of a vector. It stores only non-zero
// values, and stores them in parallel arrays of indices and values. When the
// vector is created the arrays are sorted by index.
type SparseVectorUint32 struct {
	indices  []uint32
	values   []Value
	mag      Value
	magClean bool
}

// NewSparseVector creates a new sparse vector. Pass in parallel arrays of the
// indices and values for non-zero entries. NewSparseVector will sort these
// entries by index to enable faster calculations later. It also calculates the
// vector magnitude.
func NewSparseVectorUint32(indices []uint32, values []Value) *SparseVectorUint32 {
	sv := &SparseVectorUint32{
		indices: indices,
		values:  values,
	}
	svs := sparseVectorUint32Sort{sv}
	sort.Sort(svs)
	return sv
}

type sparseVectorUint32Sort struct {
	*SparseVectorUint32
}

// We implement a sort interface to order the elements by increasing index
func (sv sparseVectorUint32Sort) Len() int           { return len(sv.indices) }
func (sv sparseVectorUint32Sort) Less(i, j int) bool { return sv.indices[i] < sv.indices[j] }
func (sv sparseVectorUint32Sort) Swap(i, j int) {
	sv.indices[i], sv.indices[j] = sv.indices[j], sv.indices[i]
	sv.values[i], sv.values[j] = sv.values[j], sv.values[i]
}

// Mag returns the magnitude of the vector. It is calculated lazily and cached.
// Note (as with the other sparse vector implementations) this means these vectors
// are not thread safe.
func (v *SparseVectorUint32) Mag() Value {
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

// Dot calculates the dot product of this vector and another sparse vector.
func (sv1 *SparseVectorUint32) Dot(sv2in Vector) Value {
	sv2 := sv2in.(*SparseVectorUint32)

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

func (sv1 *SparseVectorUint32) Add(sv2in Vector) *SparseVectorUint32 {
	sv2 := sv2in.(*SparseVectorUint32)

	var i1, i2 int
	sv1l := len(sv1.indices)
	sv2l := len(sv2.indices)

	// Our output vectors are at least as long as our longest input
	l := sv1l
	if sv2l > l {
		l = sv2l
	}
	oi := make([]uint32, 0, l)
	ov := make([]Value, 0, l)

	for {
		var sv1i, sv2i uint32
		if i1 < sv1l {
			sv1i = sv1.indices[i1]
		} else {
			// sv1 exhausted - make sure we take from sv2
			sv1i = math.MaxUint32
		}
		if i2 < sv2l {
			sv2i = sv2.indices[i2]
		} else {
			// sv2 exhausted - make sure we take from sv1
			sv2i = math.MaxUint32
			if sv1i == math.MaxUint32 {
				break // both exhausted
			}
		}

		if sv1i < sv2i {
			oi = append(oi, sv1i)
			ov = append(ov, sv1.values[i1])
			i1 += 1

		} else if sv2i < sv1i {
			oi = append(oi, sv2i)
			ov = append(ov, sv2.values[i2])
			i2 += 1

		} else {
			oi = append(oi, sv2i)
			ov = append(ov, sv2.values[i2]+sv1.values[i1])
			i1 += 1
			i2 += 1
		}
	}
	// The vector should already be sorted
	return &SparseVectorUint32{
		indices: oi,
		values:  ov,
	}
}

// Cos calculates the cosine of the angle between this sparse vector
// and another.
func (sv1 *SparseVectorUint32) Cos(sv2 Vector) Value {
	return sv1.Dot(sv2) / (sv1.Mag() * sv2.Mag())
}

// Mean() Calculates the mean element value (mean of values that are present)
func (sv *SparseVectorUint32) Mean() Value {
	var total Value
	for _, v := range sv.values {
		total += v
	}
	return total / Value(len(sv.values))
}

// AddConst adds a constant value to each of the present values in the sparse
// vector
func (sv *SparseVectorUint32) AddConst(toAdd Value) {
	for i, v := range sv.values {
		sv.values[i] = v + toAdd
	}
	sv.magClean = false
}

// SubConst subtracts a constant value to each of the present values in the sparse
// vector.
func (sv *SparseVectorUint32) SubConst(toSub Value) {
	sv.AddConst(-toSub)
}

// Iter lets you iterate over the members of the sparse vector
func (sv *SparseVectorUint32) Iter(f func(index uint32, value Value)) {
	for i, index := range sv.indices {
		value := sv.values[i]
		f(index, value)
	}
}

// Iter lets you iterate over the members of the sparse vector
func (sv *SparseVectorUint32) IterUpdate(f func(index uint32, value Value) Value) {
	for i, index := range sv.indices {
		value := sv.values[i]
		sv.values[i] = f(index, value)
	}
	sv.magClean = false
}

// GetIndices returns the array of indices of non-zero values in the vector
func (sv *SparseVectorUint32) GetIndices() []uint32 { return sv.indices }

// MapIndices applies a map to the indices of non-zero values in the vector.
// The idea is to allow you to shuffle non-zero values to the start of a set
// of vectors, potentially reducing the lengths of vectors in the set
func (sv *SparseVectorUint32) MapIndices(im map[uint32]uint32) {
	for i, idx := range sv.indices {
		sv.indices[i] = im[idx]
	}
}

var _ SparseVector = (*SparseVectorUint32)(nil)
