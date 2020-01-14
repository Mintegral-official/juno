package index

import (
	"errors"
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/log"
	"github.com/easierway/concurrent_map"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type Indexer struct {
	invertedIndex   InvertedIndex
	storageIndex    StorageIndex
	campaignMapping *concurrent_map.ConcurrentMap
	bitmap          *datastruct.BitSet
	count           document.DocId
	name            string
	kvType          *concurrent_map.ConcurrentMap
	logger          log.Logger
	aDebug          *debug.Debug
}

func NewIndex(name string, isDebug ...int) (i *Indexer) {
	i = &Indexer{
		invertedIndex:   NewInvertedIndexer(),
		storageIndex:    NewStorageIndexer(),
		campaignMapping: concurrent_map.CreateConcurrentMap(128),
		kvType:          concurrent_map.CreateConcurrentMap(128),
		bitmap:          datastruct.NewBitMap(),
		count:           1,
		name:            name,
		logger:          logrus.New(),
	}
	if len(isDebug) != 0 && isDebug[0] == 1 {
		i.aDebug = debug.NewDebug(name)
	}
	return i
}

func (i *Indexer) GetInvertedIndex() InvertedIndex {
	return i.invertedIndex
}

func (i *Indexer) GetStorageIndex() StorageIndex {
	return i.storageIndex
}

func (i *Indexer) GetCampaignMap() *concurrent_map.ConcurrentMap {
	return i.campaignMapping
}

func (i *Indexer) GetBitMap() *datastruct.BitSet {
	return i.bitmap
}

func (i *Indexer) GetName() string {
	return i.name
}

func (i *Indexer) UnsetDebug() {
	i.aDebug = debug.NewDebug(i.GetName())
}

func (i *Indexer) Add(doc *document.DocInfo) (err error) {
	if doc == nil {
		err = helpers.DocumentError
		return
	}
	for _, field := range doc.Fields {
		switch field.IndexType {
		case document.InvertedIndexType:
			err = i.invertAdd(doc.Id, field)
			if err != nil {
				if i.aDebug != nil {
					i.aDebug.AddDebugMsg(i.StringBuilder(256, "invert", doc.Id, field.Name, field.Value, err.Error()))
				}
				return
			}
		case document.StorageIndexType:
			err = i.storageAdd(doc.Id, field)
			if err != nil {
				if i.aDebug != nil {
					i.aDebug.AddDebugMsg(i.StringBuilder(256, "storage", doc.Id, field.Name, field.Value, err.Error()))
				}
				return
			}
			i.kvType.Set(concurrent_map.StrKey(field.Name), field.ValueType)
		case document.BothIndexType:
			err = i.invertAdd(doc.Id, field)
			if err != nil {
				if i.aDebug != nil {
					i.aDebug.AddDebugMsg(i.StringBuilder(256, "invert", doc.Id, field.Name, field.Value, err.Error()))
				}
				return
			}
			err = i.storageAdd(doc.Id, field)
			if err != nil {
				if i.aDebug != nil {
					i.aDebug.AddDebugMsg(i.StringBuilder(256, "storage", doc.Id, field.Name, field.Value, err.Error()))
				}
				return
			}
			i.kvType.Set(concurrent_map.StrKey(field.Name), field.ValueType)
		default:
			i.WarnStatus(field.Name, field.Value, "type is wrong")
			err = errors.New("the add doc type is wrong or nil ")
			return
		}
	}
	return
}

func (i *Indexer) Del(doc *document.DocInfo) {
	if doc == nil {
		return
	}
	for _, field := range doc.Fields {
		switch field.IndexType {
		case document.InvertedIndexType:
			i.invertDel(doc.Id, field)
		case document.StorageIndexType:
			i.storageDel(doc.Id, field)
		case document.BothIndexType:
			i.invertDel(doc.Id, field)
			i.storageDel(doc.Id, field)
		default:
			panic("the del doc type is nil or wrong")
		}
	}
}

func (i *Indexer) Update(filename string) error {
	if err := i.Dump(filename); err != nil {
		return err
	}
	return nil
}

func (i *Indexer) Dump(filename string) error {
	// TODO
	return nil
}

func (i *Indexer) Load(filename string) error {
	return nil
}

func (i *Indexer) GetDataType(fieldName string) document.FieldType {
	t, ok := i.kvType.Get(concurrent_map.StrKey(fieldName))
	if ok {
		return t.(document.FieldType)
	}
	return document.DefaultFieldType
}

func (i *Indexer) invertAdd(id document.DocId, field *document.Field) (err error) {
	switch field.Value.(type) {
	case []string:
		value, _ := field.Value.([]string)
		for _, v := range value {
			err = i.invertedIndex.Add(field.Name+"_"+v, id)
			if err != nil {
				i.WarnStatus(field.Name, v, err.Error())
				return err
			}
			i.campaignMapping.Set(DocId(id), i.count)
			i.bitmap.Set(i.count)
			i.count++
		}
	case []int64:
		value, _ := field.Value.([]int64)
		for _, v := range value {
			err = i.invertedIndex.Add(field.Name+"_"+strconv.FormatInt(v, 10), id)
			if err != nil {
				i.WarnStatus(field.Name, v, err.Error())
				return err
			}
			i.campaignMapping.Set(DocId(id), i.count)
			i.bitmap.Set(i.count)
			i.count++
		}
	case string:
		value, _ := field.Value.(string)
		err = i.invertedIndex.Add(field.Name+"_"+value, id)
		if err != nil {
			i.WarnStatus(field.Name, value, err.Error())
			return
		}
		i.campaignMapping.Set(DocId(id), i.count)
		i.bitmap.Set(i.count)
		i.count++
	case int64:
		value, _ := field.Value.(int64)
		err = i.invertedIndex.Add(field.Name+"_"+strconv.FormatInt(value, 10), id)
		if err != nil {
			i.WarnStatus(field.Name, value, err.Error())
			return
		}
		i.campaignMapping.Set(DocId(id), i.count)
		i.bitmap.Set(i.count)
		i.count++
	default:
		return errors.New("the doc is nil or type is wrong")
	}
	return
}

func (i *Indexer) storageAdd(id document.DocId, field *document.Field) (err error) {
	err = i.storageIndex.Add(field.Name, id, field.Value)
	if err != nil {
		i.WarnStatus(field.Name, field.Value, err.Error())
		return
	}
	return
}

func (i *Indexer) invertDel(id document.DocId, field *document.Field) {
	switch field.Value.(type) {
	case []string:
		value, _ := field.Value.([]string)
		for _, v := range value {
			i.invertedIndex.Del(field.Name+"_"+v, id)
			docId, ok := i.campaignMapping.Get(DocId(id))
			if ok {
				i.bitmap.Del(docId.(document.DocId))
			}
		}
	case []int64:
		value, _ := field.Value.([]int64)
		for _, v := range value {
			i.invertedIndex.Del(field.Name+"_"+strconv.FormatInt(v, 10), id)
			docId, ok := i.campaignMapping.Get(DocId(id))
			if ok {
				i.bitmap.Del(docId.(document.DocId))
			}
		}
	case string:
		value, _ := field.Value.(string)
		i.invertedIndex.Del(field.Name+"_"+value, id)
		docId, ok := i.campaignMapping.Get(DocId(id))
		if ok {
			i.bitmap.Del(docId.(document.DocId))
		}
	case int64:
		value, _ := field.Value.(int64)
		i.invertedIndex.Del(field.Name+"_"+strconv.FormatInt(value, 10), id)
		docId, ok := i.campaignMapping.Get(DocId(id))
		if ok {
			i.bitmap.Del(docId.(document.DocId))
		}
	default:
		if i.aDebug != nil {
			i.aDebug.AddDebugMsg(fmt.Sprintf("the del doc [%v] is nil or type is wrong", field))
		}
	}
}

func (i *Indexer) storageDel(id document.DocId, field *document.Field) {
	ok := i.storageIndex.Del(field.Name, id)
	if !ok {
		if i.aDebug != nil {
			i.aDebug.AddDebugMsg(fmt.Sprintf("del [%v] - [%v] failed", id, field))
		}
	}
}

func (i *Indexer) DebugInfo() *debug.Debug {
	if i.aDebug != nil {
		i.aDebug.AddDebugMsg("invert index count: " + strconv.Itoa(i.invertedIndex.Count()))
		i.aDebug.AddDebugMsg("storage index count: " + strconv.Itoa(i.storageIndex.Count()))
		i.aDebug.AddDebug(i.invertedIndex.DebugInfo(), i.storageIndex.DebugInfo())
		return i.aDebug
	}
	return nil
}

func (i *Indexer) StringBuilder(cap int, value ...interface{}) string {
	var b strings.Builder
	b.Grow(cap)
	_, _ = fmt.Fprintf(&b, "%s: ", value[0])
	b.WriteString("index add failed. ")
	_, _ = fmt.Fprintf(&b, "%s:[%v] ", "docId", value[1])
	_, _ = fmt.Fprintf(&b, "%s:[%v] ", "name", value[2])
	_, _ = fmt.Fprintf(&b, "%s:[%s] ", "value", value[3])
	_, _ = fmt.Fprintf(&b, "%s:[%s]", "reason", value[4])
	return b.String()
}

func (i *Indexer) WarnStatus(name string, value interface{}, err string) {
	if i.logger != nil {
		i.logger.Warnf("name:[%s] value:[%v] wrong reason:[%s]", name, value, err)
	}
}
