package juno

import (
	"encoding/json"
	"errors"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/index"
	"io/ioutil"
)

type Index struct {
	invertedIndex index.InvertedIndex
	storageIndex  index.StorageIndex
}

func NewIndex(name string) *Index {
	return &Index{
		invertedIndex: index.NewSimpleInvertedIndex(),
		storageIndex:  index.NewSimpleStorageIndex(),
	}

}

func (i *Index) Add(doc *document.DocInfo) error {
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

		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Index) Del(doc *document.DocInfo) error {
	if doc == nil {
		return errors.New("doc is nil")
	}
	var flag bool
	for j := range doc.Fields {
		if doc.Fields[j].IndexType == document.INVERTED_INDEX_TYPE {
			flag = i.invertedIndex.Del(doc.Fields[j].Name, doc.Id)
		}

		if doc.Fields[j].IndexType == document.STORAGE_INDEX_TYPE {
			flag = i.storageIndex.Del(doc.Fields[j].Name, doc.Id)
		}
		if !flag {
			return errors.New("del failed")
		}
	}

	return nil
}

func (i *Index) Update(filename string) error {
	if err := i.Dump(filename); err != nil {
		return err
	}
	return nil
}

func (i *Index) Dump(filename string) error {
	data, err := json.Marshal(i)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0x77)
}

func (i *Index) Load(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, i)
}

func (i *Index) Search(query *Query) *index.SearchResult {
	return nil
}
