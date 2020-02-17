package index

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
)

type InvertedIndex interface {
	Add(fieldName string, id document.DocId) error
	Del(fieldName string, id document.DocId) bool
	Update(fieldName string, ids []document.DocId)
	Iterator(name, value string) datastruct.Iterator
	Count() int
	DebugInfo() *debug.Debug
}
