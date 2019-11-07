package index

import (
	"encoding/json"
	"errors"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
)

type IndexImpl struct {
	invertedIndex InvertedIndex
	storageIndex  StorageIndex
}

func NewIndex(name string) *IndexImpl {
	return &IndexImpl{
		invertedIndex: NewSimpleInvertedIndex(),
		storageIndex:  NewSimpleStorageIndex(),
	}

}

func (i *IndexImpl) Add(doc *document.DocInfo) error {
	if doc == nil {
		return helpers.DOCUMENT_ERROR
	}
	for j := range doc.Fields {
		var err error
		if doc.Fields[j].IndexType == document.INVERTED_INDEX_TYPE {
			err = i.invertedIndex.Add(doc.Fields[j].Name, doc.Id)
		}

		if doc.Fields[j].IndexType == document.STORAGE_INDEX_TYPE {
			err = i.storageIndex.Add(doc.Fields[j].Name, doc.Id, doc.Fields[j].Value)
		}

		if doc.Fields[j].IndexType == document.INDEX_TYPE {
			err = i.invertedIndex.Add(doc.Fields[j].Name, doc.Id)
			err = i.storageIndex.Add(doc.Fields[j].Name, doc.Id, doc.Fields[j].Value)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (i *IndexImpl) Del(doc *document.DocInfo) error {
	if doc == nil {
		return helpers.DOCUMENT_ERROR
	}
	var flag bool
	for j := range doc.Fields {
		if doc.Fields[j].IndexType == document.INVERTED_INDEX_TYPE {
			flag = i.invertedIndex.Del(doc.Fields[j].Name, doc.Id)
		}

		if doc.Fields[j].IndexType == document.STORAGE_INDEX_TYPE {
			flag = i.storageIndex.Del(doc.Fields[j].Name, doc.Id)
		}
		if doc.Fields[j].IndexType == document.INDEX_TYPE {
			flag = i.invertedIndex.Del(doc.Fields[j].Name, doc.Id)
			flag = i.storageIndex.Del(doc.Fields[j].Name, doc.Id)
		}
		if !flag {
			return helpers.DEL_FAILED
		}
	}

	return nil
}

func (i *IndexImpl) Update(filename string) error {
	if err := i.Dump(filename); err != nil {
		return err
	}
	return nil
}

func (i *IndexImpl) Dump(filename string) error {
	// TODO
	return nil
}

func (i *IndexImpl) Load(filename string) error {
	return nil
}
