package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/query/check"
	"github.com/pkg/errors"
)

type AndQuery struct {
	querys   []Query
	checkers []check.Checker
	curIdx   int
}

func NewAndQuery(querys []Query, checkers []check.Checker) *AndQuery {
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
			curIdx = (curIdx + 1) % len(a.querys)
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
}

func (a *AndQuery) Current() (document.DocId, error) {
	res, err := a.querys[0].Current()
	if err != nil {
		return 0, err
	}

	for i := 1; i < len(a.querys); i++ {
		tar, err := a.querys[i].Current()
		if err != nil {
			return 0, err
		}
		if tar != res {
			return 0, helpers.ElementNotfound
		}
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
