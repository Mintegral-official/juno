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
		//v = datastruct.ElementCopy(v)
		//t.iterator.Next()
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
	if v != nil {
		k, ok := v.(*datastruct.Element)
		if !ok {
			return 0, helpers.ElementNotfound
		}
		if k == nil {
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
	v := t.iterator.Current()
	if v == nil {
		return 0, helpers.ElementNotfound
	}
	res, ok := v.(*datastruct.Element)
	if !ok {
		return 0, helpers.ElementNotfound
	}

	if res == nil {
		return 0, helpers.ElementNotfound
	}

	return res.Key().(document.DocId), nil
}

func (t *TermQuery) String() string {
	return ""
}
