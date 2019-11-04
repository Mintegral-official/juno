package index

import (
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"sync"
)

type SimpleStorageIndex struct {
	data sync.Map
}

func NewSimpleStorageIndex() *SimpleStorageIndex {
	return &SimpleStorageIndex{data: sync.Map{}}
}

func (ssi *SimpleStorageIndex) Get(filedName string, id document.DocId) interface{} {
	v, ok := ssi.data.Load(filedName)
	if ok {
		if sl, ok := v.(*SkipList); ok {
			if res, err := sl.Get(id); err != nil {
				return res
			}
			return helpers.ERROR_DOCUMENT_ERROR
		} else {
			return helpers.ERROR_PARSE_ERROR
		}
	}
	return nil
}
