package query

import (
	"github.com/MintegralTech/juno/check"
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/helpers"
	"github.com/MintegralTech/juno/index"
	"github.com/MintegralTech/juno/marshal"
)

type NotAndQuery struct {
	q        Query
	subQuery Query
	checkers []check.Checker
	curIdx   int
	label    string
	debugs   *debug.Debug
}

func NewNotAndQuery(queries []Query, checkers []check.Checker) (naq *NotAndQuery) {
	if len(queries) == 0 || queries[0] == nil {
		return nil
	}
	naq = &NotAndQuery{
		checkers: checkers,
	}
	if len(queries) == 1 {
		naq.q = queries[0]
	} else {
		naq.q = queries[0]
		naq.subQuery = NewOrQuery(queries[1:], nil)
	}
	naq.next()
	return naq
}

func (naq *NotAndQuery) SetLabel(label string) {
	if naq != nil {
		naq.label = label
	}
}

func (naq *NotAndQuery) next() {
	if naq == nil || naq.q == nil {
		return
	}
	for target, err := naq.Current(); err == nil; {
		if naq.check(target) && naq.findSubSet(target) {
			return
		}
		naq.q.Next()
		target, err = naq.Current()
	}
}

func (naq *NotAndQuery) findSubSet(id document.DocId) bool {
	if naq == nil || naq.q == nil {
		return false
	}
	if naq.subQuery == nil {
		return true
	}
	target, err := naq.subQuery.GetGE(id)
	if err != nil {
		return true
	}
	return target != id
}

func (naq *NotAndQuery) Next() {
	if naq == nil || naq.q == nil {
		return
	}
	naq.q.Next()
	naq.next()
}

func (naq *NotAndQuery) GetGE(id document.DocId) (document.DocId, error) {
	if naq == nil || naq.q == nil {
		return 0, helpers.NoMoreData
	}
	_, _ = naq.q.GetGE(id)
	naq.next()
	return naq.Current()
}

func (naq *NotAndQuery) Current() (document.DocId, error) {
	if naq == nil || naq.q == nil {
		return 0, helpers.NoMoreData
	}
	return naq.q.Current()
}

func (naq *NotAndQuery) check(id document.DocId) bool {
	if naq == nil || len(naq.checkers) == 0 {
		return true
	}
	for i, v := range naq.checkers {
		if v == nil {
			continue
		}
		if i == 0 && !v.Check(id) {
			return false
		}
		if i != 0 && v.Check(id) {
			return false
		}
	}
	return true
}

func (naq *NotAndQuery) DebugInfo() *debug.Debug {
	if naq != nil && naq.debugs != nil {
		naq.debugs.AddDebug(naq.q.DebugInfo(), naq.subQuery.DebugInfo())
		for _, v := range naq.checkers {
			naq.debugs.AddDebug(v.DebugInfo())
		}
		return naq.debugs
	}
	return nil
}

func (naq *NotAndQuery) Marshal() map[string]interface{} {
	if naq == nil {
		return map[string]interface{}{}
	}
	var queryInfo, checkInfo []map[string]interface{}
	res := make(map[string]interface{}, 2)
	queryInfo = append(queryInfo, naq.q.Marshal())
	queryInfo = append(queryInfo, naq.subQuery.Marshal())
	if naq.label != "" {
		queryInfo = append(queryInfo, map[string]interface{}{"label": naq.label})
	}
	if len(naq.checkers) != 0 {
		for _, v := range naq.checkers {
			checkInfo = append(checkInfo, v.Marshal())
		}
		res["not_and_check"] = checkInfo
	}
	res["not"] = queryInfo
	return res
}

func (naq *NotAndQuery) MarshalV2() *marshal.MarshalInfo {
	if naq == nil {
		return nil
	}
	info := &marshal.MarshalInfo{
		Operation: "not",
		Nodes:     make([]*marshal.MarshalInfo, 0),
	}
	if naq.label != "" {
		info.Label = naq.label
	}
	info.Nodes = append(info.Nodes, naq.q.MarshalV2())
	info.Nodes = append(info.Nodes, naq.subQuery.MarshalV2())

	for _, v := range naq.checkers {
		m := v.MarshalV2()
		if m != nil {
			info.Nodes = append(info.Nodes, v.MarshalV2())
		}
	}
	return info
}

func (naq *NotAndQuery) UnmarshalV2(idx index.Index, marshalInfo *marshal.MarshalInfo) Query {
	if marshalInfo == nil || marshalInfo.Operation != "not" {
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
	return NewNotAndQuery(q, c)
}

func (naq *NotAndQuery) Unmarshal(idx index.Index, res map[string]interface{}) Query {
	notAnd, ok := res["not"]
	if !ok {
		return nil
	}
	notCheck, ok := res["not_and_check"]
	var queries []Query
	var checker []check.Checker
	uq := &Unmarshal{}
	for _, v := range notAnd.([]map[string]interface{}) {
		if q := uq.Unmarshal(idx, v); q != nil {
			queries = append(queries, q.(Query))
		}
	}
	if !ok {
		return NewNotAndQuery(queries, nil)
	}
	for _, v := range notCheck.([]map[string]interface{}) {
		if c := uq.Unmarshal(idx, v); c != nil {
			checker = append(checker, c.(check.Checker))
		}
	}
	return NewNotAndQuery(queries, checker)
}

func (naq *NotAndQuery) SetDebug(level int) {
	if naq == nil {
		return
	}
	if naq.debugs == nil {
		naq.debugs = debug.NewDebug(level, "NotAndQuery")
	}
	naq.q.SetDebug(level)
	naq.subQuery.SetDebug(level)
	for _, v := range naq.checkers {
		v.SetDebug(level)
	}
}
