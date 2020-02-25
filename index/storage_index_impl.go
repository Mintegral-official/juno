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
	data      sync.Map
	fieldName []string
	aDebug    *debug.Debug
}

func NewStorageIndexer(isDebug ...int) (s *StorageIndexer) {
	s = &StorageIndexer{
		data: sync.Map{},
	}
	if len(isDebug) != 0 && isDebug[0] == 1 {
		s.aDebug = debug.NewDebug("storage index")
	}
	return s
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
	s.fieldName = append(s.fieldName, fieldName)
	if v, ok := s.data.Load(fieldName); ok {
		sl, ok := v.(*datastruct.SkipList)
		if ok {
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

func (s *StorageIndexer) GetFieldName() []string {
	return s.fieldName
}
