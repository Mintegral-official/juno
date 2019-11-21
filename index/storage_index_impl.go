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

func (siImpl *StorageIndexImpl) Get(fieldName string, id document.DocId) interface{} {
	if v, ok := siImpl.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipList); ok {
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

func (siImpl *StorageIndexImpl) Add(fieldName string, id document.DocId, value interface{}) error {
	if v, ok := siImpl.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipList); ok {
			sl.Add(id, value)
		} else {
			return helpers.ParseError
		}
	} else {
		sl, err := datastruct.NewSkipList(datastruct.DefaultMaxLevel, helpers.DocIdFunc)
		if err != nil {
			return err
		}
		sl.Add(id, value)
		siImpl.data.Store(fieldName, sl)
	}
	return nil
}

func (siImpl *StorageIndexImpl) Del(fieldName string, id document.DocId) bool {
	if v, ok := siImpl.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipList); ok {
			sl.Del(id)
			siImpl.data.Store(fieldName, sl)
			return true
		}
	}
	return false
}

func (siImpl *StorageIndexImpl) Iterator(fieldName string) datastruct.Iterator {
	if v, ok := siImpl.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipList); ok {
			return sl.Iterator()
		}
	}
	return nil
}
