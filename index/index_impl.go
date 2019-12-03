package index

import (
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/query"
)

type IndexImpl struct {
	invertedIndex InvertedIndex
	storageIndex  StorageIndex
}

func NewIndex(name string) *IndexImpl {
	return &IndexImpl{
		invertedIndex: NewInvertedIndexImpl(),
		storageIndex:  NewStorageIndexImpl(),
	}

}

func (i *IndexImpl) GetInvertedIndex() InvertedIndex {
	if i == nil || i.invertedIndex == nil {
		return nil
	}
	return i.invertedIndex
}

func (i *IndexImpl) GetStorageIndex() StorageIndex {
	if i == nil || i.storageIndex == nil {
		return nil
	}
	return i.storageIndex
}

func (i *IndexImpl) Add(doc *document.DocInfo) error {
	if doc == nil {
		return helpers.DocumentError
	}
	for j := range doc.Fields {
		var err error
		if doc.Fields[j].IndexType == document.InvertedIndexType {
			err = i.invertedIndex.Add(doc.Fields[j].Name, doc.Id)
		} else if doc.Fields[j].IndexType == document.StorageIndexType {
			err = i.storageIndex.Add(doc.Fields[j].Name, doc.Id, doc.Fields[j].Value)
		} else if doc.Fields[j].IndexType == document.BothIndexType {
			err = i.invertedIndex.Add(doc.Fields[j].Name, doc.Id)
			if err != nil {
				return err
			}
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
		return helpers.DocumentError
	}
	var flag bool
	for j := range doc.Fields {
		if doc.Fields[j].IndexType == document.InvertedIndexType {
			flag = i.invertedIndex.Del(doc.Fields[j].Name, doc.Id)
		} else if doc.Fields[j].IndexType == document.StorageIndexType {
			flag = i.storageIndex.Del(doc.Fields[j].Name, doc.Id)
		} else if doc.Fields[j].IndexType == document.BothIndexType {
			flag = i.invertedIndex.Del(doc.Fields[j].Name, doc.Id)
			if !flag {
				return helpers.DelFailed
			}
			flag = i.storageIndex.Del(doc.Fields[j].Name, doc.Id)
		}
		if !flag {
			return helpers.DelFailed
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

func (i *IndexImpl) GetDataType(fieldName string) document.FieldType {
	return 0
}

func (i *IndexImpl) Search(query query.Query) *SearchResult {
	if query == nil {
		return nil
	}
	s := &SearchResult{Docs: []document.DocId{}}
	if _, err := query.Current(); err != nil {
		return s
	}
	id, err := query.Next()
	for err == nil {
		s.Docs = append(s.Docs, id)
		id, err = query.Next()
	}
	return s
}
