package query

import (
	"container/heap"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/pkg/errors"
)

type OrQuery struct {
	querys   []Query
	checkers []Checker
	h        Heap
	curIdx   int
}

func NewOrQuery(querys []Query, checkers []Checker) *OrQuery {
	if querys == nil {
		return nil
	}
	h := &Heap{}
	// heap.Init(h)
	for i := 0; i < len(querys); i++ {
		heap.Push(h, querys[i])
	}
	return &OrQuery{
		querys:   querys,
		checkers: checkers,
		h:        *h,
	}
}

func (o *OrQuery) Next() (document.DocId, error) {

	top := o.h.Top()
	if top == nil {
		return 0, helpers.NoMoreData
	}
	q := top.(Query)
	target, err := q.Current()
	// fmt.Println(target)
	if err != nil {
		return 0, errors.Wrap(err, "no more data")
	}

	_, _ = q.Next()
	heap.Fix(&o.h, 0)

	if o.check(target) {
		return target, nil
	} else {
		return o.Next()
	}

}

func (o *OrQuery) GetGE(id document.DocId) (document.DocId, error) {

	curIdx := o.curIdx
	lastIdx := o.curIdx
	res, err := o.querys[o.curIdx].GetGE(id)
	// fmt.Println(res)
	if err != nil {
		curIdx++
	}

	for {
		curIdx = (curIdx + 1) % len(o.querys)
		cur, err := o.querys[curIdx].GetGE(id)
		//fmt.Println(cur, err)
		if err != nil {
			if res > cur {
				if o.check(res) {
					return res, nil
				}
				curIdx++
				res, err = o.querys[curIdx].Next()
				if err != nil {
					curIdx++
				}
			} else {
				return 0, helpers.NoMoreData
			}
		}

		if cur < res {
			res = cur
			lastIdx = curIdx
		}

		if (curIdx+1)%len(o.querys) == lastIdx {
			if o.check(res) {
				return res, nil
			}
			curIdx++
			res, err = o.querys[curIdx].Next()
			if err != nil {
				curIdx++
			}
		}
	}
}

func (o *OrQuery) Current() (document.DocId, error) {
	return o.querys[o.curIdx].Current()
}

func (t *OrQuery) String() string {
	return ""
}

func (o *OrQuery) check(id document.DocId) bool {
	return true
}
