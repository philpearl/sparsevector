package sparsevector

// StringIndex is an implementation of VectorIndex for string.
type StringIndex []string

func (si StringIndex) Len() int      { return len(si) }
func (si StringIndex) Swap(i, j int) { si[i], si[j] = si[j], si[i] }
func (si StringIndex) Less(i, j int) bool {
	return si[i] < si[j]
}

func (si StringIndex) LessThanOther(i int, sii2 VectorIndex, j int) bool {
	si2 := sii2.(StringIndex)
	return si[i] < si2[j]
}

func (si StringIndex) GetAtLocation(location int) interface{} {
	return si[location]
}

func (si StringIndex) New(l int) VectorIndex {
	return make(StringIndex, 0, l)
}

func (si StringIndex) Append(idx interface{}) VectorIndex {
	return append(si, idx.(string))
}

// Assert StringIndex implements VectorIndex
var _ VectorIndex = (StringIndex)(nil)
