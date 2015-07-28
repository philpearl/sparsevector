package sparsevector

// Uint32Index is an implementation of VectorIndex for uint32.
type Uint32Index []uint32

func (si Uint32Index) Len() int      { return len(si) }
func (si Uint32Index) Swap(i, j int) { si[i], si[j] = si[j], si[i] }
func (si Uint32Index) Less(i, j int) bool {
	return si[i] < si[j]
}

func (si Uint32Index) LessThanOther(i int, sii2 VectorIndex, j int) bool {
	si2 := sii2.(Uint32Index)
	return si[i] < si2[j]
}

func (si Uint32Index) GetAtLocation(location int) interface{} {
	return si[location]
}

// Assert Uint32Index implements VectorIndex
var _ VectorIndex = (Uint32Index)(nil)
