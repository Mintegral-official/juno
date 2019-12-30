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

func (s *StorageIndexer) Count() int {
	var count = 0
	s.data.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

func (s *StorageIndexer) Get(fieldName string, id document.DocId) interface{} {
	if v, ok := s.data.Load(fieldName); ok {
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

func (s *StorageIndexer) Add(fieldName string, id document.DocId, value interface{}) error {
	if v, ok := s.data.Load(fieldName); ok {
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
		s.data.Store(fieldName, sl)
	}
	return nil
}

func (s *StorageIndexer) Del(fieldName string, id document.DocId) bool {
	if v, ok := s.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipList); ok {
			sl.Del(id)
			s.data.Store(fieldName, sl)
			return true
		}
	}
	return false
}

func (s *StorageIndexer) Iterator(fieldName string) datastruct.Iterator {
	if v, ok := s.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipList); ok {
			s.aDebug.AddDebug("index: " + fieldName + " len: " + strconv.Itoa(sl.Len()))
			return sl.Iterator()
		}
	}
	s.aDebug.AddDebug("index: " + fieldName + " is nil")
	sl, _ := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
	return sl.Iterator()
}

func (s *StorageIndexer) String() string {
	if res, err := json.Marshal(s.aDebug); err == nil {
		return string(res)
	} else {
		return err.Error()
	}
}
