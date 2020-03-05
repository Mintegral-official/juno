package query

import (
	"errors"
	"fmt"
	"github.com/Mintegral-official/juno/check"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
	"strings"
)

type AndQuery struct {
	queries  []Query
	checkers []check.Checker
	curIdx   int
	debugs   *debug.Debug
}

func NewAndQuery(queries []Query, checkers []check.Checker, isDebug ...int) (aq *AndQuery) {
	aq = &AndQuery{}
	if len(isDebug) == 1 && isDebug[0] == 1 {
		aq.debugs = debug.NewDebug("AndQuery")
	}
	if len(queries) == 0 {
		return aq
	}
	aq.curIdx = 0
	aq.queries = queries
	aq.checkers = checkers
	aq.next()
	return aq
	//for {
	//	target, err := aq.queries[0].Current()
	//	if err != nil {
	//		return aq
	//	}
	//	if len(aq.queries) == 1 {
	//		for !aq.check(target) {
	//			aq.queries[0].Next()
	//			target, err = aq.queries[0].Current()
	//		}
	//		return aq
	//	}
	//	for i := 1; i < len(aq.queries); i++ {
	//		tar, _ := aq.queries[i].GetGE(target)
	//		if tar != target {
	//			aq.queries[0].Next()
	//			_, _ = aq.queries[0].Current()
	//			break
	//		} else if i == len(aq.queries)-1 {
	//			return aq
	//		}
	//	}
	//}
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
		return target, errors.New(aq.StringBuilder(256, curIdx, target, err.Error()))
	}
	for {
		curIdx = (curIdx + 1) % len(aq.queries)
		cur, err := aq.queries[curIdx].GetGE(target)
		if err != nil {
			aq.curIdx = curIdx
			return cur, errors.New(aq.StringBuilder(256, curIdx, target, err.Error()))
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
				return target, errors.New(aq.StringBuilder(256, curIdx, target, err.Error()))
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
			if v.DebugInfo() != nil {
				for key, value := range v.DebugInfo().Node {
					aq.debugs.Node[key] = append(aq.debugs.Node[key], value...)
				}
			}
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

func (aq *AndQuery) StringBuilder(cap int, value ...interface{}) string {
	var b strings.Builder
	b.Grow(cap)
	_, _ = fmt.Fprintf(&b, "queries[%d] ", value[0])
	_, _ = fmt.Fprintf(&b, "not found:[%d], ", value[1])
	_, _ = fmt.Fprintf(&b, "reason:[%s]", value[2])
	return b.String()
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

func (aq *AndQuery) SetDebug(isDebug ...int) {
	if len(isDebug) == 1 && isDebug[0] == 1 {
		aq.debugs = debug.NewDebug("AndQuery")
	}
	for _, v := range aq.queries {
		v.SetDebug(1)
	}
	for _, v := range aq.checkers {
		switch v.(type) {
		case *check.AndChecker:
			v.(*check.AndChecker).SetDebug()
		case *check.OrChecker:
			v.(*check.OrChecker).SetDebug()
		case *check.NotAndChecker:
			v.(*check.NotAndChecker).SetDebug()
		}
	}
}
