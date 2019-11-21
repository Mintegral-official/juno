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
	if iter == nil {
		return nil
	}
	return &TermQuery{
		iterator: iter,
	}
}

func (t *TermQuery) Next() (document.DocId, error) {

	if t.iterator == nil {
		return 0, helpers.DocumentError
	}

	if element := t.iterator.Next(); element != nil {
		v, ok := element.(*datastruct.Element)
		if !ok {
			return 0, helpers.ElementNotfound
		}
		if v == nil {
			return 0, helpers.ElementNotfound
		}
		if v.Key() == nil {
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

	if t.iterator == nil {
		return 0, helpers.DocumentError
	}

	if v := t.iterator.GetGE(id); v != nil {
		v, ok := v.(*datastruct.Element)
		if !ok {
			return 0, helpers.ElementNotfound
		}
		if v == nil {
			return 0, helpers.ElementNotfound
		}
		if v.Key() == nil {
			return 0, helpers.ElementNotfound
		}
		if v, ok := v.Key().(document.DocId); ok {
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
