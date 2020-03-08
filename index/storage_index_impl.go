package index

import (
	"fmt"
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
		data: sync.Map{},
	}
}

func (s *StorageIndexer) GetValueById(id document.DocId) map[string][]string {
	var res = make(map[string][]string, 16)
	s.data.Range(func(key, value interface{}) bool {
		v, ok := value.(*datastruct.SkipList)
		if !ok {
			return true
		}
		e := v.Iterator().GetGE(id)
		if e == nil {
			return true
		}
		if e.Key() != id {
			return true
		}
		res[key.(string)] = append(res[key.(string)], fmt.Sprintf("%v", e.Value()))
		return true
	})
	return res
}

func (s *StorageIndexer) Count() (count int) {
	count = 0
	s.data.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

func (s *StorageIndexer) Get(fieldName string, id document.DocId) interface{} {
	v, ok := s.data.Load(fieldName)
	if !ok {
		return nil
	}
	sl, ok := v.(*datastruct.SkipList)
	if !ok {
		return helpers.ParseError
	}
	if res, err := sl.Get(id); err == nil {
		return res
	}
	return helpers.DocumentError
}

func (s *StorageIndexer) Add(fieldName string, id document.DocId, value interface{}) (err error) {
	v, ok := s.data.Load(fieldName)
	if !ok {
		sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
		sl.Add(id, value)
		s.data.Store(fieldName, sl)
		return err
	}
	sl, ok := v.(*datastruct.SkipList)
	if !ok {
		err = helpers.ParseError
		return
	}
	sl.Add(id, value)
	return err
}

func (s *StorageIndexer) Del(fieldName string, id document.DocId) (ok bool) {
	v, ok := s.data.Load(fieldName)
	if !ok {
		return ok
	}
	if sl, ok := v.(*datastruct.SkipList); ok {
		sl.Del(id)
		s.data.Store(fieldName, sl)
		return ok
	}
	return ok
}

func (s *StorageIndexer) Iterator(fieldName string) datastruct.Iterator {
	if v, ok := s.data.Load(fieldName); ok {
		sl, ok := v.(*datastruct.SkipList)
		if ok {
			if s.aDebug != nil {
				s.aDebug.AddDebugMsg("index: " + fieldName + " len: " + strconv.Itoa(sl.Len()))
			}
			iter := sl.Iterator()
			iter.FieldName = fieldName
			return iter
		}
	}
	if s.aDebug != nil {
		s.aDebug.AddDebugMsg("index: " + fieldName + " is nil")
	}
	sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel).Iterator()
	sl.FieldName = fieldName
	return sl
}

func (s *StorageIndexer) DebugInfo() *debug.Debug {
	return s.aDebug
}

func (s *StorageIndexer) SetDebug(level int) {
	if s.aDebug == nil {
		s.aDebug = debug.NewDebug(level, "storage index")
	}
}
