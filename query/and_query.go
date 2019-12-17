package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/query/check"
	"github.com/pkg/errors"
)

type AndQuery struct {
	querySlice []Query
	checkers   []check.Checker
	curIdx     int
}

func NewAndQuery(querys []Query, checkers []check.Checker) *AndQuery {
	if querys == nil {
		return nil
	}
	return &AndQuery{
		querySlice: querys,
		checkers:   checkers,
	}
}

func (aq *AndQuery) Next() (document.DocId, error) {
	lastIdx, curIdx := aq.curIdx, aq.curIdx
	target, err := aq.querySlice[curIdx].Next()

	if err != nil {
		return 0, errors.Wrap(err, "no more data")
	}
	if curIdx == len(aq.querySlice)-1 {
		if !aq.check(target) {
			target, err = aq.querySlice[curIdx].Next()
			if err != nil {
				return 0, errors.Wrap(err, "no more data")
			}
		}
		return target, nil
	}

	for {
		curIdx = (curIdx + 1) % len(aq.querySlice)
		cur, err := aq.querySlice[curIdx].GetGE(target)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(target), curIdx))
		}
		if cur != target {
			lastIdx = curIdx
			target = cur
		}
		if (curIdx+1)%len(aq.querySlice) == lastIdx {
			if aq.check(target) {
				return target, nil
			}
			curIdx = (curIdx + 1) % len(aq.querySlice)
			target, err = aq.querySlice[curIdx].Next()
			if err != nil {
				return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(target), curIdx))
			}
		}
	}
}

func (aq *AndQuery) GetGE(id document.DocId) (document.DocId, error) {
	curIdx, lastIdx := aq.curIdx, aq.curIdx
	res, err := aq.querySlice[aq.curIdx].GetGE(id)

	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(res), curIdx))
	}

	for {
		curIdx = (curIdx + 1) % len(aq.querySlice)
		cur, err := aq.querySlice[curIdx].GetGE(res)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(res), curIdx))
		}
		if cur != res {
			lastIdx = curIdx
			res = cur
		}
		if (curIdx+1)%len(aq.querySlice) == lastIdx {
			if aq.check(res) {
				return res, nil
			}
			curIdx++
			res, err = aq.querySlice[curIdx].Next()
			if err != nil {
				return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(res), curIdx))
			}
		}
	}
}

func (aq *AndQuery) Current() (document.DocId, error) {
	res, err := aq.querySlice[0].Current()
	if err != nil {
		return 0, err
	}

	for i := 1; i < len(aq.querySlice); i++ {
		tar, err := aq.querySlice[i].Current()
		if err != nil {
			return 0, err
		}
		if tar != res {
			return 0, helpers.ElementNotfound
		}
	}
	return res, nil
}

func (aq *AndQuery) String() string {
	return ""
}

func (aq *AndQuery) check(id document.DocId) bool {
	if aq.checkers == nil {
		return true
	}
	for _, c := range aq.checkers {
		if !c.Check(id) {
			return false
		}
	}
	return true
}
