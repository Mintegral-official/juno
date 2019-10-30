package index

import "github.com/Mintegral-official/juno/document"

type StorageIndex interface {
	Get(filedName string, id document.DocId) interface{}
}
