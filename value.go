package sparsevector

// Value is the type we use for vector values. We defined a type so we can easily
// change types, until such times as our Go overlords figure out generics.
type Value float32

// An operation on a pair of values.
type ValueOp func(v1, v2 Value) Value

func AddOp(v1, v2 Value) Value { return v1 + v2 }
func SubOp(v1, v2 Value) Value { return v1 - v2 }
