package index

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"sync"
)

type InvertedIndexImpl struct {
	data sync.Map
}

func NewInvertedIndexImpl() *InvertedIndexImpl {
	return &InvertedIndexImpl{data: sync.Map{}}
}

func (iiImpl *InvertedIndexImpl) Add(fieldName string, id document.DocId) error {
	if v, ok := iiImpl.data.Load(fieldName); ok {
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
		iiImpl.data.Store(fieldName, sl)
	}
	return nil
}

func (iiImpl *InvertedIndexImpl) Del(fieldName string, id document.DocId) bool {

	if v, ok := iiImpl.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipList); ok {
			sl.Del(id)
			iiImpl.data.Store(fieldName, sl)
			return true
		}
	}
	return false
}

func (iiImpl *InvertedIndexImpl) Iterator(fieldName string) datastruct.Iterator {
	if v, ok := iiImpl.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipList); ok {
			return sl.Iterator()
		}
	}
	return nil
}
