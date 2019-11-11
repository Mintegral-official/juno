package query

import (
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/index"
)

type TermQuery struct {
	index.InvertedIterator
}

func NewTermQuery() *TermQuery {
	return &TermQuery{
		index.NewSkipListIterator(index.DEFAULT_MAX_LEVEL, helpers.DocIdFunc),
	}
}

func (t *TermQuery) Next() (document.DocId, error) {
	element := t.InvertedIterator.Next()
	if element != nil {
		v, ok := element.(*index.Element)
		if !ok {
			return 0, helpers.ElementNotfound
		}
		if v == nil {
			return 0, helpers.ElementNotfound
		}
		if v, ok := v.Key().(document.DocId); ok {
			return v, nil
		} else {
			return 0, helpers.DocIdNotFound
		}
	}
	return 0, helpers.ElementNotfound
}

func (t *TermQuery) GetGE(id document.DocId) (document.DocId, error) {
	v := t.InvertedIterator.GetGE(id)
	// fmt.Printf("%T, %v\n", v, v)
	if v != nil {
		k, ok := v.(*index.Element)
		if !ok {
			return 0, helpers.ElementNotfound
		}
		if v, ok := k.Key().(document.DocId); ok {
			return v, nil
		}
		return  0, helpers.DocIdNotFound
	}
	return  0, helpers.ElementNotfound
}

func (t *TermQuery) String() string {
	return ""
}
