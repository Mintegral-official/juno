package index

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"sync"
)

type InvertedIndexer struct {
	data sync.Map
}

func NewInvertedIndexer() *InvertedIndexer {
	return &InvertedIndexer{data: sync.Map{}}
}

func (iIndexer *InvertedIndexer) Count() int {
	var count = 0
	iIndexer.data.Range(func(key, value interface{}) bool {
		if key != nil {
			count++
			return true
		}
		return false
	})
	return count
}

func (iIndexer *InvertedIndexer) Add(fieldName string, id document.DocId) error {
	if v, ok := iIndexer.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipList); ok {
			sl.Add(id, nil)
		} else {
			return helpers.ParseError
		}
	} else {
		sl, err := datastruct.NewSkipList(datastruct.DefaultMaxLevel, helpers.DocIdFunc)
		if err != nil {
			return err
		}
		sl.Add(id, nil)
		iIndexer.data.Store(fieldName, sl)
	}
	return nil
}

func (iIndexer *InvertedIndexer) Del(fieldName string, id document.DocId) bool {

	if v, ok := iIndexer.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipList); ok {
			sl.Del(id)
			iIndexer.data.Store(fieldName, sl)
			return true
		}
	}
	return false
}

func (iIndexer *InvertedIndexer) Iterator(fieldName string) datastruct.Iterator {
	if v, ok := iIndexer.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipList); ok {
			return sl.Iterator()
		}
	}
	return nil
}
