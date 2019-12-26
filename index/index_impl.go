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
	"github.com/sirupsen/logrus"
	"strconv"
)

type Indexer struct {
	invertedIndex   InvertedIndex
	storageIndex    StorageIndex
	campaignMapping map[document.DocId]document.DocId
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
		campaignMapping: make(map[document.DocId]document.DocId, 2000000),
		bitmap:          datastruct.NewBitMap(),
		count:           1,
		name:            name,
		logger:          logrus.New(),
		aDebug:          debug.NewDebug(name),
	}

}

func (indexer *Indexer) GetInvertedIndex() InvertedIndex {
	return indexer.invertedIndex
}

func (indexer *Indexer) GetStorageIndex() StorageIndex {
	return indexer.storageIndex
}

func (indexer *Indexer) GetCampaignMap() map[document.DocId]document.DocId {
	return indexer.campaignMapping
}

func (indexer *Indexer) GetBitMap() *datastruct.BitSet {
	return indexer.bitmap
}

func (indexer *Indexer) GetName() string {
	return indexer.name
}

func (indexer *Indexer) Add(doc *document.DocInfo) error {
	if doc == nil {
		return helpers.DocumentError
	}
	for j := range doc.Fields {
		var err error
		if doc.Fields[j].IndexType == document.InvertedIndexType {
			if err = indexer.invertedIndex.Add(doc.Fields[j].Name+"_"+fmt.Sprint(doc.Fields[j].Value), indexer.count); err != nil {
				indexer.WarnStatus("invert index", doc.Fields[j].Name, fmt.Sprint(doc.Fields[j].Value), err.Error())
				return err
			}
			indexer.campaignMapping[doc.Id] = indexer.count
			indexer.bitmap.Set(uint64(indexer.count))
			indexer.count++
		} else if doc.Fields[j].IndexType == document.StorageIndexType {
			if err = indexer.storageIndex.Add(doc.Fields[j].Name, doc.Id, doc.Fields[j].Value); err != nil {
				indexer.WarnStatus("storage index", doc.Fields[j].Name, fmt.Sprint(doc.Fields[j].Value), err.Error())
				return err
			}
		} else if doc.Fields[j].IndexType == document.BothIndexType {
			if err = indexer.invertedIndex.Add(doc.Fields[j].Name+"_"+fmt.Sprint(doc.Fields[j].Value), indexer.count); err != nil {
				indexer.WarnStatus("invert index", doc.Fields[j].Name, fmt.Sprint(doc.Fields[j].Value), err.Error())
				return err
			}
			indexer.campaignMapping[doc.Id] = indexer.count
			indexer.bitmap.Set(uint64(indexer.count))
			indexer.count++
			if err = indexer.storageIndex.Add(doc.Fields[j].Name, doc.Id, doc.Fields[j].Value); err != nil {
				indexer.WarnStatus("storage index", doc.Fields[j].Name, fmt.Sprint(doc.Fields[j].Value), err.Error())
				return err
			}
		} else {
			indexer.WarnStatus("index", doc.Fields[j].Name, fmt.Sprint(doc.Fields[j].Value),
				fmt.Sprint("index type ", doc.Fields[j].IndexType, " is wrong"))
			return errors.New("the add doc type is wrong or nil ")
		}
	}
	return nil
}

func (indexer *Indexer) Del(doc *document.DocInfo) {
	if doc == nil {
		return
	}
	for j := range doc.Fields {
		if doc.Fields[j].Value == nil {
			continue
		}
		if doc.Fields[j].IndexType == document.InvertedIndexType {
			indexer.invertedIndex.Del(doc.Fields[j].Name+"_"+fmt.Sprint(doc.Fields[j].Value), indexer.count)
			indexer.bitmap.Del(uint64(indexer.count))
		} else if doc.Fields[j].IndexType == document.StorageIndexType {
			indexer.storageIndex.Del(doc.Fields[j].Name, doc.Id)
		} else if doc.Fields[j].IndexType == document.BothIndexType {
			indexer.invertedIndex.Del(doc.Fields[j].Name+"_"+fmt.Sprint(doc.Fields[j].Value), indexer.count)
			indexer.bitmap.Del(uint64(indexer.count))
			indexer.storageIndex.Del(doc.Fields[j].Name, doc.Id)
		} else {
			panic("the del doc type is nil or wrong")
		}
	}
}

func (indexer *Indexer) Update(filename string) error {
	if err := indexer.Dump(filename); err != nil {
		return err
	}
	return nil
}

func (indexer *Indexer) Dump(filename string) error {
	// TODO
	return nil
}

func (indexer *Indexer) Load(filename string) error {
	return nil
}

func (indexer *Indexer) String() string {
	indexer.aDebug.AddDebug("invert index count: " + strconv.Itoa(indexer.invertedIndex.Count()) +
		", storage index count: " + strconv.Itoa(indexer.storageIndex.Count()))
	var (
		invert  = debug.NewDebug("invert index")
		storage = debug.NewDebug("storage index")
	)
	invert.AddDebug(indexer.invertedIndex.String())
	storage.AddDebug(indexer.storageIndex.String())
	invert.Node, indexer.aDebug.Node = storage, invert
	if res, err := json.Marshal(indexer.aDebug); err == nil {
		return string(res)
	} else {
		return err.Error()
	}
}

func (indexer *Indexer) GetDataType(fieldName string) document.FieldType {
	return 0
}

func (indexer *Indexer) WarnStatus(idxType, name, value, err string) {
	if indexer.logger != nil {
		indexer.logger.Warnf("[%s]: name:[%s] value:[%s] wrong reason:[%s]", idxType, name, value, err)
	}
}
