package index

import (
	"errors"
	"fmt"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/log"
	"github.com/easierway/concurrent_map"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"sync/atomic"
)

const SEP = "\007"

type Indexer struct {
	invertedIndex   InvertedIndex
	storageIndex    StorageIndex
	campaignMapping *concurrent_map.ConcurrentMap
	bitmap          *concurrent_map.ConcurrentMap
	count           uint64
	name            string
	kvType          *concurrent_map.ConcurrentMap
	logger          log.Logger
	aDebug          *debug.Debug
}

func NewIndex(name string) (i *Indexer) {
	i = &Indexer{
		invertedIndex:   NewInvertedIndexer(),
		storageIndex:    NewStorageIndexer(),
		campaignMapping: concurrent_map.CreateConcurrentMap(128),
		kvType:          concurrent_map.CreateConcurrentMap(128),
		bitmap:          concurrent_map.CreateConcurrentMap(128),
		count:           1,
		name:            name,
		logger:          logrus.New(),
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

func (i *Indexer) GetBitMap() *concurrent_map.ConcurrentMap {
	return i.bitmap
}

func (i *Indexer) GetName() string {
	return i.name
}

func (i *Indexer) SetDebug(level int) {
	if i.aDebug == nil {
		i.aDebug = debug.NewDebug(level, i.GetName())
		i.invertedIndex.SetDebug(level)
		i.storageIndex.SetDebug(level)
	}

}

func (i *Indexer) GetValueById(id document.DocId) [2]map[string][]string {
	var res [2]map[string][]string
	docId, ok := i.campaignMapping.Get(DocId(id))
	if ok {
		if _, ok := i.bitmap.Get(DocId(docId.(document.DocId))); !ok {
			return res
		}
		res[0] = i.GetInvertedIndex().GetValueById(docId.(document.DocId))
		res[1] = i.GetStorageIndex().GetValueById(docId.(document.DocId))
	}
	return res
}

func (i *Indexer) UpdateIds(fieldName string, ids []document.DocId) {
	var idList []document.DocId
	for _, id := range ids {
		if v, ok := i.campaignMapping.Get(DocId(id)); ok {
			if _, ok := i.bitmap.Get(DocId(v.(document.DocId))); ok {
				idList = append(idList, v.(document.DocId))
			} else {
				i.campaignMapping.Set(DocId(id), document.DocId(i.count))
				i.bitmap.Set(DocId(document.DocId(i.count)), id)
				idList = append(idList, document.DocId(i.count))
				atomic.AddUint64(&i.count, 1)
			}
		} else {
			i.campaignMapping.Set(DocId(id), document.DocId(i.count))
			i.bitmap.Set(DocId(document.DocId(i.count)), id)
			idList = append(idList, document.DocId(i.count))
			atomic.AddUint64(&i.count, 1)
		}
	}
	i.invertedIndex.Update(fieldName, idList)
}

func (i *Indexer) Delete(fieldName string) {
	i.invertedIndex.Delete(fieldName)
}

func (i *Indexer) Add(doc *document.DocInfo) (err error) {
	if doc == nil {
		return helpers.DocumentError
	}
	i.campaignMapping.Set(DocId(doc.Id), document.DocId(i.count))
	i.bitmap.Set(DocId(document.DocId(i.count)), doc.Id)
	for _, field := range doc.Fields {
		switch field.IndexType {
		case document.InvertedIndexType:
			i.invertAdd(document.DocId(i.count), field)
		case document.StorageIndexType:
			i.storageAdd(document.DocId(i.count), field)
			i.kvType.Set(concurrent_map.StrKey(field.Name), field.ValueType)
		case document.BothIndexType:
			i.invertAdd(document.DocId(i.count), field)
			i.storageAdd(document.DocId(i.count), field)
			i.kvType.Set(concurrent_map.StrKey(field.Name), field.ValueType)
		default:
			i.WarnStatus(field.Name, field.Value, "type is wrong")
			return errors.New("the add doc type is wrong or nil ")
		}
	}
	atomic.AddUint64(&i.count, 1)
	return err
}

func (i *Indexer) Del(doc *document.DocInfo) {
	if doc == nil {
		return
	}
	i.invertDel(doc.Id)
	for _, field := range doc.Fields {
		if field.IndexType == document.StorageIndexType || field.IndexType == document.BothIndexType {
			i.storageDel(doc.Id, field)
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

func (i *Indexer) invertAdd(id document.DocId, field *document.Field) {
	switch field.Value.(type) {
	case []string:
		value, _ := field.Value.([]string)
		for _, v := range value {
			if err := i.invertedIndex.Add(field.Name+SEP+v, id); err != nil {
				i.WarnStatus(field.Name, v, err.Error())
			}
		}
	case []int64:
		value, _ := field.Value.([]int64)
		for _, v := range value {
			if err := i.invertedIndex.Add(field.Name+SEP+strconv.FormatInt(v, 10), id); err != nil {
				i.WarnStatus(field.Name, v, err.Error())
			}
		}
	case string:
		value, _ := field.Value.(string)
		if err := i.invertedIndex.Add(field.Name+SEP+value, id); err != nil {
			i.WarnStatus(field.Name, value, err.Error())
		}
	case int64:
		value, _ := field.Value.(int64)
		if err := i.invertedIndex.Add(field.Name+SEP+strconv.FormatInt(value, 10), id); err != nil {
			i.WarnStatus(field.Name, value, err.Error())
		}
	default:
		i.WarnStatus(field.Name, field.Value, errors.New("the doc is nil or type is wrong").Error())
	}
}

func (i *Indexer) storageAdd(id document.DocId, field *document.Field) {
	if err := i.storageIndex.Add(field.Name, id, field.Value); err != nil {
		i.WarnStatus(field.Name, field.Value, err.Error())
	}
}

func (i *Indexer) invertDel(id document.DocId) {
	if docId, ok := i.campaignMapping.Get(DocId(id)); ok {
		i.bitmap.Del(DocId(docId.(document.DocId)))
	}
}

func (i *Indexer) storageDel(id document.DocId, field *document.Field) {
	if docId, ok := i.campaignMapping.Get(DocId(id)); ok {
		if ok := i.storageIndex.Del(field.Name, docId.(document.DocId)); !ok {
			if i.aDebug != nil {
				i.aDebug.AddDebugMsg(fmt.Sprintf("del [%v] - [%v] failed", id, field))
			}
		}
		i.bitmap.Del(DocId(docId.(document.DocId)))
	}
}

func (i *Indexer) DebugInfo() *debug.Debug {
	if i.aDebug != nil {
		i.aDebug.AddDebugMsg("invert index count: " + strconv.Itoa(i.invertedIndex.Count()))
		i.aDebug.AddDebugMsg("storage index count: " + strconv.Itoa(i.storageIndex.Count()))
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
