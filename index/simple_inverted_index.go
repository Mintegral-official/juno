package index

type SimpleInvertedIndex struct {
}

func (sii *SimpleInvertedIndex) HasNext() bool {
	return false
}
func (sii *SimpleInvertedIndex) Next() DocId {
	return 0
}
func (sii *SimpleInvertedIndex) GetGE(id DocId) DocId {
	return 0
}
