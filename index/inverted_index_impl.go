package index

import (
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"sync"
)

type SimpleInvertedIndex struct {
	data sync.Map
}

func NewSimpleInvertedIndex() *SimpleInvertedIndex {
	return &SimpleInvertedIndex{data: sync.Map{}}
}

func (sii *SimpleInvertedIndex) Add(fieldName string, id document.DocId) error {
	if v, ok := sii.data.Load(fieldName); ok {
		if sl, ok := v.(*SkipList); ok {
			sl.Add(id, nil)
		} else {
			return helpers.ERROR_PARSE_ERROR
		}
	} else {
		sl := NewSkipList(DEFAULT_MAX_LEVEL, helpers.DocIdFunc)
		sl.Add(id, nil)
		sii.data.Store(fieldName, sl)
	}
	return nil
}

func (sii *SimpleInvertedIndex) Del(fieldName string, id document.DocId) bool {

	if v, ok := sii.data.Load(fieldName); ok {
		if sl, ok := v.(*SkipList); ok {
			sl.Del(id)
			sii.data.Store(fieldName, sl)
			return true
		}
	}
    return false
}

func (slii *SimpleInvertedIndex) Iterator(fieldName string) InvertedIterator {

	if v, ok := slii.data.Load(fieldName); ok {
		if sl, ok := v.(*SkipListIterator); ok {
			return sl.Iterator()
		}
	}
	return nil
}
