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

func NewIndex(name string, isDebug ...int) *Indexer {
	i := &Indexer{
		invertedIndex:   NewInvertedIndexer(),
		storageIndex:    NewStorageIndexer(),
		campaignMapping: concurrent_map.CreateConcurrentMap(199),
		kvType:          concurrent_map.CreateConcurrentMap(199),
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
		return helpers.DocumentError
	}
	for _, field := range doc.Fields {
		switch field.IndexType {
		case document.InvertedIndexType:
			err = i.invertAdd(doc.Id, field)
			if err != nil {
				if i.aDebug != nil {
					i.aDebug.AddDebugMsg(i.StringBuilder(256, "invert", doc.Id, field.Name, field.Value, err.Error()))
				}
				return err
			}
		case document.StorageIndexType:
			err = i.storageAdd(doc.Id, field)
			if err != nil {
				if i.aDebug != nil {
					i.aDebug.AddDebugMsg(i.StringBuilder(256, "storage", doc.Id, field.Name, field.Value, err.Error()))
				}
				return err
			}
			i.kvType.Set(concurrent_map.StrKey(field.Name), field.ValueType)
		case document.BothIndexType:
			err = i.invertAdd(doc.Id, field)
			if err != nil {
				if i.aDebug != nil {
					i.aDebug.AddDebugMsg(i.StringBuilder(256, "invert", doc.Id, field.Name, field.Value, err.Error()))
				}
				return err
			}
			err = i.storageAdd(doc.Id, field)
			if err != nil {
				if i.aDebug != nil {
					i.aDebug.AddDebugMsg(i.StringBuilder(256, "storage", doc.Id, field.Name, field.Value, err.Error()))
				}
				return err
			}
			i.kvType.Set(concurrent_map.StrKey(field.Name), field.ValueType)
		default:
			i.WarnStatus(field.Name, field.Value, "type is wrong")
			return errors.New("the add doc type is wrong or nil ")
		}
	}
	return nil
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
	if t, ok := i.kvType.Get(concurrent_map.StrKey(fieldName)); ok {
		return t.(document.FieldType)
	}
	return document.DefaultFieldType
}

func (i *Indexer) invertAdd(id document.DocId, field *document.Field) (err error) {
	if v, ok := field.Value.([]string); ok {
		for _, s := range v {
			if err = i.invertedIndex.Add(field.Name+"_"+s, id); err != nil {
				i.WarnStatus(field.Name, s, err.Error())
				return err
			}
			i.campaignMapping.Set(DocId(id), i.count)
			i.bitmap.Set(i.count)
			i.count++
		}
	} else if v, ok := field.Value.([]int64); ok {
		for _, s := range v {
			if err = i.invertedIndex.Add(field.Name+"_"+strconv.FormatInt(s, 10), id); err != nil {
				i.WarnStatus(field.Name, s, err.Error())
				return err
			}
			i.campaignMapping.Set(DocId(id), i.count)
			i.bitmap.Set(i.count)
			i.count++
		}
	} else if v, ok := field.Value.(string); ok {
		if err = i.invertedIndex.Add(field.Name+"_"+v, id); err != nil {
			i.WarnStatus(field.Name, v, err.Error())
			return err
		}
		i.campaignMapping.Set(DocId(id), i.count)
		i.bitmap.Set(i.count)
		i.count++
	} else if v, ok := field.Value.(int64); ok {
		if err = i.invertedIndex.Add(field.Name+"_"+strconv.FormatInt(v, 10), id); err != nil {
			i.WarnStatus(field.Name, v, err.Error())
			return err
		}
		i.campaignMapping.Set(DocId(id), i.count)
		i.bitmap.Set(i.count)
		i.count++
	} else {
		return errors.New("the doc is nil or type is wrong")
	}
	return nil
}

func (i *Indexer) storageAdd(id document.DocId, field *document.Field) (err error) {
	if err = i.storageIndex.Add(field.Name, id, field.Value); err != nil {
		i.WarnStatus(field.Name, field.Value, err.Error())
		return err
	}
	return nil
}

func (i *Indexer) invertDel(id document.DocId, field *document.Field) {
	if v, ok := field.Value.([]string); ok {
		for _, s := range v {
			i.invertedIndex.Del(field.Name+"_"+s, id)
			if v, ok := i.GetCampaignMap().Get(DocId(id)); ok {
				i.bitmap.Del(v.(document.DocId))
			}
		}
	} else if v, ok := field.Value.([]int64); ok {
		for _, s := range v {
			i.invertedIndex.Del(field.Name+"_"+fmt.Sprint(s), id)
			if v, ok := i.GetCampaignMap().Get(DocId(id)); ok {
				i.bitmap.Del(v.(document.DocId))
			}
		}
	} else if v, ok := field.Value.(string); ok {
		i.invertedIndex.Del(field.Name+"_"+v, id)
		if v, ok := i.GetCampaignMap().Get(DocId(id)); ok {
			i.bitmap.Del(v.(document.DocId))
		}
	} else if v, ok := field.Value.(int64); ok {
		i.invertedIndex.Del(field.Name+"_"+fmt.Sprint(v), id)
		if v, ok := i.GetCampaignMap().Get(DocId(id)); ok {
			i.bitmap.Del(v.(document.DocId))
		}
	}
}

func (i *Indexer) storageDel(id document.DocId, field *document.Field) {
	i.storageIndex.Del(field.Name, id)
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
