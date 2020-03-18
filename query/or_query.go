package query

import (
	"container/heap"
	"github.com/MintegralTech/juno/check"
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/helpers"
	"github.com/MintegralTech/juno/index"
	"github.com/MintegralTech/juno/marshal"
)

type OrQuery struct {
	checkers []check.Checker
	h        Heap
	label    string
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
	if h.Len() == 0 {
		return nil
	}
	oq.h = *h
	oq.next()
	return oq
}

func (oq *OrQuery) SetLabel(label string) {
	if oq != nil {
		oq.label = label
	}
}

func (oq *OrQuery) Next() {
	if oq == nil || oq.h.Len() == 0 {
		return
	}
	top := oq.h.Top()
	if top != nil {
		q := top.(Query)
		q.Next()
		heap.Fix(&oq.h, 0)
	}
	oq.next()
}

func (oq *OrQuery) next() {
	if oq == nil || oq.h.Len() == 0 {
		return
	}
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
		}
	}
}

func (oq *OrQuery) getGE(id document.DocId) {
	if oq == nil || oq.h.Len() == 0 {
		return
	}
	top := oq.h.Top()
	if top != nil {
		q := top.(Query)
		_, _ = q.GetGE(id)
		heap.Fix(&oq.h, 0)
	}
}

func (oq *OrQuery) GetGE(id document.DocId) (document.DocId, error) {
	if oq == nil || oq.h.Len() == 0 {
		return 0, helpers.ElementNotfound
	}
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
	if oq == nil || oq.h.Len() == 0 {
		return 0, helpers.ElementNotfound
	}
	top := oq.h.Top()
	if top == nil {
		return 0, helpers.NoMoreData
	}
	q := top.(Query)
	return q.Current()
}

func (oq *OrQuery) DebugInfo() *debug.Debug {
	if oq != nil && oq.debugs != nil {
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
	if oq == nil || len(oq.checkers) == 0 {
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
	if oq == nil {
		return map[string]interface{}{}
	}
	var queryInfo, checkInfo []map[string]interface{}
	res := make(map[string]interface{}, len(oq.h))
	for _, v := range oq.h {
		queryInfo = append(queryInfo, v.Marshal())
	}
	if oq.label != "" {
		queryInfo = append(queryInfo, map[string]interface{}{"label": oq.label})
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

func (oq *OrQuery) MarshalV2() *marshal.MarshalInfo {
	if oq == nil {
		return nil
	}
	info := &marshal.MarshalInfo{
		Operation: "or",
		Nodes:     make([]*marshal.MarshalInfo, 0),
	}
	if oq.label != "" {
		info.Label = oq.label
	}
	for _, v := range oq.h {
		m := v.MarshalV2()
		if m != nil {
			info.Nodes = append(info.Nodes, v.MarshalV2())
		}
	}
	for _, v := range oq.checkers {
		m := v.MarshalV2()
		if m != nil {
			info.Nodes = append(info.Nodes, v.MarshalV2())
		}
	}
	return info
}

func (oq *OrQuery) UnmarshalV2(idx index.Index, marshalInfo *marshal.MarshalInfo) Query {
	if marshalInfo == nil || marshalInfo.Operation != "or" {
		return nil
	}
	var q []Query
	var c []check.Checker
	uq := &UnmarshalV2{}
	for _, v := range marshalInfo.Nodes {
		m := uq.UnmarshalV2(idx, v)
		if m != nil {
			if res, ok := m.(Query); ok {
				q = append(q, res)
			} else if res, ok := m.(check.Checker); ok {
				c = append(c, res)
			}
		}
	}
	return NewOrQuery(q, c)
}

func (oq *OrQuery) Unmarshal(idx index.Index, res map[string]interface{}) Query {
	or, ok := res["or"]
	if !ok {
		return nil
	}
	var queries []Query
	var checker []check.Checker
	uq := &Unmarshal{}
	for _, v := range or.([]map[string]interface{}) {
		if q := uq.Unmarshal(idx, v); q != nil {
			queries = append(queries, q.(Query))
		}
	}
	orCheck, ok := res["or_check"]
	if !ok {
		return NewOrQuery(queries, nil)
	}
	for _, v := range orCheck.([]map[string]interface{}) {
		if q := uq.Unmarshal(idx, v); q != nil {
			checker = append(checker, q.(check.Checker))
		}
	}
	return NewOrQuery(queries, checker)
}

func (oq *OrQuery) SetDebug(level int) {
	if oq == nil {
		return
	}
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
