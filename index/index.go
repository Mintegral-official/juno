package index

import (
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
)

type IndexInfo struct {
	DocSize           int `json:"doc_size"`
	InvertedIndexSize int `json:"inverted_index_szie"`
	StorageIndex      int `json:"storage_index_size"`
}

type Index interface {
	Add(docInfo *document.DocInfo) error
	UpdateIds(fieldName string, ids []document.DocId)
	Delete(fieldName string)
	Del(docInfo *document.DocInfo)

	GetName() string
	GetIndexInfo() *IndexInfo
	GetInvertedIndex() InvertedIndex
	GetStorageIndex() StorageIndex
	GetDataType(fieldName string) document.FieldType
	GetValueById(id document.DocId) [2]map[string][]string
	GetId(id document.DocId) (document.DocId, error)
	GetInnerId(id document.DocId) (document.DocId, error)

	Dump(filename string) error
	Load(filename string) error

	DebugInfo() *debug.Debug
}

func NewIndex(name string) Index {
	return NewIndexV2(name)
}
