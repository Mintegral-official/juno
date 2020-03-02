package index

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"strconv"
	"strings"
	"sync"
)

type InvertedIndexer struct {
	data   sync.Map
	aDebug *debug.Debug
}

func NewInvertedIndexer(isDebug ...int) (i *InvertedIndexer) {
	i = &InvertedIndexer{
		data: sync.Map{},
	}
	if len(isDebug) != 0 && isDebug[0] == 1 {
		i.aDebug = debug.NewDebug("invert index")
	}
	return i
}

func (i *InvertedIndexer) Count() (count int) {
	count = 0
	i.data.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

func (i *InvertedIndexer) GetValueById(id document.DocId) map[string][]string {
	var res = make(map[string][]string, 16)
	i.data.Range(func(key, value interface{}) bool {
		v, ok := value.(*datastruct.SkipList)
		if !ok {
			return true
		}
		e := v.Iterator().GetGE(id)
		if e == nil {
			return true
		}
		keys := strings.Split(key.(string), SEP)
		res[keys[0]] = append(res[keys[0]], keys[1])
		return true
	})
	return res
}

func (i *InvertedIndexer) Add(fieldName string, id document.DocId) (err error) {
	v, ok := i.data.Load(fieldName)
	if !ok {
		sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
		sl.Add(id, nil)
		i.data.Store(fieldName, sl)
		return err
	}
	sl, ok := v.(*datastruct.SkipList)
	if !ok {
		err = helpers.ParseError
		return err
	}
	sl.Add(id, nil)
	return err
}

func (i *InvertedIndexer) Del(fieldName string, id document.DocId) (ok bool) {
	v, ok := i.data.Load(fieldName)
	if !ok {
		return ok
	}
	if sl, ok := v.(*datastruct.SkipList); ok {
		sl.Del(id)
		i.data.Store(fieldName, sl)
		return ok
	}
	return ok
}

func (i *InvertedIndexer) Update(fieldName string, ids []document.DocId) {
	v, ok := i.data.Load(fieldName)
	if !ok {
		sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
		for _, id := range ids {
			sl.Add(id, nil)
		}
		i.data.Store(fieldName, sl)
		return
	}
	if sl, ok := v.(*datastruct.SkipList); ok {
		sl = datastruct.NewSkipList(datastruct.DefaultMaxLevel)
		for _, id := range ids {
			sl.Add(id, nil)
		}
		i.data.Store(fieldName, sl)
	}
}

func (i *InvertedIndexer) Delete(fieldName string) {
	i.data.Delete(fieldName)
}

func (i *InvertedIndexer) Iterator(name, value string) datastruct.Iterator {
	var fieldName = name + SEP + value
	if v, ok := i.data.Load(fieldName); ok {
		sl, ok := v.(*datastruct.SkipList)
		if ok {
			if i.aDebug != nil {
				i.aDebug.AddDebugMsg("index[" + fieldName + "] len: " + strconv.Itoa(sl.Len()))
			}
			iter := sl.Iterator()
			iter.FieldName = fieldName
			return iter
		}
	}
	if i.aDebug != nil {
		i.aDebug.AddDebugMsg("index: " + fieldName + " is nil")
	}
	sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel).Iterator()
	sl.FieldName = fieldName
	return sl
}

func (i *InvertedIndexer) DebugInfo() *debug.Debug {
	return i.aDebug
}
