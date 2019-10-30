package juno

import (
	"errors"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
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

	return nil
}

func (i *Index) Del(doc *document.DocInfo) error {
	return nil
}

func (i *Index) Update(filename string) error {
	return nil
}

func (i *Index) Dump(filename string) error {
	return nil
}

func (i *Index) Load(filename string) error {
	return nil
}

func (i *Index) Search(query *Query) *index.SearchResult {
	return nil
}
