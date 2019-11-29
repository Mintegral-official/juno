package query

import (
	"container/heap"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
)

type OrQuery struct {
	checkers []Checker
	h        Heap
}

func NewOrQuery(querys []Query, checkers []Checker) *OrQuery {
	if querys == nil {
		return nil
	}
	h := &Heap{}
	for i := 0; i < len(querys); i++ {
		heap.Push(h, querys[i])
	}
	return &OrQuery{
		checkers: checkers,
		h:        *h,
	}
}

func (o *OrQuery) Next() (document.DocId, error) {
	for target, err := o.Current(); err == nil; {
		o.next()
		if o.check(target) {
			return target, nil
		}
		target, err = o.Current()
	}
	return 0, helpers.NoMoreData
}

func (o *OrQuery) next() {
	top := o.h.Top()
	if top != nil {
		q := top.(Query)
		_, _ = q.Next()
		heap.Fix(&o.h, 0)
	}
}

func (o *OrQuery) getGE(id document.DocId) {
	top := o.h.Top()
	if top != nil {
		q := top.(Query)
		_, _ = q.GetGE(id)
		heap.Fix(&o.h, 0)
	}
}

func (o *OrQuery) GetGE(id document.DocId) (document.DocId, error) {
	target, err := o.Current()
	for err == nil && target < id {
		o.getGE(id)
		target, err = o.Current()
	}
	for err == nil && !o.check(target) {
		target, err = o.Next()
	}
	return target, err
}

func (o *OrQuery) Current() (document.DocId, error) {
	top := o.h.Top()
	if top == nil {
		return 0, helpers.NoMoreData
	}
	q := top.(Query)
	return q.Current()
}

func (t *OrQuery) String() string {
	return ""
}

func (o *OrQuery) check(id document.DocId) bool {
	if o.checkers == nil {
		return true
	}
	for _, v := range o.checkers {
		if v.Check(id) {
			return true
		}
	}
	return false
}
