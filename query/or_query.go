package query

import (
	"container/heap"
	"fmt"
	"github.com/Mintegral-official/juno/check"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
)

type OrQuery struct {
	checkers []check.Checker
	h        Heap
	debugs   *debug.Debug
	lastId   *document.DocId
}

func NewOrQuery(queries []Query, checkers []check.Checker) (oq *OrQuery) {
	if len(queries) == 0 {
		return nil
	}
	oq = &OrQuery{
		checkers: checkers,
	}
	h := &Heap{}
	for i := 0; i < len(queries); i++ {
		if queries[i] == nil {
			continue
		}
		heap.Push(h, queries[i])
	}
	oq.h = *h
	oq.next()
	return oq
}

func (oq *OrQuery) Next() {
	top := oq.h.Top()
	if top != nil {
		q := top.(Query)
		q.Next()
		heap.Fix(&oq.h, 0)
	}
	oq.next()
}

func (oq *OrQuery) next() {
	for target, err := oq.Current(); err == nil; {
		if (oq.lastId == nil || *oq.lastId != target) && oq.check(target) {
			oq.lastId = &target
			return
		}
		top := oq.h.Top()
		if top != nil {
			q := top.(Query)
			q.Next()
			heap.Fix(&oq.h, 0)
			target, err = oq.Current()
		} else {
			if oq.debugs != nil {
				oq.debugs.AddDebugMsg(fmt.Sprintf("query is nil"))
			}
		}
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
	target, err := oq.Current()
	for err == nil && target < id {
		oq.getGE(id)
		target, err = oq.Current()
	}
	for err == nil && !oq.check(target) {
		oq.Next()
		target, err = oq.Current()
	}
	return target, err
}

func (oq *OrQuery) Current() (document.DocId, error) {
	top := oq.h.Top()
	if top == nil {
		return 0, helpers.NoMoreData
	}
	q := top.(Query)
	return q.Current()
}

func (oq *OrQuery) DebugInfo() *debug.Debug {
	if oq.debugs != nil {
		for _, v := range oq.h {
			oq.debugs.AddDebug(v.DebugInfo())
		}
		for _, v := range oq.checkers {
			oq.debugs.AddDebug(v.DebugInfo())
		}
		return oq.debugs
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

func (oq *OrQuery) Marshal() map[string]interface{} {
	var queryInfo, checkInfo []map[string]interface{}
	res := make(map[string]interface{}, len(oq.h))
	for _, v := range oq.h {
		queryInfo = append(queryInfo, v.Marshal())
	}
	if len(oq.checkers) != 0 {
		for _, v := range oq.checkers {
			checkInfo = append(checkInfo, v.Marshal())
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
	r := or.([]map[string]interface{})
	var q []Query
	var c []check.Checker
	for _, v := range r {
		if _, ok := v["and"]; ok {
			var tmp = &AndQuery{}
			q = append(q, tmp.Unmarshal(idx, v, e))
		} else if _, ok := v["or"]; ok {
			var tmp = &OrQuery{}
			q = append(q, tmp.Unmarshal(idx, v, e))
		} else if _, ok := v["not"]; ok {
			var tmp = &NotAndQuery{}
			q = append(q, tmp.Unmarshal(idx, v, e))
		} else if _, ok := v["="]; ok {
			var tmp = &TermQuery{}
			q = append(q, tmp.Unmarshal(idx, v, e))
		}
	}
	orCheck, ok := res["or_check"]
	if !ok {
		return NewOrQuery(q, nil)
	}
	checks := orCheck.([]map[string]interface{})
	for _, v := range checks {
		if _, ok := v["and_check"]; ok {
			var tmp = &check.AndChecker{}
			c = append(c, tmp.Unmarshal(idx, v, e))
		} else if _, ok := v["or_check"]; ok {
			var tmp = &check.OrChecker{}
			c = append(c, tmp.Unmarshal(idx, v, e))
		} else if _, ok := v["in_check"]; ok {
			var tmp = &check.InChecker{}
			c = append(c, tmp.Unmarshal(idx, v, e))
		} else if _, ok := v["not_check"]; ok {
			var tmp = &check.NotChecker{}
			c = append(c, tmp.Unmarshal(idx, v, e))
		} else if _, ok := v["check"]; ok {
			var tmp = &check.CheckerImpl{}
			c = append(c, tmp.Unmarshal(idx, v, e))
		} else if _, ok := v["not_and_check"]; ok {
			var tmp = check.NotAndChecker{}
			c = append(c, tmp.Unmarshal(idx, v, e))
		}
	}
	return NewOrQuery(q, c)
}

func (oq *OrQuery) SetDebug(level int) {
	if oq.debugs == nil {
		oq.debugs = debug.NewDebug(level, "OrQuery")
	}
	for _, v := range oq.h {
		v.SetDebug(level)
	}
	for _, v := range oq.checkers {
		v.SetDebug(level)
	}
}
