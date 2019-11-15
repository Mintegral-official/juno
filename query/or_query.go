package query

import (
	"container/heap"
	"fmt"
	"github.com/Mintegral-official/juno/document"
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
	heap.Init(h)
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
	lastIdx := o.curIdx
	curIdx := o.curIdx
	top := o.h.Top().(Query)
	target, err := top.Next()
	// fmt.Println(target)
	if err != nil {
		return 0, errors.Wrap(err, "no more data")
	}
	heap.Fix(&o.h, 0)
	if curIdx == lastIdx {
		if o.check(target) {
			return target, nil
		}

	}
	return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(target), curIdx))
	//for {
	//	curIdx = (curIdx + 1) % len(o.querys)
	//	top := o.h.Top().(Query)
	//	target, err = top.Next()
	//	heap.Fix(&o.h, 0)
	//	if err != nil {
	//		return 0, errors.Wrap(err, "no more data")
	//	}
	//	if curIdx == lastIdx {
	//		if o.check(target) {
	//			return target, nil
	//		}
	//		return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(target), curIdx))
	//	}
	//}
}

func (o *OrQuery) GetGE(id document.DocId) (document.DocId, error) {
	curIdx := o.curIdx
	lastIdx := o.curIdx
	res, err := o.querys[o.curIdx].GetGE(id)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(res), curIdx))
	}

	for {
		curIdx = (curIdx + 1) % len(o.querys)
		cur, err := o.querys[curIdx].GetGE(id)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(res), curIdx))
		}
		if cur != res {
			if cur <= res {
				res = cur
			}
			lastIdx = curIdx
		}
		if (curIdx+1)%len(o.querys) == lastIdx {
			if o.check(res) {
				return res, nil
			}
			curIdx++
			res, err = o.querys[curIdx].Next()
			if err != nil {
				return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(res), curIdx))
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
