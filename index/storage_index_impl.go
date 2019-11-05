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

func (ssi *SimpleStorageIndex) Get(fieldName string, id document.DocId) interface{} {
	if v, ok := ssi.data.Load(fieldName); ok {
		if sl, ok := v.(*SkipList); ok {
			if res, err := sl.Get(id); err != nil {
				return res
			}
			return helpers.DOCUMENT_ERROR
		} else {
			return helpers.PARSE_ERROR
		}
	}
	return nil
}

func (ssi *SimpleStorageIndex) Add(fieldName string, id document.DocId, value interface{}) error {
	if v, ok := ssi.data.Load(fieldName); ok {
		if sl, ok := v.(*SkipList); ok {
			sl.Add(id, value)
		} else {
			return helpers.PARSE_ERROR
		}
	} else {
		sl := NewSkipList(DEFAULT_MAX_LEVEL, helpers.DocIdFunc)
		sl.Add(id, value)
		ssi.data.Store(fieldName, sl)
	}
	return nil
}

func (ssi *SimpleStorageIndex) Del(fieldName string, id document.DocId) bool {
	if v, ok := ssi.data.Load(fieldName); ok {
		if sl, ok := v.(*SkipList); ok {
			sl.Del(id)
			ssi.data.Store(fieldName, sl)
			return true
		}
	}
	return false
}

func (ssi *SimpleStorageIndex) Iterator(fieldName string) InvertedIterator {
	if v, ok := ssi.data.Load(fieldName); ok {
		if sl, ok := v.(*SkipListIterator); ok {
			return sl.Iterator()
		}
	}
	return nil
}
