package index

import (
	"github.com/Mintegral-official/juno/document"
)

type Index interface {
	Add(docInfo *document.DocInfo) error
	Del(docInfo *document.DocInfo)
	GetDataType(fieldName string) document.FieldType
	Dump(filename string) error
	Load(filename string) error
}
