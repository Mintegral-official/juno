package index

import "github.com/Mintegral-official/juno/document"

type InvertedIndex interface {
	Add(fieldName string, id document.DocId) error
	Del(fieldName string, id document.DocId) bool
	Iterator(fieldName string) InvertedIterator
}
