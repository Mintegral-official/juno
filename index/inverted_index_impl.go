package index

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"strconv"
	"sync"
)

type InvertedIndexer struct {
	data   sync.Map
	aDebug *debug.Debug
}

func NewInvertedIndexer(isDebug ...int) *InvertedIndexer {
	i := &InvertedIndexer{
		data: sync.Map{},
	}
	if len(isDebug) != 0 && isDebug[0] == 1 {
		i.aDebug = debug.NewDebug("invert index")
	}
	return i
}

func (i *InvertedIndexer) Count() int {
	var count = 0
	i.data.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

func (i *InvertedIndexer) Add(fieldName string, id document.DocId) error {
	v, ok := i.data.Load(fieldName)
	if !ok {
		sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
		sl.Add(id, nil)
		i.data.Store(fieldName, sl)
		return nil
	}
	sl, ok := v.(*datastruct.SkipList)
	if !ok {
		return helpers.ParseError
	}
	sl.Add(id, nil)
	return nil
}

func (i *InvertedIndexer) Del(fieldName string, id document.DocId) bool {
	v, ok := i.data.Load(fieldName)
	if !ok {
		return false
	}
	sl, ok := v.(*datastruct.SkipList)
	if ok {
		sl.Del(id)
		i.data.Store(fieldName, sl)
		return true
	}
	return false
}

func (i *InvertedIndexer) Iterator(name, value string) datastruct.Iterator {
	var fieldName = name + "_" + value
	v, ok := i.data.Load(fieldName)
	if ok {
		sl, ok := v.(*datastruct.SkipList)
		if ok {
			if i.aDebug != nil {
				i.aDebug.AddDebugMsg("index[" + fieldName + "] len: " + strconv.Itoa(sl.Len()))
			}
			return sl.Iterator()
		}
	}
	if i.aDebug != nil {
		i.aDebug.AddDebugMsg("index: " + fieldName + " is nil")
	}
	sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
	return sl.Iterator()
}

func (i *InvertedIndexer) DebugInfo() *debug.Debug {
	return i.aDebug
}
