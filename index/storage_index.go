package index

import "github.com/Mintegral-official/juno/document"

type StorageIndex interface {
	Get(filedName string, id document.DocId) interface{}
	Add(fieldName string, id document.DocId, value interface{}) error
	Del(fieldName string, id document.DocId) bool
	Iterator(fieldName string) InvertedIterator
}
