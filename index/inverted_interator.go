package index

type InvertedIterator interface {
	HasNext() bool
	Next() interface{}
	GetGE(id interface{}) interface{}
}
