package query

import (
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
	aq.queries = queries
	aq.checkers = checkers
	return aq
}

func (aq *AndQuery) Next() (document.DocId, error) {
	lastIdx, curIdx := aq.curIdx, aq.curIdx
	target, err := aq.queries[curIdx].Next()
	if err != nil {
		return target, helpers.NoMoreData
	}
	for {
		curIdx = (curIdx + 1) % len(aq.queries)
		cur, err := aq.queries[curIdx].GetGE(target)
		if err != nil {
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
			target, err = aq.queries[curIdx].Next()
			if err != nil {
				return target, errors.New(aq.StringBuilder(256, curIdx, target, err.Error()))
			}
		}
	}
}

func (aq *AndQuery) GetGE(id document.DocId) (document.DocId, error) {
	curIdx, lastIdx := aq.curIdx, aq.curIdx
	res, err := aq.queries[aq.curIdx].GetGE(id)
	if err != nil {
		return res, errors.New(aq.StringBuilder(256, curIdx, res, err.Error()))
	}
	for {
		curIdx = (curIdx + 1) % len(aq.queries)
		cur, err := aq.queries[curIdx].GetGE(res)
		if err != nil {
			return cur, errors.New(aq.StringBuilder(256, curIdx, res, err.Error()))
		}
		if cur != res {
			lastIdx = curIdx
			res = cur
		}
		if (curIdx+1)%len(aq.queries) == lastIdx {
			if res != 0 && aq.check(res) {
				return res, nil
			}
			curIdx = (curIdx + 1) % len(aq.queries)
			res, err = aq.queries[curIdx].Next()
			if err != nil {
				return res, errors.New(aq.StringBuilder(256, curIdx, res, err.Error()))
			}
		}
	}
}

func (aq *AndQuery) Current() (document.DocId, error) {
	res, err := aq.queries[0].Current()
	if err != nil {
		return res, err
	}
	for i := 1; i < len(aq.queries); i++ {
		tar, err := aq.queries[i].GetGE(res)
		if err != nil {
			return tar, err
		}
		if tar != res {
			return res, errors.New(fmt.Sprintf("%d in queries[%d] is different with %d in queries[%d]", res, i, tar, i+1))
		}
	}
	if aq.check(res) {
		return res, nil
	}
	return res, err
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
	if aq.debugs != nil {
		var msg []string
		var flag = true
		msg = append(msg, "and check result: true")
		for i, c := range aq.checkers {
			if c == nil {
				msg = append(msg, fmt.Sprintf("check[%d] is nil", i))
				continue
			}
			if !c.Check(id) {
				flag = false
			}
			msg = append(msg, c.DebugInfo()+"\tcheck result: "+strconv.FormatBool(c.Check(id)))
		}
		if flag {
			aq.debugs.Node[id] = append(aq.debugs.Node[id], msg)
		} else {
			msg[0] = "and check result: false"
			aq.debugs.Node[id] = append(aq.debugs.Node[id], msg)
		}
		return flag
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

func (aq *AndQuery) Marshal(idx *index.Indexer) map[string]interface{} {
	var queryInfo, checkInfo []map[string]interface{}
	res := make(map[string]interface{}, len(aq.queries))
	for _, v := range aq.queries {
		queryInfo = append(queryInfo, v.Marshal(idx))
	}
	if len(aq.checkers) != 0 {
		for _, v := range aq.checkers {
			checkInfo = append(checkInfo, v.Marshal(idx))
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
	for i, v := range aq.queries {
		q = append(q, v.Unmarshal(idx, r[i], nil))
	}
	if !ok {
		return NewAndQuery(q, nil, 1)
	}
	checks := andCheck.([]map[string]interface{})
	for i, v := range aq.checkers {
		c = append(c, v.Unmarshal(idx, checks[i], e))
	}
	return NewAndQuery(q, c, 1)
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
		}

	}
}

func (aq *AndQuery) UnsetDebug() {
	aq.debugs = nil
	for _, v := range aq.queries {
		v.UnsetDebug()
	}
	for _, v := range aq.checkers {
		switch v.(type) {
		case *check.AndChecker:
			v.(*check.AndChecker).UnsetDebug()
		case *check.OrChecker:
			v.(*check.OrChecker).UnsetDebug()
		}

	}
}
