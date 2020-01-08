package index

import (
	"errors"
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/log"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
)

type Indexer struct {
	invertedIndex   InvertedIndex
	storageIndex    StorageIndex
	campaignMapping *sync.Map
	bitmap          *datastruct.BitSet
	count           document.DocId
	name            string
	kvType          *sync.Map
	logger          log.Logger
	aDebug          *debug.Debug
}

func NewIndex(name string) *Indexer {
	return &Indexer{
		invertedIndex:   NewInvertedIndexer(),
		storageIndex:    NewStorageIndexer(),
		campaignMapping: &sync.Map{},
		kvType:          &sync.Map{},
		bitmap:          datastruct.NewBitMap(),
		count:           1,
		name:            name,
		logger:          logrus.New(),
		aDebug:          debug.NewDebug(name),
	}

}

func (i *Indexer) GetInvertedIndex() InvertedIndex {
	return i.invertedIndex
}

func (i *Indexer) GetStorageIndex() StorageIndex {
	return i.storageIndex
}

func (i *Indexer) GetCampaignMap() *sync.Map {
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
				i.aDebug.AddDebugMsg(fmt.Sprintf("invert index add failed: docId[%d], Name[%s], value[%s], reason[%s]",
					int(doc.Id), field.Name, fmt.Sprint(field.Value), err.Error()))
				return err
			}
		case document.StorageIndexType:
			err = i.storageAdd(doc.Id, field)
			if err != nil {
				i.aDebug.AddDebugMsg(fmt.Sprintf("storage index add failed: docId[%d], Name[%s], value[%s], reason[%s]",
					int(doc.Id), field.Name, fmt.Sprint(field.Value), err.Error()))
				return err
			}
			i.kvType.Store(field.Name, field.ValueType)
		case document.BothIndexType:
			err = i.invertAdd(doc.Id, field)
			if err != nil {
				i.aDebug.AddDebugMsg(fmt.Sprintf("invert index add failed: docId[%d], Name[%s], value[%s], reason[%s]",
					int(doc.Id), field.Name, fmt.Sprint(field.Value), err.Error()))
				return err
			}
			err = i.storageAdd(doc.Id, field)
			if err != nil {
				i.aDebug.AddDebugMsg(fmt.Sprintf("storage index add failed: docId[%d], Name[%s], value[%s], reason[%s]",
					int(doc.Id), field.Name, fmt.Sprint(field.Value), err.Error()))
				return err
			}
			i.kvType.Store(field.Name, field.ValueType)
		default:
			i.WarnStatus("index", field.Name,
				fmt.Sprint(field.Value), fmt.Sprint("index type ", field.IndexType, " is wrong"))
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

func (i *Indexer) DebugInfo() *debug.Debug {
	i.aDebug.AddDebugMsg("invert index count: " + strconv.Itoa(i.invertedIndex.Count()))
	i.aDebug.AddDebugMsg("storage index count: " + strconv.Itoa(i.storageIndex.Count()))
	i.aDebug.AddDebug(i.invertedIndex.DebugInfo(), i.storageIndex.DebugInfo())
	return i.aDebug
}

func (i *Indexer) GetDataType(fieldName string) document.FieldType {
	if t, ok := i.kvType.Load(fieldName); ok {
		return t.(document.FieldType)
	}
	return document.DefaultFieldType
}

func (i *Indexer) invertAdd(id document.DocId, field *document.Field) (err error) {
	if v, ok := field.Value.([]string); ok {
		for _, s := range v {
			if err = i.invertedIndex.Add(field.Name+"_"+s, id); err != nil {
				i.WarnStatus("invert index", field.Name, s, err.Error())
				return err
			}
			i.campaignMapping.Store(id, i.count)
			i.bitmap.Set(i.count)
			i.count++
		}
	} else if v, ok := field.Value.([]int64); ok {
		for _, s := range v {
			if err = i.invertedIndex.Add(field.Name+"_"+fmt.Sprint(s), id); err != nil {
				i.WarnStatus("invert index", field.Name, fmt.Sprint(s), err.Error())
				return err
			}
			i.campaignMapping.Store(id, i.count)
			i.bitmap.Set(i.count)
			i.count++
		}
	} else if v, ok := field.Value.(string); ok {
		if err = i.invertedIndex.Add(field.Name+"_"+v, id); err != nil {
			i.WarnStatus("invert index", field.Name, v, err.Error())
			return err
		}
		i.campaignMapping.Store(id, i.count)
		i.bitmap.Set(i.count)
		i.count++
	} else if v, ok := field.Value.(int64); ok {
		if err = i.invertedIndex.Add(field.Name+"_"+fmt.Sprint(v), id); err != nil {
			i.WarnStatus("invert index", field.Name, fmt.Sprint(v), err.Error())
			return err
		}
		i.campaignMapping.Store(id, i.count)
		i.bitmap.Set(i.count)
		i.count++
	} else {
		return errors.New("the doc is nil or type is wrong")
	}
	return nil
}

func (i *Indexer) storageAdd(id document.DocId, field *document.Field) (err error) {
	if err = i.storageIndex.Add(field.Name, id, field.Value); err != nil {
		i.WarnStatus("storage index", field.Name, fmt.Sprint(field.Value), err.Error())
		return err
	}
	return nil
}

func (i *Indexer) invertDel(id document.DocId, field *document.Field) {
	if v, ok := field.Value.([]string); ok {
		for _, s := range v {
			i.invertedIndex.Del(field.Name+"_"+s, id)
			v, _ := i.GetCampaignMap().Load(id)
			i.bitmap.Del(v.(document.DocId))
		}
	} else if v, ok := field.Value.([]int64); ok {
		for _, s := range v {
			i.invertedIndex.Del(field.Name+"_"+fmt.Sprint(s), id)
			v, _ := i.GetCampaignMap().Load(id)
			i.bitmap.Del(v.(document.DocId))
		}
	} else if v, ok := field.Value.(string); ok {
		i.invertedIndex.Del(field.Name+"_"+v, id)
		v, _ := i.GetCampaignMap().Load(id)
		i.bitmap.Del(v.(document.DocId))
	} else if v, ok := field.Value.(int64); ok {
		i.invertedIndex.Del(field.Name+"_"+fmt.Sprint(v), id)
		v, _ := i.GetCampaignMap().Load(id)
		i.bitmap.Del(v.(document.DocId))
	}
}

func (i *Indexer) storageDel(id document.DocId, field *document.Field) {
	i.storageIndex.Del(field.Name, id)
}

func (i *Indexer) WarnStatus(idxType, name, value, err string) {
	if i.logger != nil {
		i.logger.Warnf("[%s]: name:[%s] value:[%s] wrong reason:[%s]", idxType, name, value, err)
	}
}
