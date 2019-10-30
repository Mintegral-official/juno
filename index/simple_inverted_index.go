package index

import (
	"errors"
	"github.com/Mintegral-official/juno/document"
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
			sl.Add(id)
		} else {
			return errors.New("Parse  error")
		}
	} else {
		sl := NewSkipList(12)
		sl.Add(id)
		sii.data.Store(fieldName, sl)
	}
	return nil
}

func (sii *SimpleInvertedIndex) Del(fieldName string, id document.DocId) {
	v, ok := sii.data.Load(fieldName)
	if ok {
		sl, ok := v.(*SkipList)
		if ok {
			sl.Del(id)
		}
	}
}

func (sii *SimpleInvertedIndex) Iterator(fieldName string) InvertedIterator {
	v, ok := sii.data.Load(fieldName)
	if ok {
		sl, ok := v.(*SkipList)
		if ok {
			return sl.Iterator()
		}
	}
	return nil
}
