package index

import (
	"errors"
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/log"
	"github.com/sirupsen/logrus"
)

type IndexImpl struct {
	invertedIndex   InvertedIndex
	storageIndex    StorageIndex
	campaignMapping map[document.DocId]document.DocId
	bitmap          *datastruct.BitSet
	count           document.DocId
	name            string
	logger          log.Logger
}

func NewIndex(name string) *IndexImpl {
	return &IndexImpl{
		invertedIndex:   NewInvertedIndexImpl(),
		storageIndex:    NewStorageIndexImpl(),
		campaignMapping: make(map[document.DocId]document.DocId, 2000000),
		bitmap:          datastruct.NewBitMap(),
		count:           1,
		name:            name,
		logger:          logrus.New(),
	}

}

func (ii *IndexImpl) GetInvertedIndex() InvertedIndex {
	return ii.invertedIndex
}

func (ii *IndexImpl) GetStorageIndex() StorageIndex {
	return ii.storageIndex
}

func (ii *IndexImpl) GetCampaignMap() map[document.DocId]document.DocId {
	return ii.campaignMapping
}

func (ii *IndexImpl) GetBitMap() *datastruct.BitSet {
	return ii.bitmap
}

func (ii *IndexImpl) GetName() string {
	return ii.name
}

func (ii *IndexImpl) Add(doc *document.DocInfo) error {
	if doc == nil {
		return helpers.DocumentError
	}
	for j := range doc.Fields {
		var err error
		if doc.Fields[j].IndexType == document.InvertedIndexType {
			if err = ii.invertedIndex.Add(doc.Fields[j].Name+"_"+fmt.Sprint(doc.Fields[j].Value), ii.count); err != nil {
				ii.WarnStatus("invert index", doc.Fields[j].Name, fmt.Sprint(doc.Fields[j].Value), err.Error())
				return err
			}
			ii.campaignMapping[doc.Id] = ii.count
			ii.bitmap.Set(uint64(ii.count))
			ii.count++
		} else if doc.Fields[j].IndexType == document.StorageIndexType {
			if err = ii.storageIndex.Add(doc.Fields[j].Name, doc.Id, doc.Fields[j].Value); err != nil {
				ii.WarnStatus("storage index", doc.Fields[j].Name, fmt.Sprint(doc.Fields[j].Value), err.Error())
				return err
			}
		} else if doc.Fields[j].IndexType == document.BothIndexType {
			if err = ii.invertedIndex.Add(doc.Fields[j].Name+"_"+fmt.Sprint(doc.Fields[j].Value), ii.count); err != nil {
				ii.WarnStatus("invert index", doc.Fields[j].Name, fmt.Sprint(doc.Fields[j].Value), err.Error())
				return err
			}
			ii.campaignMapping[doc.Id] = ii.count
			ii.bitmap.Set(uint64(ii.count))
			ii.count++
			if err = ii.storageIndex.Add(doc.Fields[j].Name, doc.Id, doc.Fields[j].Value); err != nil {
				ii.WarnStatus("storage index", doc.Fields[j].Name, fmt.Sprint(doc.Fields[j].Value), err.Error())
				return err
			}
		} else {
			ii.WarnStatus("index", doc.Fields[j].Name, fmt.Sprint(doc.Fields[j].Value),
				fmt.Sprint("index type ", doc.Fields[j].IndexType, " is wrong"))
			return errors.New("the add doc type is wrong or nil ")
		}
	}
	return nil
}

func (ii *IndexImpl) Del(doc *document.DocInfo) {
	if doc == nil {
		return
	}
	for j := range doc.Fields {
		if doc.Fields[j].Value == nil {
			continue
		}
		if doc.Fields[j].IndexType == document.InvertedIndexType {
			ii.invertedIndex.Del(doc.Fields[j].Name+"_"+fmt.Sprint(doc.Fields[j].Value), ii.count)
			ii.bitmap.Del(uint64(ii.count))
		} else if doc.Fields[j].IndexType == document.StorageIndexType {
			ii.storageIndex.Del(doc.Fields[j].Name, doc.Id)
		} else if doc.Fields[j].IndexType == document.BothIndexType {
			ii.invertedIndex.Del(doc.Fields[j].Name+"_"+fmt.Sprint(doc.Fields[j].Value), ii.count)
			ii.bitmap.Del(uint64(ii.count))
			ii.storageIndex.Del(doc.Fields[j].Name, doc.Id)
		} else {
			panic("the del doc type is nil or wrong")
		}
	}
}

func (ii *IndexImpl) Update(filename string) error {
	if err := ii.Dump(filename); err != nil {
		return err
	}
	return nil
}

func (ii *IndexImpl) Dump(filename string) error {
	// TODO
	return nil
}

func (ii *IndexImpl) Load(filename string) error {
	return nil
}

func (ii *IndexImpl) GetDataType(fieldName string) document.FieldType {
	return 0
}

func (ii *IndexImpl) WarnStatus(idxType, name, value, err string) {
	if ii.logger != nil {
		ii.logger.Warnf("[%s]: name:[%s] value:[%s] wrong reason:[%s]", idxType, name, value, err)
	}
}
