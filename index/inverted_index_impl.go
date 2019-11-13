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

func (sii *InvertedIndexImpl) Add(fieldName string, id document.DocId) error {
	if v, ok := sii.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipListIterator); ok {
			sl.Add(id, nil)
		} else {
			return helpers.ParseError
		}
	} else {
		sl := datastruct.NewSkipListIterator(datastruct.DEFAULT_MAX_LEVEL, helpers.DocIdFunc)
		sl.Add(id, nil)
		sii.data.Store(fieldName, sl)
	}
	return nil
}

func (sii *InvertedIndexImpl) Del(fieldName string, id document.DocId) bool {

	if v, ok := sii.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipListIterator); ok {
			sl.Del(id)
			sii.data.Store(fieldName, sl)
			return true
		}
	}
	return false
}

func (slii *InvertedIndexImpl) Iterator(fieldName string) datastruct.Iterator {
	if v, ok := slii.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipListIterator); ok {
			return sl.Iterator()
		}
	}
	return nil
}
