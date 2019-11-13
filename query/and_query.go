package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/pkg/errors"
)

type AndQuery struct {
	querys   []Query
	checkers []Checker
	curIdx   int
}

func NewAndQuery(querys []Query, checkers []Checker) *AndQuery {
	if querys == nil {
		return nil
	}
	return &AndQuery{
		querys:   querys,
		checkers: checkers,
	}
}

func (a *AndQuery) Next() (document.DocId, error) {
	lastIdx := a.curIdx
	curIdx := a.curIdx
	target, err := a.querys[curIdx].Next()
	if err != nil {
		return 0, errors.Wrap(err, "no more data")
	}
	for {
		curIdx = (curIdx + 1) % len(a.querys)
		cur, err := a.querys[curIdx].GetGE(target)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(target), curIdx))
		}
		if cur != target {
			lastIdx = curIdx
			target = cur
		}
		if (curIdx+1)%len(a.querys) == lastIdx {
			if a.check(target) {
				return target, nil
			}
			curIdx++
			target, err = a.querys[curIdx].Next()
			if err != nil {
				return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(target), curIdx))
			}
		}
	}
}

func (a *AndQuery) GetGE(id document.DocId) (document.DocId, error) {
	curIdx := a.curIdx
	res, err := a.querys[a.curIdx].GetGE(id)
	// fmt.Println(err)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(res), curIdx))
	}
	curIdx++

	//TODO 这块逻辑有问题，这样计算的是并集而不是交集
	for curIdx < len(a.querys) {
		cur, err := a.querys[a.curIdx].GetGE(res)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(cur), curIdx))
		}

		for !a.check(cur) {
			cur, err = a.querys[a.curIdx].Next()
			return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(cur), curIdx))
		}

		if cur != res {
			return 0, helpers.DocIdNotFound
		}
		curIdx++
	}
	return res, nil
}

func (a *AndQuery) String() string {
	return ""
}

func (a *AndQuery) check(id document.DocId) bool {
	if a.checkers == nil {
		return true
	}
	for _, c := range a.checkers {
		if !c.Check(id) {
			return false
		}
	}
	return true
}
