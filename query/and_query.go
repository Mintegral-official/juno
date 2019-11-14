package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/document"
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
	if curIdx == len(a.querys) - 1 {
		return target, nil
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
	lastIdx := a.curIdx
	res, err := a.querys[a.curIdx].GetGE(id)
	// fmt.Println(err)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(res), curIdx))
	}

	for {
		curIdx = (curIdx + 1) % len(a.querys)
		cur, err := a.querys[curIdx].GetGE(res)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(res), curIdx))
		}
		if cur != res {
			lastIdx = curIdx
			res = cur
		}
		if (curIdx+1)%len(a.querys) == lastIdx {
			if a.check(res) {
				return res, nil
			}
			curIdx++
			res, err = a.querys[curIdx].Next()
			if err != nil {
				return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(res), curIdx))
			}
		}
	}
	//return res, nil
}

func (a *AndQuery) Current() (document.DocId, error) {
	return 0, nil
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
