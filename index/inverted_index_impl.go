package index

import (
	"encoding/json"
	"fmt"
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

func NewInvertedIndexer() *InvertedIndexer {
	return &InvertedIndexer{
		data:   sync.Map{},
		aDebug: debug.NewDebug("invert index"),
	}
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
	if v, ok := i.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipList); ok {
			sl.Add(id, nil)
		} else {
			return helpers.ParseError
		}
	} else {
		sl, err := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
		if err != nil {
			return err
		}
		sl.Add(id, nil)
		i.data.Store(fieldName, sl)
	}
	return nil
}

func (i *InvertedIndexer) Del(fieldName string, id document.DocId) bool {

	if v, ok := i.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipList); ok {
			sl.Del(id)
			i.data.Store(fieldName, sl)
			return true
		}
	}
	return false
}

func (i *InvertedIndexer) Iterator(name string, value interface{}) datastruct.Iterator {
	var fieldName = name + "_" + fmt.Sprint(value)
	if v, ok := i.data.Load(fieldName); ok {
		if sl, ok := v.(*datastruct.SkipList); ok {
			i.aDebug.AddDebug("index[" + fieldName + "] len: " + strconv.Itoa(sl.Len()))
			return sl.Iterator()
		}
	}
	i.aDebug.AddDebug("index: " + fieldName + " is nil")
	sl, _ := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
	return sl.Iterator()
}

func (i *InvertedIndexer) String() string {
	if res, err := json.Marshal(i.aDebug); err == nil {
		return string(res)
	} else {
		return err.Error()
	}
}
