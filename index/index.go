package index

import (
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
)

type Index interface {
	Add(docInfo *document.DocInfo) error
	UpdateIds(fieldName string, ids []document.DocId)
	Delete(fieldName string)
	Del(docInfo *document.DocInfo)
	GetDataType(fieldName string) document.FieldType
	GetValueById(id document.DocId) [2]map[string][]string
	Dump(filename string) error
	Load(filename string) error
	DebugInfo() *debug.Debug
}
