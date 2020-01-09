package query

import (
	"container/heap"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/query/check"
	"strconv"
)

type OrQuery struct {
	checkers []check.Checker
	h        Heap
	debugs   *debug.Debugs
}

func NewOrQuery(queries []Query, checkers []check.Checker, isDebug ...int) *OrQuery {
	oq := &OrQuery{}
	if len(isDebug) != 0 && isDebug[0] == 1 {
		oq.debugs = debug.NewDebugs(debug.NewDebug("OrQuery"))
	}
	if len(queries) == 0 {
		return oq
	}
	h := &Heap{}
	for i := 0; i < len(queries); i++ {
		if queries[i] == nil {
			continue
		}
		heap.Push(h, queries[i])
	}
	oq.h = *h
	oq.checkers = checkers
	return oq
}

func (oq *OrQuery) Next() (document.DocId, error) {
	if oq.debugs != nil {
		oq.debugs.NextNum++
	}
	for target, err := oq.Current(); err == nil; {
		oq.next()
		if oq.check(target) {
			for cur, err := oq.Current(); err == nil; {
				if cur != target {
					break
				}
				oq.next()
				cur, err = oq.Current()
			}
			return target, nil
		}
		if oq.debugs != nil {
			oq.debugs.DebugInfo.AddDebugMsg(strconv.FormatInt(int64(target), 10) + "has been filtered out")
		}
		target, err = oq.Current()
	}
	return 0, helpers.NoMoreData
}

func (oq *OrQuery) next() {
	top := oq.h.Top()
	if top != nil {
		q := top.(Query)
		_, _ = q.Next()
		heap.Fix(&oq.h, 0)
	}
}

func (oq *OrQuery) getGE(id document.DocId) {
	top := oq.h.Top()
	if top != nil {
		q := top.(Query)
		_, _ = q.GetGE(id)
		heap.Fix(&oq.h, 0)
	}
}

func (oq *OrQuery) GetGE(id document.DocId) (document.DocId, error) {
	if oq.debugs != nil {
		oq.debugs.GetNum++
	}
	target, err := oq.Current()
	for err == nil && target < id {
		oq.getGE(id)
		target, err = oq.Current()
	}
	for err == nil && !oq.check(target) {
		if oq.debugs != nil {
			oq.debugs.DebugInfo.AddDebugMsg(strconv.FormatInt(int64(target), 10) + "has been filtered out")
		}
		target, err = oq.Next()
	}
	return target, err
}

func (oq *OrQuery) Current() (document.DocId, error) {
	if oq.debugs != nil {
		oq.debugs.CurNum++
	}
	top := oq.h.Top()
	if top == nil {
		return 0, helpers.NoMoreData
	}
	q := top.(Query)
	res, err := q.Current()
	if err != nil {
		return 0, err
	}
	if oq.check(res) {
		return res, nil
	}
	if oq.debugs != nil {
		oq.debugs.DebugInfo.AddDebugMsg(strconv.FormatInt(int64(res), 10) + "has been filtered out")
	}
	return 0, nil
}

func (oq *OrQuery) DebugInfo() *debug.Debug {
	if oq.debugs != nil {
		oq.debugs.DebugInfo.AddDebugMsg("next has been called: " + strconv.Itoa(oq.debugs.NextNum))
		oq.debugs.DebugInfo.AddDebugMsg("get has been called: " + strconv.Itoa(oq.debugs.GetNum))
		oq.debugs.DebugInfo.AddDebugMsg("current has been called: " + strconv.Itoa(oq.debugs.CurNum))
		for i := 0; i < oq.h.Len(); i++ {
			oq.debugs.DebugInfo.AddDebug(oq.h[i].DebugInfo())
		}
		return oq.debugs.DebugInfo
	}
	return nil
}

func (oq *OrQuery) check(id document.DocId) bool {
	if len(oq.checkers) == 0 {
		return true
	}
	for _, v := range oq.checkers {
		if v == nil {
			continue
		}
		if v.Check(id) {
			return true
		}
	}
	return false
}
