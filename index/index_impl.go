package index

import (
	"errors"
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/query"
	"time"
)

type IndexImpl struct {
	invertedIndex   InvertedIndex
	storageIndex    StorageIndex
	campaignMapping map[document.DocId]document.DocId
	bitmap          *datastruct.BitMap
	count           document.DocId
}

func NewIndex(name string) *IndexImpl {
	return &IndexImpl{
		invertedIndex:   NewInvertedIndexImpl(),
		storageIndex:    NewStorageIndexImpl(),
		campaignMapping: make(map[document.DocId]document.DocId, 2000000),
		bitmap:          datastruct.NewBitMap(2000000),
		count:           1,
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

func (ii *IndexImpl) GetBitMap() *datastruct.BitMap {
	return ii.bitmap
}

func (ii *IndexImpl) Add(doc *document.DocInfo) error {
	if doc == nil {
		return helpers.DocumentError
	}
	for j := range doc.Fields {
		var err error
		if doc.Fields[j].IndexType == document.InvertedIndexType {
			ii.campaignMapping[doc.Id] = ii.count
			if err = ii.invertedIndex.Add(doc.Fields[j].Name+"_"+fmt.Sprint(doc.Fields[j].Value), ii.count); err != nil {
				return err
			}
			if !ii.bitmap.IsExist(int(ii.count)) {
				ii.bitmap.Set(int(ii.count))
			}
			ii.count++
		} else if doc.Fields[j].IndexType == document.StorageIndexType {
			if err = ii.storageIndex.Add(doc.Fields[j].Name, doc.Id, doc.Fields[j].Value); err != nil {
				return err
			}
		} else if doc.Fields[j].IndexType == document.BothIndexType {
			ii.campaignMapping[doc.Id] = ii.count
			if err = ii.invertedIndex.Add(doc.Fields[j].Name+"_"+fmt.Sprint(doc.Fields[j].Value), ii.count); err != nil {
				return err
			}
			if !ii.bitmap.IsExist(int(ii.count)) {
				ii.bitmap.Set(int(ii.count))
			}
			ii.count++
			if err = ii.storageIndex.Add(doc.Fields[j].Name, doc.Id, doc.Fields[j].Value); err != nil {
				return err
			}
		} else {
			return errors.New("the add doc type is nil or wrong")
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
			ii.bitmap.Del(int(ii.count))
		} else if doc.Fields[j].IndexType == document.StorageIndexType {
			ii.storageIndex.Del(doc.Fields[j].Name, doc.Id)
		} else if doc.Fields[j].IndexType == document.BothIndexType {
			ii.invertedIndex.Del(doc.Fields[j].Name+"_"+fmt.Sprint(doc.Fields[j].Value), ii.count)
			ii.bitmap.Del(int(ii.count))
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

func (ii *IndexImpl) Search(query query.Query) *SearchResult {
	if query == nil {
		return nil
	}
	s, now := &SearchResult{Docs: []document.DocId{}}, time.Now()
	if _, err := query.Current(); err != nil {
		return s
	}
	id, err := query.Next()
	for err == nil {
		if !ii.GetBitMap().IsExist(int(ii.GetCampaignMap()[id])) {
			continue
		}
		s.Docs = append(s.Docs, id)
		id, err = query.Next()
	}
	s.Time = time.Since(now)
	return s
}
