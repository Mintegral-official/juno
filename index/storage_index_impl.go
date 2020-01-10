package index

import (
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

func NewStorageIndexer(isDebug ...int) *StorageIndexer {
	s := &StorageIndexer{
		data: sync.Map{},
	}
	if len(isDebug) != 0 && isDebug[0] == 1 {
		s.aDebug = debug.NewDebug("storage index")
	}
	return s
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
		sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
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
			if s.aDebug != nil {
				s.aDebug.AddDebugMsg("index: " + fieldName + " len: " + strconv.Itoa(sl.Len()))
			}
			return sl.Iterator()
		}
	}
	if s.aDebug != nil {
		s.aDebug.AddDebugMsg("index: " + fieldName + " is nil")
	}
	sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
	return sl.Iterator()
}

func (s *StorageIndexer) DebugInfo() *debug.Debug {
	return s.aDebug
}
