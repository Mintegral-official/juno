package query

import (
	"container/heap"
	"errors"
	"fmt"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/query/check"
)

type OrQuery struct {
	checkers []check.Checker
	h        Heap
}

func NewOrQuery(queries []Query, checkers []check.Checker) *OrQuery {
	if len(queries) == 0 {
		return &OrQuery{}
	}
	h := &Heap{}
	for i := 0; i < len(queries); i++ {
		heap.Push(h, queries[i])
	}
	return &OrQuery{
		checkers: checkers,
		h:        *h,
	}
}

func (oq *OrQuery) Next() (document.DocId, error) {
	for target, err := oq.Current(); err == nil; {
		oq.next()
		if oq.check(target) {
			return target, nil
		}
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
	return 0, errors.New(fmt.Sprintf("the result [%d] is filter", res))

}

func (oq *OrQuery) String() string {
	return ""
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
