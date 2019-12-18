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

func (tq *TermQuery) Next() (document.DocId, error) {

	if tq.iterator == nil {
		return 0, helpers.DocumentError
	}

	if element := tq.iterator.Next(); element != nil {
		v, ok := element.(*datastruct.Element)
		if !ok {
			return 0, helpers.ElementNotfound
		}
		if v == nil || v.Key() == nil {
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

func (tq *TermQuery) GetGE(id document.DocId) (document.DocId, error) {

	if tq.iterator == nil {
		return 0, helpers.DocumentError
	}

	if v := tq.iterator.GetGE(id); v != nil {
		v, ok := v.(*datastruct.Element)
		if !ok || v == nil || v.Key() == nil {
			return 0, helpers.ElementNotfound
		}
		if v, ok := v.Key().(document.DocId); ok {
			return v, nil
		}
		return 0, helpers.DocIdNotFound
	}
	return 0, helpers.ElementNotfound
}

func (tq *TermQuery) Current() (document.DocId, error) {
	v := tq.iterator.Current()
	if v == nil {
		return 0, helpers.ElementNotfound
	}
	res, ok := v.(*datastruct.Element)
	if !ok {
		return 0, helpers.ElementNotfound
	}
	if res, ok := res.Key().(document.DocId); ok {
		return res, nil
	}
	return 0, helpers.ElementNotfound
}

func (tq *TermQuery) String() string {
	return ""
}
