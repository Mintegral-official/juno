package query

import (
	"container/heap"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/query/check"
)

type OrQuery struct {
	checkers []check.Checker
	h        Heap
	aDebug   *debug.Debug
}

func NewOrQuery(queries []Query, checkers []check.Checker) *OrQuery {
	oq := &OrQuery{aDebug: &debug.Debug{
		Name: "NewOrQuery",
		Msg:  []string{},
	},}
	if len(queries) == 0 {
		return oq
	}
	h := &Heap{}
	for i := 0; i < len(queries); i++ {
		heap.Push(h, queries[i])
	}
	oq.h = *h
	oq.checkers = checkers
	return oq
}

func (oq *OrQuery) Next() (document.DocId, error) {
	for target, err := oq.Current(); err == nil; {
		oq.next()
		if oq.check(target) {
			return target, nil
		}
		oq.aDebug.AddDebug(fmt.Sprintf("docID[%d] is filter", target))
		target, err = oq.Current()
	}
	return 0, helpers.NoMoreData
}

func (oq *OrQuery) next() {
	top := oq.h.Top()
	if top != nil {
		q := top.(Query)
		_, _ = q.Next()
		heap.Fix(&oq.h, 0)
	}
}

func (oq *OrQuery) getGE(id document.DocId) {
	top := oq.h.Top()
	if top != nil {
		q := top.(Query)
		_, _ = q.GetGE(id)
		heap.Fix(&oq.h, 0)
	}
}

func (oq *OrQuery) GetGE(id document.DocId) (document.DocId, error) {
	target, err := oq.Current()
	for err == nil && target < id {
		oq.getGE(id)
		target, err = oq.Current()
	}
	for err == nil && !oq.check(target) {
		oq.aDebug.AddDebug(fmt.Sprintf("docID[%d] is not suitable", target))
		target, err = oq.Next()
	}
	return target, err
}

func (oq *OrQuery) Current() (document.DocId, error) {
	top := oq.h.Top()
	if top == nil {
		return 0, helpers.NoMoreData
	}
	q := top.(Query)
	res, err := q.Current()
	if err != nil {
		return 0, err
	}
	if oq.check(res) {
		return res, nil
	}
	oq.aDebug.AddDebug(fmt.Sprintf("docID[%d] is not suitable", res))
	return 0, errors.New(fmt.Sprintf("the result [%d] is not suitable", res))

}

func (oq *OrQuery) String() string {
	for i := 0; i < oq.h.Len(); i++ {
		oq.aDebug.AddDebug(oq.h[i].String())
	}

	if res, err := json.Marshal(oq.aDebug); err == nil {
		return string(res)
	} else {
		return err.Error()
	}
}

func (oq *OrQuery) check(id document.DocId) bool {
	if oq.checkers == nil {
		return true
	}
	for _, v := range oq.checkers {
		if v.Check(id) {
			return true
		}
	}
	return false
}
