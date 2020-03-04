package query

import (
	"github.com/Mintegral-official/juno/check"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
)

type NotAndQuery struct {
	queries  []Query
	checkers []check.Checker
	curIdx   int
	debugs   *debug.Debug
}

func NewNotAndQuery(queries []Query, checkers []check.Checker, isDebug ...int) (naq *NotAndQuery) {
	naq = &NotAndQuery{}
	if len(isDebug) == 1 && isDebug[0] == 1 {
		naq.debugs = debug.NewDebug("NotAndQuery")
	}
	if len(queries) == 0 {
		return naq
	}
	naq.checkers = checkers
	if len(queries) == 1 {
		naq.queries = queries
	} else {
		naq.queries = append(naq.queries, queries[0])
		naq.queries = append(naq.queries, NewOrQuery(queries[1:], nil))
	}

	for {
		target, err := naq.queries[0].Current()
		if err != nil {
			return naq
		}
		if len(naq.queries) == 1 {
			for !naq.check(target) {
				target, err = naq.queries[0].Current()
				naq.queries[0].Next()
			}
			return naq
		}
		for {
			tar, err := naq.queries[1].GetGE(target)
			if err != nil {
				return naq
			}
			if tar == target {
				target, err = queries[0].Current()
				naq.queries[0].Next()
				continue
			}
			break
		}
		return naq
	}
}

func (naq *NotAndQuery) Next() {
	if len(naq.queries) == 0 {
		return
	}
	target, err := naq.queries[0].Current()
	if len(naq.queries) == 1 {
		for err == nil && !naq.check(target) {
			naq.queries[0].Next()
			target, err = naq.queries[0].Current()
		}
		naq.queries[0].Next()
	}

	if err != nil {
		return
	}
	for {
		cur, err := naq.queries[1].GetGE(target)
		if cur != target {
			for err == nil && !naq.check(target) {
				naq.queries[0].Next()
				target, err = naq.queries[0].Current()
			}
			naq.queries[0].Next()
			return
		}
		naq.queries[0].Next()
		target, err = naq.queries[0].Current()
	}
}

func (naq *NotAndQuery) GetGE(id document.DocId) (document.DocId, error) {
	if len(naq.queries) == 0 {
		return 0, helpers.NoMoreData
	}
	target, err := naq.queries[0].GetGE(id)
	if len(naq.queries) == 1 {
		for err == nil && !naq.check(target) {
			naq.queries[0].Next()
			target, err = naq.queries[0].Current()
		}
		return target, err
	}

	if err != nil {
		return 0, helpers.NoMoreData
	}
	for {
		cur, err := naq.queries[1].GetGE(target)
		if cur != target {
			for err == nil && !naq.check(target) {
				naq.queries[0].Next()
				target, err = naq.queries[0].Current()
			}
			return target, nil
		}
		target, err = naq.queries[0].Current()
		naq.queries[0].Next()
	}
}

func (naq *NotAndQuery) Current() (document.DocId, error) {
	return naq.queries[0].Current()
}

func (naq *NotAndQuery) DebugInfo() *debug.Debug {
	if naq.debugs != nil {
		for _, v := range naq.queries {
			if v.DebugInfo() != nil {
				for key, value := range v.DebugInfo().Node {
					naq.debugs.Node[key] = append(naq.debugs.Node[key], value...)
				}
			}
		}
		return naq.debugs
	}
	return nil
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

func (naq *NotAndQuery) Marshal() map[string]interface{} {
	var queryInfo, checkInfo []map[string]interface{}
	res := make(map[string]interface{}, len(naq.queries))
	for _, v := range naq.queries {
		queryInfo = append(queryInfo, v.Marshal())
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

func (naq *NotAndQuery) Unmarshal(idx *index.Indexer, res map[string]interface{}, e operation.Operation) Query {
	notAnd, ok := res["not"]
	if !ok {
		return nil
	}
	notCheck, ok := res["not_and_check"]
	r := notAnd.([]map[string]interface{})
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
		return NewNotAndQuery(q, nil)
	}
	checks := notCheck.([]map[string]interface{})
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
			var tmp = &check.NotAndChecker{}
			c = append(c, tmp.Unmarshal(idx, v, e))
		}
	}
	return NewNotAndQuery(q, c)
}

func (naq *NotAndQuery) SetDebug(isDebug ...int) {
	if len(isDebug) == 1 && isDebug[0] == 1 {
		naq.debugs = debug.NewDebug("NotAndQuery")
	}
	for _, v := range naq.queries {
		v.SetDebug(1)
	}
	for _, v := range naq.checkers {
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
