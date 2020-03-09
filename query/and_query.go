package query

import (
	"fmt"
	"github.com/MintegralTech/juno/check"
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/index"
	"github.com/MintegralTech/juno/operation"
)

type AndQuery struct {
	queries  []Query
	checkers []check.Checker
	curIdx   int
	debugs   *debug.Debug
}

func NewAndQuery(queries []Query, checkers []check.Checker) (aq *AndQuery) {
	if len(queries) == 0 {
		return nil
	}
	aq = &AndQuery{
		curIdx:   0,
		queries:  queries,
		checkers: checkers,
	}
	aq.next()
	return aq
}

func (aq *AndQuery) Next() {
	aq.queries[aq.curIdx].Next()
	aq.next()
}

func (aq *AndQuery) next() {
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
				aq.curIdx = curIdx
				return
			}
		}
	}
}

func (aq *AndQuery) GetGE(id document.DocId) (document.DocId, error) {
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
			if err != nil {
				return target, err
			}
		}
	}
}

func (aq *AndQuery) Current() (document.DocId, error) {
	return aq.queries[aq.curIdx].Current()
}

func (aq *AndQuery) DebugInfo() *debug.Debug {
	if aq.debugs != nil {
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

func (aq *AndQuery) Marshal() map[string]interface{} {
	var queryInfo, checkInfo []map[string]interface{}
	res := make(map[string]interface{}, len(aq.queries))
	for _, v := range aq.queries {
		if m := v.Marshal(); m != nil {
			queryInfo = append(queryInfo, m)
		}
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

func (aq *AndQuery) Unmarshal(idx *index.Indexer, res map[string]interface{}, e operation.Operation) Query {
	and, ok := res["and"]
	if !ok {
		return nil
	}
	andCheck, ok := res["and_check"]
	r := and.([]map[string]interface{})
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
	if !ok {
		return NewAndQuery(q, nil)
	}
	checks := andCheck.([]map[string]interface{})
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
	return NewAndQuery(q, c)
}

func (aq *AndQuery) SetDebug(level int) {
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
