package sparsevector

// IntIndex is an implementation of VectorIndex for int.
type IntIndex []int

func (si IntIndex) Len() int      { return len(si) }
func (si IntIndex) Swap(i, j int) { si[i], si[j] = si[j], si[i] }
func (si IntIndex) Less(i, j int) bool {
	return si[i] < si[j]
}

func (si IntIndex) LessThanOther(i int, sii2 VectorIndex, j int) bool {
	si2 := sii2.(IntIndex)
	return si[i] < si2[j]
}

func (si IntIndex) GetAtLocation(location int) interface{} {
	return si[location]
}

// Assert IntIndex implements VectorIndex
var _ VectorIndex = (IntIndex)(nil)
