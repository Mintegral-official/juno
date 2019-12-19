package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/query/check"
	"github.com/pkg/errors"
)

type AndQuery struct {
	queries  []Query
	checkers []check.Checker
	curIdx   int
}

func NewAndQuery(queries []Query, checkers []check.Checker) *AndQuery {
	if queries == nil {
		return nil
	}
	return &AndQuery{
		queries:  queries,
		checkers: checkers,
	}
}

func (aq *AndQuery) Next() (document.DocId, error) {
	lastIdx, curIdx := aq.curIdx, aq.curIdx
	target, err := aq.queries[curIdx].Next()

	if err != nil {
		return 0, helpers.NoMoreData
	}
	if curIdx == len(aq.queries)-1 {
		if !aq.check(target) {
			target, err = aq.queries[curIdx].Next()
			if err != nil {
				return 0, helpers.NoMoreData
			}
		}
		return target, nil
	}

	for {
		curIdx = (curIdx + 1) % len(aq.queries)
		cur, err := aq.queries[curIdx].GetGE(target)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in queries[%d]", int64(target), curIdx))
		}
		if cur != target {
			lastIdx = curIdx
			target = cur
		}
		if (curIdx+1)%len(aq.queries) == lastIdx {
			if aq.check(target) {
				return target, nil
			}
			curIdx = (curIdx + 1) % len(aq.queries)
			target, err = aq.queries[curIdx].Next()
			if err != nil {
				return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in queries[%d]", int64(target), curIdx))
			}
		}
	}
}

func (aq *AndQuery) GetGE(id document.DocId) (document.DocId, error) {
	curIdx, lastIdx := aq.curIdx, aq.curIdx
	res, err := aq.queries[aq.curIdx].GetGE(id)

	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in queries[%d]", int64(res), curIdx))
	}

	for {
		curIdx = (curIdx + 1) % len(aq.queries)
		cur, err := aq.queries[curIdx].GetGE(res)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in queries[%d]", int64(res), curIdx))
		}
		if cur != res {
			lastIdx = curIdx
			res = cur
		}
		if (curIdx+1)%len(aq.queries) == lastIdx {
			if aq.check(res) {
				return res, nil
			}
			curIdx = (curIdx + 1) % len(aq.queries)
			res, err = aq.queries[curIdx].Next()
			if err != nil {
				return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in queries[%d]", int64(res), curIdx))
			}
		}
	}
}

func (aq *AndQuery) Current() (document.DocId, error) {
	res, err := aq.queries[0].Current()
	if err != nil {
		return 0, err
	}

	for i := 1; i < len(aq.queries); i++ {
		tar, err := aq.queries[i].Current()
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
