package index

import (
	"github.com/Mintegral-official/juno/document"
	"sync"
)

type SimpleStorageIndex struct {
	data sync.Map
}

func NewSimpleStorageIndex() *SimpleStorageIndex {
	return &SimpleStorageIndex{}
}

func (ssi *SimpleStorageIndex) Get(filedName string, id document.DocId) interface{} {
	return nil
}
