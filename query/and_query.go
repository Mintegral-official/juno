package query

import (
	"encoding/json"
	"fmt"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/query/check"
	"github.com/pkg/errors"
)

type AndQuery struct {
	queries  []Query
	checkers []check.Checker
	curIdx   int
	aDebug   *debug.Debug
}

func NewAndQuery(queries []Query, checkers []check.Checker) *AndQuery {
	aq := &AndQuery{}
	if len(queries) == 0 {
		return aq
	}
	aq.queries = queries
	aq.checkers = checkers
	aq.aDebug = &debug.Debug{
		Name: "NewAndQuery",
		Msg:  []string{},
	}
	return aq
}

func (aq *AndQuery) Next() (document.DocId, error) {
	lastIdx, curIdx := aq.curIdx, aq.curIdx
	target, err := aq.queries[curIdx].Next()
	if err != nil {
		return 0, helpers.NoMoreData
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
			aq.aDebug.AddDebug(fmt.Sprintf("docID[%d] is not suitable", target))
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
			aq.aDebug.AddDebug(fmt.Sprintf("docID[%d] is not suitable", res))
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
		tar, err := aq.queries[i].GetGE(res)
		if err != nil {
			return 0, err
		}
		if tar != res {
			if i == len(aq.queries)-1 {
				return 0, errors.New("no suitable num")
			}
			continue
		}
	}
	if aq.check(res) {
		return res, nil
	}
	aq.aDebug.AddDebug(fmt.Sprintf("docID[%d] is not suitable", res))
	return 0, errors.New(fmt.Sprintf("the result [%d] is not suitable", res))
}

func (aq *AndQuery) String() string {
	for i := 0; i < len(aq.queries); i++ {
		aq.aDebug.AddDebug(aq.queries[i].String())
	}

	if res, err := json.Marshal(aq.aDebug); err == nil {
		return string(res)
	} else {
		return err.Error()
	}
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
