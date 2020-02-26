package query

import (
	"container/heap"
	"errors"
	"fmt"
	"github.com/Mintegral-official/juno/check"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
	"strconv"
	"strings"
)

type OrQuery struct {
	checkers []check.Checker
	h        Heap
	debugs   *debug.Debugs
}

func NewOrQuery(queries []Query, checkers []check.Checker, isDebug ...int) (oq *OrQuery) {
	oq = &OrQuery{}
	if len(isDebug) == 1 && isDebug[0] == 1 {
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
		return res, err
	}
	if oq.check(res) {
		return res, nil
	}
	return res, errors.New(strconv.FormatInt(int64(res), 10) + " has been filtered out")
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
	if oq.debugs != nil {
		var msg []string
		for i, v := range oq.checkers {
			if v == nil {
				msg = append(msg, fmt.Sprintf("check[%d] is nil", i))
				continue
			}
			if c, ok := v.(*check.AndChecker); ok {
				msg = append(msg, c.DebugInfo())
			} else if c, ok := v.(*check.OrChecker); ok {
				msg = append(msg, c.DebugInfo())
			} else {
				msg = append(msg, strconv.FormatBool(v.Check(id)))
			}
		}
		oq.debugs.DebugInfo.AddDebugMsg(strconv.FormatInt(int64(id), 10) + " check: [" + strings.Join(msg, ",") + "]")
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

func (oq *OrQuery) Marshal(idx *index.Indexer) map[string]interface{} {
	var queryInfo, checkInfo []map[string]interface{}
	res := make(map[string]interface{}, len(oq.h))
	for _, v := range oq.h {
		queryInfo = append(queryInfo, v.Marshal(idx))
	}
	if len(oq.checkers) != 0 {
		for _, v := range oq.checkers {
			checkInfo = append(checkInfo, v.Marshal(idx))
		}
		res["or_check"] = checkInfo
	}
	res["or"] = queryInfo
	return res
}

func (oq *OrQuery) Unmarshal(idx *index.Indexer, res map[string]interface{}, e operation.Operation) Query {
	or, ok := res["or"]
	if !ok {
		return nil
	}
	orCheck, ok := res["or_check"]
	r := or.([]map[string]interface{})
	var q []Query
	var c []check.Checker
	for i, v := range oq.h {
		q = append(q, v.Unmarshal(idx, r[i], nil))
	}
	if !ok {
		return NewOrQuery(q, nil, 1)
	}
	checks := orCheck.([]map[string]interface{})
	for i, v := range oq.checkers {
		c = append(c, v.Unmarshal(idx, checks[i], e))
	}
	return NewOrQuery(q, c, 1)
}

func (oq *OrQuery) SetDebug(isDebug ...int) {
	if len(isDebug) == 1 && isDebug[0] == 1 {
		oq.debugs = debug.NewDebugs(debug.NewDebug("OrQuery"))
	}
	for _, v := range oq.h {
		v.SetDebug(1)
	}
	for _, v := range oq.checkers {
		switch v.(type) {
		case *check.AndChecker:
			v.(*check.AndChecker).SetDebug()
		case *check.OrChecker:
			v.(*check.OrChecker).SetDebug()
		}
	}
}

func (oq *OrQuery) UnsetDebug() {
	oq.debugs = nil
	for _, v := range oq.h {
		v.UnsetDebug()
	}
	for _, v := range oq.checkers {
		switch v.(type) {
		case *check.AndChecker:
			v.(*check.AndChecker).UnsetDebug()
		case *check.OrChecker:
			v.(*check.OrChecker).UnsetDebug()
		}
	}
}
