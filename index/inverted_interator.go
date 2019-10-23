package index

type InvertedIterator interface {
	HasNext() bool
	Next() DocId
	GetGE(id DocId) DocId
}
