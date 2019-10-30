package index

import (
	"errors"
	"github.com/Mintegral-official/juno/document"
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
			return sl.GetK(id)
		} else {
			return errors.New("Parse  error")
		}
	}
	return nil
}
