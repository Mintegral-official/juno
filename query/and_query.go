package query

import (
	"fmt"
	"github.com/MintegralTech/juno/check"
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/helpers"
	"github.com/MintegralTech/juno/index"
)

type AndQuery struct {
	queries  []Query
	checkers []check.Checker
	curIdx   int
	label    string
	debugs   *debug.Debug
}

func NewAndQuery(queries []Query, checkers []check.Checker) (aq *AndQuery) {
	if len(queries) == 0 {
		return nil
	}
	aq = &AndQuery{
		curIdx:   0,
		checkers: checkers,
	}
	for i := 0; i < len(queries); i++ {
		if queries[i] != nil {
			aq.queries = append(aq.queries, queries[i])
		}
	}
	if len(aq.queries) == 0 {
		return nil
	}
	aq.next()
	return aq
}

func (aq *AndQuery) SetLabel(label string) {
	if aq != nil {
		aq.label = label
	}
}

func (aq *AndQuery) Next() {
	if aq == nil || len(aq.queries) == 0 {
		return
	}
	aq.queries[aq.curIdx].Next()
	aq.next()
}

func (aq *AndQuery) next() {
	if aq == nil || len(aq.queries) == 0 {
		return
	}
	lastIdx, curIdx := aq.curIdx, aq.curIdx
	target, err := aq.queries[curIdx].Current()
	if err != nil {
		return
	}
	for {
		curIdx = (curIdx + 1) % len(aq.queries)
		cur, err := aq.queries[curIdx].GetGE(target)
		if err != nil {
			if aq.debugs != nil {
				aq.debugs.AddDebugMsg(fmt.Sprintf("%d not found in queries[%d]", target, curIdx))
			}
			aq.curIdx = curIdx
			return
		}
		if cur != target {
			lastIdx = curIdx
			target = cur
		}
		if (curIdx+1)%len(aq.queries) == lastIdx {
			if aq.check(target) {
				return
			}
			curIdx = (curIdx + 1) % len(aq.queries)
			aq.queries[curIdx].Next()
			target, err = aq.queries[curIdx].Current()
			if err != nil {
				if aq.debugs != nil {
					aq.debugs.AddDebugMsg(fmt.Sprintf("%d not found in queries[%d]", target, curIdx))
				}
				aq.curIdx = curIdx
				return
			}
		}
	}
}

func (aq *AndQuery) GetGE(id document.DocId) (document.DocId, error) {
	if aq == nil || len(aq.queries) == 0 {
		return 0, helpers.ElementNotfound
	}
	curIdx, lastIdx := 0, 0
	target, err := aq.queries[aq.curIdx].GetGE(id)
	if err != nil {
		return target, err
	}
	for {
		curIdx = (curIdx + 1) % len(aq.queries)
		cur, err := aq.queries[curIdx].GetGE(target)
		if err != nil {
			if aq.debugs != nil {
				aq.debugs.AddDebugMsg(fmt.Sprintf("%d not found in queries[%d]", target, curIdx))
			}
			aq.curIdx = curIdx
			return cur, err
		}
		if cur != target {
			lastIdx = curIdx
			target = cur
		}
		if (curIdx+1)%len(aq.queries) == lastIdx {
			if aq.check(target) {
				return target, nil
			}
			curIdx = (curIdx + 1) % len(aq.queries)
			aq.queries[curIdx].Next()
			target, err = aq.queries[curIdx].Current()
			if aq.debugs != nil {
				aq.debugs.AddDebugMsg(fmt.Sprintf("%d not found in queries[%d]", target, curIdx))
			}
			if err != nil {
				return target, err
			}
		}
	}
}

func (aq *AndQuery) Current() (document.DocId, error) {
	if aq == nil || len(aq.queries) == 0 {
		return 0, helpers.ElementNotfound
	}
	return aq.queries[aq.curIdx].Current()
}

func (aq *AndQuery) DebugInfo() *debug.Debug {
	if aq != nil && aq.debugs != nil {
		for _, v := range aq.queries {
			aq.debugs.AddDebug(v.DebugInfo())
		}
		for _, v := range aq.checkers {
			aq.debugs.AddDebug(v.DebugInfo())
		}
		return aq.debugs
	}
	return nil
}

func (aq *AndQuery) check(id document.DocId) bool {
	if aq == nil || len(aq.checkers) == 0 {
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

func (aq *AndQuery) Marshal() map[string]interface{} {
	if aq == nil {
		return map[string]interface{}{}
	}
	var queryInfo, checkInfo []map[string]interface{}
	res := make(map[string]interface{}, len(aq.queries))
	for _, v := range aq.queries {
		if m := v.Marshal(); m != nil {
			queryInfo = append(queryInfo, m)
		}
	}
	if aq.label != "" {
		queryInfo = append(queryInfo, map[string]interface{}{"label": aq.label})
	}
	if len(aq.checkers) != 0 {
		for _, v := range aq.checkers {
			if m := v.Marshal(); m != nil {
				checkInfo = append(checkInfo, m)
			}
		}
		res["and_check"] = checkInfo
	}
	res["and"] = queryInfo
	return res
}

func (aq *AndQuery) Unmarshal(idx index.Index, res map[string]interface{}) Query {
	and, ok := res["and"]
	if !ok {
		return nil
	}
	var queries []Query
	var checkers []check.Checker
	uq := &Unmarshal{}
	for _, v := range and.([]map[string]interface{}) {
		if q := uq.Unmarshal(idx, v); q != nil {
			queries = append(queries, q.(Query))
		}
	}
	andCheck, ok := res["and_check"]
	if !ok {
		return NewAndQuery(queries, nil)
	}
	for _, v := range andCheck.([]map[string]interface{}) {
		if c := uq.Unmarshal(idx, v); c != nil {
			checkers = append(checkers, c.(check.Checker))
		}
	}
	return NewAndQuery(queries, checkers)
}

func (aq *AndQuery) SetDebug(level int) {
	if aq == nil {
		return
	}
	if aq.debugs == nil {
		aq.debugs = debug.NewDebug(level, "AndQuery")
	}
	for _, v := range aq.queries {
		v.SetDebug(1)
	}
	for _, v := range aq.checkers {
		v.SetDebug(level)
	}
}
