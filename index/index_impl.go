package index

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/log"
	"sync"

	//cmap "github.com/easierway/concurrent_map"
	"github.com/sirupsen/logrus"
	"strconv"
)

type Indexer struct {
	invertedIndex   InvertedIndex
	storageIndex    StorageIndex
	campaignMapping *sync.Map
	bitmap          *datastruct.BitSet
	count           document.DocId
	name            string
	logger          log.Logger
	aDebug          *debug.Debug
}

func NewIndex(name string) *Indexer {
	return &Indexer{
		invertedIndex:   NewInvertedIndexer(),
		storageIndex:    NewStorageIndexer(),
		campaignMapping: &sync.Map{},
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

func (i *Indexer) Add(doc *document.DocInfo) error {
	if doc == nil {
		return helpers.DocumentError
	}
	for j := range doc.Fields {
		var err error
		if doc.Fields[j].IndexType == document.InvertedIndexType {
			if err = i.invertedIndex.Add(doc.Fields[j].Name+"_"+fmt.Sprint(doc.Fields[j].Value), i.count); err != nil {
				i.WarnStatus("invert index", doc.Fields[j].Name, fmt.Sprint(doc.Fields[j].Value), err.Error())
				return err
			}
			i.campaignMapping.Store(doc.Id, i.count)
			i.bitmap.Set(i.count)
			i.count++
		} else if doc.Fields[j].IndexType == document.StorageIndexType {
			if err = i.storageIndex.Add(doc.Fields[j].Name, doc.Id, doc.Fields[j].Value); err != nil {
				i.WarnStatus("storage index", doc.Fields[j].Name, fmt.Sprint(doc.Fields[j].Value), err.Error())
				return err
			}
		} else if doc.Fields[j].IndexType == document.BothIndexType {
			if err = i.invertedIndex.Add(doc.Fields[j].Name+"_"+fmt.Sprint(doc.Fields[j].Value), i.count); err != nil {
				i.WarnStatus("invert index", doc.Fields[j].Name, fmt.Sprint(doc.Fields[j].Value), err.Error())
				return err
			}
			i.campaignMapping.Store(doc.Id, i.count)
			i.bitmap.Set(i.count)
			i.count++
			if err = i.storageIndex.Add(doc.Fields[j].Name, doc.Id, doc.Fields[j].Value); err != nil {
				i.WarnStatus("storage index", doc.Fields[j].Name, fmt.Sprint(doc.Fields[j].Value), err.Error())
				return err
			}
		} else {
			i.WarnStatus("index", doc.Fields[j].Name, fmt.Sprint(doc.Fields[j].Value),
				fmt.Sprint("index type ", doc.Fields[j].IndexType, " is wrong"))
			return errors.New("the add doc type is wrong or nil ")
		}
	}
	return nil
}

func (i *Indexer) Del(doc *document.DocInfo) {
	if doc == nil {
		return
	}
	for j := range doc.Fields {
		if doc.Fields[j].Value == nil {
			continue
		}
		if doc.Fields[j].IndexType == document.InvertedIndexType {
			i.invertedIndex.Del(doc.Fields[j].Name+"_"+fmt.Sprint(doc.Fields[j].Value), i.count)
			i.bitmap.Del(i.count)
		} else if doc.Fields[j].IndexType == document.StorageIndexType {
			i.storageIndex.Del(doc.Fields[j].Name, doc.Id)
		} else if doc.Fields[j].IndexType == document.BothIndexType {
			i.invertedIndex.Del(doc.Fields[j].Name+"_"+fmt.Sprint(doc.Fields[j].Value), i.count)
			i.bitmap.Del(i.count)
			i.storageIndex.Del(doc.Fields[j].Name, doc.Id)
		} else {
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

func (i *Indexer) String() string {
	i.aDebug.AddDebug("invert index count: " + strconv.Itoa(i.invertedIndex.Count()) +
		", storage index count: " + strconv.Itoa(i.storageIndex.Count()))
	var (
		invert  = debug.NewDebug("invert index")
		storage = debug.NewDebug("storage index")
	)
	invert.AddDebug(i.invertedIndex.String())
	storage.AddDebug(i.storageIndex.String())
	invert.Node, i.aDebug.Node = storage, invert
	if res, err := json.Marshal(i.aDebug); err == nil {
		return string(res)
	} else {
		return err.Error()
	}
}

func (i *Indexer) GetDataType(fieldName string) document.FieldType {
	return 0
}

func (i *Indexer) WarnStatus(idxType, name, value, err string) {
	if i.logger != nil {
		i.logger.Warnf("[%s]: name:[%s] value:[%s] wrong reason:[%s]", idxType, name, value, err)
	}
}
