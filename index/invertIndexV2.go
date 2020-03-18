package index

import (
	"github.com/MintegralTech/juno/datastruct"
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/helpers"
	"strconv"
	"strings"
	"sync"
)

type InvertedIndexerV2 struct {
	data   sync.Map
	aDebug *debug.Debug
}

func NewInvertedIndexV2() *InvertedIndexerV2 {
	return &InvertedIndexerV2{
		data: sync.Map{},
	}
}

func (i *InvertedIndexerV2) Count() (count int) {
	count = 0
	i.data.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

func (i *InvertedIndexerV2) Range(f func(key, value interface{}) bool) {
	i.data.Range(f)
}

func (i *InvertedIndexerV2) GetInvertIndexDebugInfoById(id document.DocId) map[string][]string {
	var res = make(map[string][]string, 16)
	i.data.Range(func(key, value interface{}) bool {
		v, ok := value.(*datastruct.Slice)
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
		keys := strings.Split(key.(string), SEP)
		if len(keys) == 2 {
			res[keys[0]] = append(res[keys[0]], keys[1])
		}
		return true
	})
	return res
}

func (i *InvertedIndexerV2) Add(fieldName string, id document.DocId) (err error) {
	v, ok := i.data.Load(fieldName)
	if !ok {
		sl := datastruct.NewSlice()
		sl.Add(id, nil)
		i.data.Store(fieldName, sl)
		return err
	}
	sl, ok := v.(*datastruct.Slice)
	if !ok {
		err = helpers.ParseError
		return err
	}
	sl.Add(id, nil)
	return err
}

func (i *InvertedIndexerV2) Del(fieldName string, id document.DocId) (ok bool) {
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

func (i *InvertedIndexerV2) Update(fieldName string, ids []document.DocId) {
	v, ok := i.data.Load(fieldName)
	if !ok {
		sl := datastruct.NewSlice()
		for _, id := range ids {
			sl.Add(id, nil)
		}
		i.data.Store(fieldName, sl)
		return
	}
	if sl, ok := v.(*datastruct.Slice); ok {
		sl = datastruct.NewSlice()
		for _, id := range ids {
			sl.Add(id, nil)
		}
		i.data.Store(fieldName, sl)
	}
}

func (i *InvertedIndexerV2) Delete(fieldName string) {
	i.data.Delete(fieldName)
}

func (i *InvertedIndexerV2) Iterator(name, value string) datastruct.Iterator {
	var fieldName = name + SEP + value
	if v, ok := i.data.Load(fieldName); ok {
		sl, ok := v.(*datastruct.Slice)
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

func (i *InvertedIndexerV2) DebugInfo() *debug.Debug {
	return i.aDebug
}

func (i *InvertedIndexerV2) SetDebug(level int) {
	if i.aDebug == nil {
		i.aDebug = debug.NewDebug(level, "invert index")
	}
}
