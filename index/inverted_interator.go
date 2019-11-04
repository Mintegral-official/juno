package index

type InvertedIterator interface {
	HasNext() bool
	Next() *Element
	GetGE(id interface{}) interface{}
}
