package index

import (
	"errors"
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
	v, ok := sii.data.Load(fieldName)
	if ok {
		sl, ok := v.(*SkipList)
		if ok {
			sl.Add(id, nil)
		} else {
			return errors.New("Parse  error")
		}
	} else {
		sl := NewSkipList(12, helpers.DocIdFunc)
		sl.Add(id, nil)
		sii.data.Store(fieldName, sl)
	}
	return nil
}

func (sii *SimpleInvertedIndex) Del(fieldName string, id document.DocId) bool {
	v, ok := sii.data.Load(fieldName)
	if !ok {
		return false
	}
	sl, ok := v.(*SkipList)
	if !ok {
		return false
	}
	ok = sl.Del(id)
	if ok {
		return true
	}
	return false
}

func (slii *SimpleInvertedIndex) Iterator(fieldName string) InvertedIterator {
	v, ok := slii.data.Load(fieldName)
	if ok {
		sl, ok := v.(*SkipListIterator)
		if ok {
			return sl.Iterator()
		}
	}
	return nil
}
