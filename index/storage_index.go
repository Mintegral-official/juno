package index

import (
	"github.com/MintegralTech/juno/datastruct"
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
)

type StorageIndex interface {
	Get(filedName string, id document.DocId) interface{}
	Add(fieldName string, id document.DocId, value interface{}) error
	Del(fieldName string, id document.DocId) bool
	Iterator(fieldName string) datastruct.Iterator
	Count() int
	Range(func(key, value interface{}) bool)
	GetValueById(id document.DocId) map[string][]string
	SetDebug(level int)
	DebugInfo() *debug.Debug
}
