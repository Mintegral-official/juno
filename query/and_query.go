package query

import (
	"errors"
	"fmt"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/query/check"
)

type AndQuery struct {
	queries  []Query
	checkers []check.Checker
	curIdx   int
	debugs   *debug.Debugs
}

func NewAndQuery(queries []Query, checkers []check.Checker) *AndQuery {
	aq := &AndQuery{
		debugs: debug.NewDebugs(debug.NewDebug("AndQuery")),
	}
	if len(queries) == 0 {
		return aq
	}
	aq.queries = queries
	aq.checkers = checkers
	return aq
}

func (aq *AndQuery) Next() (document.DocId, error) {
	aq.debugs.NextNum++
	lastIdx, curIdx := aq.curIdx, aq.curIdx
	target, err := aq.queries[curIdx].Next()
	if err != nil {
		return 0, helpers.NoMoreData
	}

	for {
		curIdx = (curIdx + 1) % len(aq.queries)
		cur, err := aq.queries[curIdx].GetGE(target)
		if err != nil {
			return 0, errors.New(fmt.Sprintf("not find [%d] in queries[%d], err: %s", int64(target), curIdx, err.Error()))
		}
		if cur != target {
			lastIdx = curIdx
			target = cur
		}
		if (curIdx+1)%len(aq.queries) == lastIdx {
			if aq.check(target) {
				return target, nil
			}
			aq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("docID[%d] is filtered out", target))
			curIdx = (curIdx + 1) % len(aq.queries)
			target, err = aq.queries[curIdx].Next()
			if err != nil {
				return 0, errors.New(fmt.Sprintf("not find [%d] in queries[%d], err: %s", int64(target), curIdx, err.Error()))
			}
		}
	}
}

func (aq *AndQuery) GetGE(id document.DocId) (document.DocId, error) {
	aq.debugs.GetNum++
	curIdx, lastIdx := aq.curIdx, aq.curIdx
	res, err := aq.queries[aq.curIdx].GetGE(id)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("not find [%d] in queries[%d], err: %s", int64(res), curIdx, err.Error()))
	}

	for {
		curIdx = (curIdx + 1) % len(aq.queries)
		cur, err := aq.queries[curIdx].GetGE(res)
		if err != nil {
			return 0, errors.New(fmt.Sprintf("not find [%d] in queries[%d], err: %s", int64(res), curIdx, err.Error()))
		}
		if cur != res {
			lastIdx = curIdx
			res = cur
		}
		if (curIdx+1)%len(aq.queries) == lastIdx {
			if aq.check(res) {
				return res, nil
			}
			aq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("docID[%d] is filtered out", res))
			curIdx = (curIdx + 1) % len(aq.queries)
			res, err = aq.queries[curIdx].Next()
			if err != nil {
				return 0, errors.New(fmt.Sprintf("not find [%d] in queries[%d], err: %s", int64(res), curIdx, err.Error()))
			}
		}
	}
}

func (aq *AndQuery) Current() (document.DocId, error) {
	aq.debugs.CurNum++
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
			return 0, errors.New("no suitable num")
		}
	}
	if aq.check(res) {
		return res, nil
	}
	aq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("docID[%d] is filtered out", res))
	return 0, errors.New(fmt.Sprintf("the result [%d] is filtered out", res))
}

func (aq *AndQuery) DebugInfo() *debug.Debug {
	aq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("next has been called: %d", aq.debugs.NextNum))
	aq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("get has been called: %d", aq.debugs.GetNum))
	aq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("current has been called: %d", aq.debugs.CurNum))
	for i := 0; i < len(aq.queries); i++ {
		aq.debugs.DebugInfo.AddDebug(aq.queries[i].DebugInfo())
	}
	return aq.debugs.DebugInfo
}

func (aq *AndQuery) check(id document.DocId) bool {
	if len(aq.checkers) == 0 {
		return true
	}
	for _, c := range aq.checkers {
		if c == nil {
			continue
		}
		if !c.Check(id) {
			return false
		}
	}
	return true
}
