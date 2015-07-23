package sparsevector

import (
	"sort"
)

type Value float32

type Vector interface {
	Cos(v Vector) Value
	Mag() Value
	Dot(v Vector) Value
	Mean() Value
	AddConst(toAdd Value)
	SubConst(toSub Value)
	GetIndex() VectorIndex
	Iter(f func(index interface{}, value Value))
	IterUpdate(f func(index interface{}, value Value) Value)
}

type VectorIndex interface {
	sort.Interface
	LessThanOther(i int, v2 VectorIndex, j int) bool
	GetAtLocation(i int) interface{}
}
