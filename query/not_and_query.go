package query

import (
	"errors"
	"fmt"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/query/check"
)

type NotAndQuery struct {
	queries  []Query
	checkers []check.Checker
	curIdx   int
	debugs   *debug.Debugs
}

func NewNotAndQuery(queries []Query, checkers []check.Checker) *NotAndQuery {
	naq := &NotAndQuery{
		debugs: debug.NewDebugs(debug.NewDebug("NotAndQuery")),
	}
	if len(queries) == 0 {
		return naq
	}
	naq.queries = queries
	naq.checkers = checkers
	return naq
}

func (naq *NotAndQuery) Next() (document.DocId, error) {
	naq.debugs.NextNum++
	for {
		target, err := naq.queries[0].Current()
		if err != nil {
			return 0, helpers.NoMoreData
		}
		if len(naq.queries) == 1 {
			_, _ = naq.queries[0].Next()
			if naq.check(target) {
				return target, nil
			}
			naq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("docID[%d] is filtered out", target))
		}
		for i := 1; i < len(naq.queries); i++ {
			cur, err := naq.queries[i].GetGE(target)
			if (target != cur || err != nil) && i == len(naq.queries)-1 {
				_, _ = naq.queries[0].Next()
				for !naq.check(target) {
					naq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("docID[%d] is filtered out", target))
					target, err = naq.queries[0].Current()
					if err != nil {
						return 0, err
					}
					_, _ = naq.queries[0].Next()
				}
				return target, nil
			}
		}
		target, err = naq.queries[0].Next()
	}
}

func (naq *NotAndQuery) GetGE(id document.DocId) (document.DocId, error) {
	naq.debugs.GetNum++
	for {
		target, err := naq.queries[0].GetGE(id)
		if err != nil {
			return 0, helpers.NoMoreData
		}
		if len(naq.queries) == 1 {
			for !naq.check(target) {
				naq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("docID[%d] is filtered out", target))
				target, err = naq.queries[0].Next()
			}
			return target, nil
		}
		for i := 1; i < len(naq.queries); i++ {
			if _, err := naq.queries[i].Current(); err != nil {
				for !naq.check(target) {
					naq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("docID[%d] is filtered out", target))
					target, err = naq.queries[0].Next()
				}
				return target, nil
			}
			cur, err := naq.queries[i].GetGE(target)
			if (target != cur || err != nil) && i == len(naq.queries)-1 {
				if naq.check(target) {
					return target, nil
				}
				naq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("docID[%d] is filtered out", target))
			}
		}
		_, _ = naq.queries[0].Next()
	}
}

func (naq *NotAndQuery) Current() (document.DocId, error) {
	naq.debugs.CurNum++
	res, err := naq.queries[0].Current()
	if err != nil {
		return 0, err
	}
	for i := 1; i < len(naq.queries); i++ {
		tar, err := naq.queries[i].GetGE(res)
		_, _ = naq.queries[0].Next()
		if err != nil {
			continue
		}
		if tar == res {
			return 0, errors.New("this target is not result")
		} else if i == len(naq.queries)-1 {
			if naq.check(res) {
				return res, nil
			}
			naq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("docID[%d] is filtered out", res))
		}
	}
	naq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("docID[%d] is filtered out", res))
	return 0, errors.New("current data is filtered out")
}

func (naq *NotAndQuery) DebugInfo() *debug.Debug {
	naq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("next has been called: %d", naq.debugs.NextNum))
	naq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("get has been called: %d", naq.debugs.GetNum))
	naq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("current has been called: %d", naq.debugs.CurNum))
	for i := 0; i < len(naq.queries); i++ {
		naq.debugs.DebugInfo.AddDebug(naq.queries[i].DebugInfo())
	}
	return naq.debugs.DebugInfo
}

func (naq *NotAndQuery) check(id document.DocId) bool {
	if len(naq.checkers) == 0 {
		return true
	}
	for _, v := range naq.checkers {
		if v == nil {
			continue
		}
		if v.Check(id) {
			return false
		}
	}
	return true
}
