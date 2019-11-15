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
	curIdx   int
}

func NewOrQuery(querys []Query, checkers []Checker) *OrQuery {
	if querys == nil {
		return nil
	}
	return &OrQuery{
		querys:   querys,
		checkers: checkers,
	}
}

func (o *OrQuery) Next() (document.DocId, error) {
	lastIdx := o.curIdx
	curIdx := o.curIdx
	h := &Heap{}
	heap.Init(h)
	for i := 0; i < len(o.querys); i++ {
		heap.Push(h, o.querys[i])
	}

	q := heap.Pop(h).(Query)
	target, err := q.Next()
	if err != nil {
		target, err = q.Next()
	}

	if len(o.querys) == 1 {
		if o.check(target) {
			return target, err
		}
		return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(target), curIdx))
	}
	heap.Push(h, q)
	for {
		curIdx = (curIdx + 1) % len(o.querys)
		tmp := h.Pop().(Query)
		cur, err := tmp.Current()
		//cur, err := tmp.GetGE(target)
		for err != nil {
			cur, err = tmp.Next()
		}

		if cur < target {
			target = cur
			lastIdx = curIdx
		}

		heap.Push(h, tmp)

		if (curIdx+1)%len(o.querys) == lastIdx {

			if o.check(target) {
				return target, nil
			}
			curIdx++
		}
		//fmt.Println(curIdx)
	}
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
