package query

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
)

type TermQuery struct {
	iterator datastruct.Iterator
}

func NewTermQuery(iter datastruct.Iterator) *TermQuery {
	return &TermQuery{
		iterator: iter,
	}
}

func (t *TermQuery) Next() (document.DocId, error) {
	element := t.iterator.Next()
	if element != nil {
		v, ok := element.(*datastruct.Element)
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
	v := t.iterator.GetGE(id)
	// fmt.Printf("%T, %v\n", v, v)
	if v != nil {
		k, ok := v.(*datastruct.Element)
		if !ok {
			return 0, helpers.ElementNotfound
		}
		if v, ok := k.Key().(document.DocId); ok {
			return v, nil
		}
		return 0, helpers.DocIdNotFound
	}
	return 0, helpers.ElementNotfound
}

func (t *TermQuery) Current() (document.DocId, error) {
	v, err := t.iterator.Current()
	if err != nil {
		return 0, err
	}
	return v.(document.DocId), err
}

func (t *TermQuery) String() string {
	return ""
}
