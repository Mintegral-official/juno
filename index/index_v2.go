package index

import (
	"errors"
	"fmt"
	"github.com/MintegralTech/juno/datastruct"
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/helpers"
	"github.com/MintegralTech/juno/log"
	"github.com/easierway/concurrent_map"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"sync/atomic"
)

const (
	MaxNumIndex = 50000
)

type IndexerV2 struct {
	invertedIndex   InvertedIndex
	storageIndex    StorageIndex
	campaignMapping *concurrent_map.ConcurrentMap
	idMap           []document.DocId
	count           uint64
	name            string
	kvType          *concurrent_map.ConcurrentMap
	logger          log.Logger
	aDebug          *debug.Debug
}

func NewIndexV2(name string) (i *IndexerV2) {
	i = &IndexerV2{
		invertedIndex:   NewInvertedIndexer(),
		storageIndex:    NewStorageIndexer(),
		campaignMapping: concurrent_map.CreateConcurrentMap(128),
		kvType:          concurrent_map.CreateConcurrentMap(128),
		idMap:           make([]document.DocId, MaxNumIndex),
		count:           0,
		name:            name,
		logger:          logrus.New(),
	}
	return i
}

func (i *IndexerV2) GetInvertedIndex() InvertedIndex {
	return i.invertedIndex
}

func (i *IndexerV2) GetStorageIndex() StorageIndex {
	return i.storageIndex
}

func (i *IndexerV2) GetCampaignMap() *concurrent_map.ConcurrentMap {
	return i.campaignMapping
}

func (i *IndexerV2) GetName() string {
	return i.name
}

func (i *IndexerV2) SetDebug(level int) {
	if i.aDebug == nil {
		i.aDebug = debug.NewDebug(level, i.GetName())
		i.invertedIndex.SetDebug(level)
		i.storageIndex.SetDebug(level)
	}
}

func (i *IndexerV2) GetValueById(id document.DocId) [2]map[string][]string {
	var res [2]map[string][]string
	docId, ok := i.campaignMapping.Get(DocId(id))
	if ok {
		res[0] = i.GetInvertedIndex().GetValueById(docId.(document.DocId))
		res[1] = i.GetStorageIndex().GetValueById(docId.(document.DocId))
	}
	return res
}

func (i *IndexerV2) UpdateIds(fieldName string, ids []document.DocId) {
	panic("method not support")
}

func (i *IndexerV2) Delete(fieldName string) {
	panic("method not support")
}

func (i *IndexerV2) Add(doc *document.DocInfo) error {
	if doc == nil {
		return helpers.DocumentError
	}
	if i.count >= MaxNumIndex {
		return fmt.Errorf("index is full: current[%d], max[%d]", i.count, MaxNumIndex)
	}
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
		}
	}
	i.campaignMapping.Set(DocId(doc.Id), document.DocId(i.count))
	i.idMap[i.count] = doc.Id
	atomic.AddUint64(&i.count, 1)
	return nil
}

func (i *IndexerV2) Del(doc *document.DocInfo) {
	i.campaignMapping.Del(DocId(doc.Id))
}

func (i *IndexerV2) Update(filename string) error {
	if err := i.Dump(filename); err != nil {
		return err
	}
	return nil
}

func (i *IndexerV2) Dump(filename string) error {
	// TODO
	return nil
}

func (i *IndexerV2) Load(filename string) error {
	return nil
}

func (i *IndexerV2) GetDataType(fieldName string) document.FieldType {
	if t, ok := i.kvType.Get(concurrent_map.StrKey(fieldName)); ok {
		return t.(document.FieldType)
	}
	return document.DefaultFieldType
}

func (i *IndexerV2) invertAdd(id document.DocId, field *document.Field) {
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

func (i *IndexerV2) storageAdd(id document.DocId, field *document.Field) {
	if err := i.storageIndex.Add(field.Name, id, field.Value); err != nil {
		i.WarnStatus(field.Name, field.Value, err.Error())
	}
}

func (i *IndexerV2) DebugInfo() *debug.Debug {
	if i.aDebug != nil {
		i.aDebug.AddDebugMsg("invert index count: " + strconv.Itoa(i.invertedIndex.Count()))
		i.aDebug.AddDebugMsg("storage index count: " + strconv.Itoa(i.storageIndex.Count()))
		return i.aDebug
	}
	return nil
}

func (i *IndexerV2) WarnStatus(name string, value interface{}, err string) {
	if i.logger != nil {
		i.logger.Warnf("name:[%s] value:[%v] wrong reason:[%s]", name, value, err)
	}
}

func (i *IndexerV2) MergeIndex(target *IndexerV2) error {

	invertIters := make(map[string]datastruct.Iterator, target.invertedIndex.Count())
	target.invertedIndex.Range(func(key, value interface{}) bool {
		k := key.(string)
		items := strings.Split(k, SEP)
		iter := target.invertedIndex.Iterator(items[0], items[1])
		invertIters[k] = iter
		return true
	})

	storageIters := make(map[string]datastruct.Iterator, target.storageIndex.Count())
	target.storageIndex.Range(func(key, value interface{}) bool {
		k := key.(string)
		iter := target.storageIndex.Iterator(k)
		storageIters[k] = iter
		return true
	})

	// merge by id
	for id := uint64(0); id < target.count; id++ {
		docId := target.idMap[id]
		// already deleted
		if _, ok := target.campaignMapping.Get(DocId(docId)); !ok {
			continue
		}

		// new index updated
		if _, ok := i.campaignMapping.Get(DocId(docId)); ok {
			continue
		}

		// invert List
		for k, v := range invertIters {
			if id == uint64(v.Current().Key()) {
				// add invert index
				if e := i.invertedIndex.Add(k, document.DocId(i.count)); e != nil {
					i.logger.Warnf("MergeIndex add inverted index error, docId[%d], id[%d]", docId, i.count)
				}
				v.Next()
				continue
			}
		}

		// storage List
		for k, v := range storageIters {
			if id == uint64(v.Current().Key()) {
				// add storage index
				if e := i.storageIndex.Add(k, document.DocId(i.count), v.Current().Value()); e != nil {
					i.logger.Warnf("MergeIndex add inverted index error, docId[%d], id[%d]", docId, i.count)
				}
				v.Next()
				continue
			}
		}
		i.campaignMapping.Set(DocId(docId), document.DocId(i.count))
		i.idMap[i.count] = docId
		i.count++

		if i.count > MaxNumIndex {
			return fmt.Errorf("merge index error, index is full, maxsize[%d], current[%d]", MaxNumIndex, i.count)
		}
	}

	return nil
}

func (i *IndexerV2) GetId(id document.DocId) (document.DocId, error) {
	if uint64(id) >= i.count {
		return 0, errors.New("id not found")
	}
	return i.idMap[id], nil
}

func (i *IndexerV2) GetInnerId(id document.DocId) (document.DocId, error) {
	v, ok := i.GetCampaignMap().Get(DocId(id))
	if !ok {
		return 0, errors.New("id not found")
	}
	return v.(document.DocId), nil
}

func (i *IndexerV2) IndexInfo() string {
	var builder strings.Builder
	builder.WriteString("index[")
	builder.WriteString(strconv.FormatInt(int64(i.count), 10))
	builder.WriteString("], invertIndex[")
	builder.WriteString(strconv.Itoa(i.GetInvertedIndex().Count()))
	builder.WriteString("], storageIndex[")
	builder.WriteString(strconv.Itoa(i.GetStorageIndex().Count()))
	builder.WriteString("]")
	return builder.String()
}
