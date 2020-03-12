package index

import (
	"github.com/MintegralTech/juno/datastruct"
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
)

type InvertedIndex interface {
	Add(fieldName string, id document.DocId) error
	Del(fieldName string, id document.DocId) bool
	Update(fieldName string, ids []document.DocId)
	Delete(fieldName string)
	Iterator(name, value string) datastruct.Iterator
	Range(func(key, value interface{}) bool)
	Count() int
	GetValueById(id document.DocId) map[string][]string
	SetDebug(level int)
	DebugInfo() *debug.Debug
}
