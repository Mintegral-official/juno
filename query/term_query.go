package query

import (
	"encoding/json"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
)

type TermQuery struct {
	iterator datastruct.Iterator
	aDebug   *debug.Debug
}

func NewTermQuery(iter datastruct.Iterator) *TermQuery {
	tq := &TermQuery{
		aDebug: debug.NewDebug("TermQuery"),
	}
	if iter == nil {
		tq.aDebug.AddDebug("the iterator is nil")
		return tq
	}
	tq.iterator = iter
	return tq
}

func (tq *TermQuery) Next() (document.DocId, error) {

	if tq.iterator == nil {
		return 0, helpers.DocumentError
	}

	tq.iterator.Next()
	if element := tq.iterator.Current(); element != nil {
		v, ok := element.(*datastruct.Element)
		if !ok || v == nil || v.Key() == 0 {
			return 0, helpers.ElementNotfound
		}
		return v.Key(), nil
	}
	return 0, helpers.ElementNotfound
}

func (tq *TermQuery) GetGE(id document.DocId) (document.DocId, error) {

	if tq.iterator == nil {
		return 0, helpers.DocumentError
	}

	if v := tq.iterator.GetGE(id); v != nil {
		v, ok := v.(*datastruct.Element)
		if !ok || v == nil || v.Key() == 0 {
			return 0, helpers.ElementNotfound
		}
		return v.Key(), nil
	}
	return 0, helpers.ElementNotfound
}

func (tq *TermQuery) Current() (document.DocId, error) {
	if tq.iterator == nil {
		return 0, helpers.DocumentError
	}
	if v := tq.iterator.Current(); v != nil {
		v, ok := v.(*datastruct.Element)
		if !ok || v == nil || v.Key() == 0 {
			return 0, helpers.ElementNotfound
		}
		return v.Key(), nil
	}
	return 0, helpers.ElementNotfound
}

func (tq *TermQuery) String() string {
	if res, err := json.Marshal(tq.aDebug); err == nil {
		return string(res)
	} else {
		return err.Error()
	}
}
