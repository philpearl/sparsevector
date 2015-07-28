package sparsevector

import (
	"sort"
)

// Value is the type we use for vector values. We defined a type so we can easily
// change types, until such times as our Go overlords figure out generics.
type Value float32

// Vector is a general interface to a vector. Our focus on similarity calculations
// means that so far we've only defined cosine, magnitude and dot-product methods.
type Vector interface {
	// Cos calculates the cosine of the angle between two vectors
	Cos(v Vector) Value

	// Mag gives the magnitude of a vector
	Mag() Value

	// Dot calculates the dot product of two vectors
	Dot(v Vector) Value
}

// SparseVector extends the Vector interface to include some functions that only
// apply to sparse vectors, in particular ones the operate only on present values.
// The intention of Mean, AddConst and SubConst is to allow you to mean-center
// your vectors.
type SparseVector interface {
	Vector

	// Mean calculates the mean of non-zero values in a vector
	Mean() Value

	// AddConst adds a constant value to non-zero values in the vector
	AddConst(toAdd Value)

	// SubConst subtracts a constant value from non-zero values in the vector
	SubConst(toSub Value)
}

// IndexSparseVector is an additional interface for sparse vectors that use
// the VectorIndex interface to describe their indices
type IndexSparseVector interface {
	// GetIndex returns the index values of the sparse vector
	GetIndex() VectorIndex

	// Iter calls a function for each value present in a sparse vector
	Iter(f func(index interface{}, value Value))

	// IterUpdate calls a function for each value present in a sparse vector. The
	// value is replaced by the result of the function
	IterUpdate(f func(index interface{}, value Value) Value)
}

// VectorIndex is an interface for the index of a sparse vector.
// The idea here is to allow the index and values to be sorted in parallel.
// The mental model for this interface is that the index is implemented
// as an array or slice, and that index values at points in the array can be
// extracted or compared.
type VectorIndex interface {
	// VectorIndex must implement the standard sort interface
	sort.Interface

	// LessThanOther compares values at a location in this index with
	// a value at a second location in another index.
	// In most cases the two VectorIndexs must be the same underlying
	// type.
	LessThanOther(i int, v2 VectorIndex, j int) bool

	// GetAtLocation returns the index value currently at location i
	GetAtLocation(i int) interface{}
}
