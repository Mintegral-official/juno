package query

import (
	"errors"
	"fmt"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/query/check"
	"strconv"
	"strings"
)

type AndQuery struct {
	queries  []Query
	checkers []check.Checker
	curIdx   int
	debugs   *debug.Debugs
}

func NewAndQuery(queries []Query, checkers []check.Checker, isDebug ...int) *AndQuery {
	aq := &AndQuery{}
	if len(isDebug) == 1 && isDebug[0] == 1 {
		aq.debugs = debug.NewDebugs(debug.NewDebug("AndQuery"))
	}
	if len(queries) == 0 {
		return aq
	}
	aq.queries = queries
	aq.checkers = checkers
	return aq
}

func (aq *AndQuery) Next() (document.DocId, error) {
	if aq.debugs != nil {
		aq.debugs.NextNum++
	}
	lastIdx, curIdx := aq.curIdx, aq.curIdx
	target, err := aq.queries[curIdx].Next()
	if err != nil {
		return 0, helpers.NoMoreData
	}

	for {
		curIdx = (curIdx + 1) % len(aq.queries)
		cur, err := aq.queries[curIdx].GetGE(target)
		if err != nil {
			return 0, errors.New(aq.StringBuilder(256, curIdx, target, err.Error()))
		}
		if cur != target {
			lastIdx = curIdx
			target = cur
		}
		if (curIdx+1)%len(aq.queries) == lastIdx {
			if aq.check(target) {
				return target, nil
			}
			if aq.debugs != nil {
				aq.debugs.DebugInfo.AddDebugMsg(strconv.FormatInt(int64(target), 10) + "has been filtered out")
			}
			curIdx = (curIdx + 1) % len(aq.queries)
			target, err = aq.queries[curIdx].Next()
			if err != nil {
				return 0, errors.New(aq.StringBuilder(256, curIdx, target, err.Error()))
			}
		}
	}
}

func (aq *AndQuery) GetGE(id document.DocId) (document.DocId, error) {
	if aq.debugs != nil {
		aq.debugs.GetNum++
	}
	curIdx, lastIdx := aq.curIdx, aq.curIdx
	res, err := aq.queries[aq.curIdx].GetGE(id)
	if err != nil {
		return 0, errors.New(aq.StringBuilder(256, curIdx, res, err.Error()))
	}

	for {
		curIdx = (curIdx + 1) % len(aq.queries)
		cur, err := aq.queries[curIdx].GetGE(res)
		if err != nil {
			return 0, errors.New(aq.StringBuilder(256, curIdx, res, err.Error()))
		}
		if cur != res {
			lastIdx = curIdx
			res = cur
		}
		if (curIdx+1)%len(aq.queries) == lastIdx {
			if aq.check(res) {
				return res, nil
			}
			if aq.debugs != nil {
				aq.debugs.DebugInfo.AddDebugMsg(strconv.FormatInt(int64(res), 10) + "has been filtered out")
			}
			curIdx = (curIdx + 1) % len(aq.queries)
			res, err = aq.queries[curIdx].Next()
			if err != nil {
				return 0, errors.New(aq.StringBuilder(256, curIdx, res, err.Error()))
			}
		}
	}
}

func (aq *AndQuery) Current() (document.DocId, error) {
	if aq.debugs != nil {
		aq.debugs.CurNum++
	}
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
			return 0, nil
		}
	}
	if aq.check(res) {
		return res, nil
	}
	if aq.debugs != nil {
		aq.debugs.DebugInfo.AddDebugMsg(aq.StringBuilder(128, res))
	}
	return 0, nil
}

func (aq *AndQuery) DebugInfo() *debug.Debug {
	if aq.debugs != nil {
		aq.debugs.DebugInfo.AddDebugMsg("next has been called: " + strconv.Itoa(aq.debugs.NextNum))
		aq.debugs.DebugInfo.AddDebugMsg("get has been called: " + strconv.Itoa(aq.debugs.GetNum))
		aq.debugs.DebugInfo.AddDebugMsg("current has been called: " + strconv.Itoa(aq.debugs.CurNum))
		for i := 0; i < len(aq.queries); i++ {
			aq.debugs.DebugInfo.AddDebug(aq.queries[i].DebugInfo())
		}
		return aq.debugs.DebugInfo
	}
	return nil
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

func (aq *AndQuery) StringBuilder(cap int, value ...interface{}) string {
	var b strings.Builder
	b.Grow(cap)
	_, _ = fmt.Fprintf(&b, "queries[%d] ", value[0])
	_, _ = fmt.Fprintf(&b, "not found:[%d], ", value[1])
	_, _ = fmt.Fprintf(&b, "reason:[%s]", value[2])
	return b.String()
}
