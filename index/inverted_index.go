package index

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
)

type InvertedIndex interface {
	Add(fieldName string, id document.DocId) error
	Del(fieldName string, id document.DocId) bool
	Iterator(name string, value interface{}) datastruct.Iterator
	Count() int
	String() string
}
