package query

import (
	"container/heap"
	"errors"
	"fmt"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/query/check"
)

type OrQuery struct {
	checkers []check.Checker
	h        Heap
	debugs   *debug.Debugs
}

func NewOrQuery(queries []Query, checkers []check.Checker) *OrQuery {
	oq := &OrQuery{
		debugs: debug.NewDebugs(debug.NewDebug("OrQuery")),
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
	oq.debugs.NextNum++
	for target, err := oq.Current(); err == nil; {
		oq.next()
		if oq.check(target) {
			return target, nil
		}
		oq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("docID[%d] is filtered out", target))
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
	oq.debugs.GetNum++
	target, err := oq.Current()
	for err == nil && target < id {
		oq.getGE(id)
		target, err = oq.Current()
	}
	for err == nil && !oq.check(target) {
		oq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("docID[%d] is filtered out", target))
		target, err = oq.Next()
	}
	return target, err
}

func (oq *OrQuery) Current() (document.DocId, error) {
	oq.debugs.CurNum++
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
	oq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("docID[%d] is filtered out", res))
	return 0, errors.New(fmt.Sprintf("the result [%d] is filtered out", res))

}

func (oq *OrQuery) DebugInfo() *debug.Debug {
	oq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("next has been called: %d", oq.debugs.NextNum))
	oq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("get has been called: %d", oq.debugs.GetNum))
	oq.debugs.DebugInfo.AddDebugMsg(fmt.Sprintf("current has been called: %d", oq.debugs.CurNum))
	for i := 0; i < oq.h.Len(); i++ {
		oq.debugs.DebugInfo.AddDebug(oq.h[i].DebugInfo())
	}
	return oq.debugs.DebugInfo
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
