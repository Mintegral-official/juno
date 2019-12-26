package index

import (
	"encoding/json"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"strconv"
	"sync"
)

type StorageIndexer struct {
	data   sync.Map
	aDebug *debug.Debug
}

func NewStorageIndexer() *StorageIndexer {
	return &StorageIndexer{
		data:   sync.Map{},
		aDebug: debug.NewDebug("storage index"),
	}
}

func (sIndexer *StorageIndexer) Count() int {
	var count = 0
	sIndexer.data.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

func (sIndexer *StorageIndexer) Get(fieldName string, id document.DocId) interface{} {
	if v, ok := sIndexer.data.Load(fieldName); ok {
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

func (sIndexer *StorageIndexer) Add(fieldName string, id document.DocId, value interface{}) error {
	if v, ok := sIndexer.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipList); ok {
			sl.Add(id, value)
		} else {
			return helpers.ParseError
		}
	} else {
		sl, err := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
		if err != nil {
			return err
		}
		sl.Add(id, value)
		sIndexer.data.Store(fieldName, sl)
	}
	return nil
}

func (sIndexer *StorageIndexer) Del(fieldName string, id document.DocId) bool {
	if v, ok := sIndexer.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipList); ok {
			sl.Del(id)
			sIndexer.data.Store(fieldName, sl)
			return true
		}
	}
	return false
}

func (sIndexer *StorageIndexer) Iterator(fieldName string) datastruct.Iterator {
	if v, ok := sIndexer.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipList); ok {
			sIndexer.aDebug.AddDebug("index: " + fieldName + " len: " + strconv.Itoa(sl.Len()))
			return sl.Iterator()
		}
	}
	sIndexer.aDebug.AddDebug("index: " + fieldName + " is nil")
	sl, _ := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
	return sl.Iterator()
}

func (sIndexer *StorageIndexer) String() string {
	if res, err := json.Marshal(sIndexer.aDebug); err == nil {
		return string(res)
	} else {
		return err.Error()
	}
}
