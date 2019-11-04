package juno

import (
	"encoding/json"
	"errors"
	"github.com/Mintegral-official/juno/document"
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
		return errors.New("doc is nil")
	}
	fileName := "" //TODO
	err := i.invertedIndex.Add(fileName, doc.Id)
	if err != nil {
		return err
	}
	_ = i.Dump(fileName)
	return nil
}

func (i *Index) Del(doc *document.DocInfo) error {
	if doc == nil {
		return errors.New("doc is nil")
	}
	fileName := "" // TODO
	i.invertedIndex.Del(fileName, doc.Id)
	_ = i.Dump(fileName)
	return nil
}

func (i *Index) Update(filename string) error {
    err := i.Dump(filename)
	if err != nil {
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
