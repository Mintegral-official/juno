package index

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"sync"
)

type StorageIndexImpl struct {
	data sync.Map
}

func NewStorageIndexImpl() *StorageIndexImpl {
	return &StorageIndexImpl{data: sync.Map{}}
}

func (ssi *StorageIndexImpl) Get(fieldName string, id document.DocId) interface{} {
	if v, ok := ssi.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipListIterator); ok {
			if res, err := sl.Get(id); err == nil {
				return res
			}
			return helpers.DocumentError
		} else {
			return helpers.ParseError
		}
	}
	return nil
}

func (ssi *StorageIndexImpl) Add(fieldName string, id document.DocId, value interface{}) error {
	if v, ok := ssi.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipListIterator); ok {
			sl.Add(id, value)
		} else {
			return helpers.ParseError
		}
	} else {
		sl := datastruct.NewSkipListIterator(datastruct.DEFAULT_MAX_LEVEL, helpers.DocIdFunc)
		sl.Add(id, value)
		ssi.data.Store(fieldName, sl)
	}
	return nil
}

func (ssi *StorageIndexImpl) Del(fieldName string, id document.DocId) bool {
	if v, ok := ssi.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipListIterator); ok {
			sl.Del(id)
			ssi.data.Store(fieldName, sl)
			return true
		}
	}
	return false
}

func (ssi *StorageIndexImpl) Iterator(fieldName string) datastruct.Iterator {
	if v, ok := ssi.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipListIterator); ok {
			return sl.Iterator()
		}
	}
	return nil
}
