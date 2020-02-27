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
	naq.queries = queries
	naq.checkers = checkers
	return naq
}

func (naq *NotAndQuery) Next() (document.DocId, error) {
label:
	for {
		target, err := naq.queries[0].Current()
		if err != nil {
			return target, helpers.NoMoreData
		}
		if len(naq.queries) == 1 {
			_, _ = naq.queries[0].Next()
			if naq.check(target) {
				return target, nil
			}
		}
		for i := 1; i < len(naq.queries); i++ {
			cur, err := naq.queries[i].GetGE(target)
			if target == cur {
				_, _ = naq.queries[0].Next()
				goto label
			}
			if (target != cur || err != nil) && i == len(naq.queries)-1 {
				_, _ = naq.queries[0].Next()
				for !naq.check(target) {
					target, err = naq.queries[0].Current()
					if err != nil {
						return target, err
					}
					_, _ = naq.queries[0].Next()
				}
				return target, nil
			}
		}
		target, err = naq.queries[0].Next()
	}
}

func (naq *NotAndQuery) GetGE(id document.DocId) (document.DocId, error) {
	for {
		target, err := naq.queries[0].GetGE(id)
		if err != nil {
			return target, helpers.NoMoreData
		}
		if len(naq.queries) == 1 {
			for !naq.check(target) {
				target, err = naq.queries[0].Next()
			}
			return target, nil
		}
		for i := 1; i < len(naq.queries); i++ {
			if _, err := naq.queries[i].Current(); err != nil {
				for !naq.check(target) {
					target, err = naq.queries[0].Next()
				}
				return target, nil
			}
			cur, err := naq.queries[i].GetGE(target)
			if (target != cur || err != nil) && i == len(naq.queries)-1 {
				if naq.check(target) {
					return target, nil
				}
			}
		}
		_, _ = naq.queries[0].Next()
	}
}

func (naq *NotAndQuery) Current() (document.DocId, error) {
	res, err := naq.queries[0].Current()
	if err != nil {
		return res, err
	}
	for i := 1; i < len(naq.queries); i++ {
		tar, err := naq.queries[i].GetGE(res)
		_, _ = naq.queries[0].Next()
		if err != nil {
			continue
		}
		if tar == res {
			return res, errors.New("this target is not result")
		} else if i == len(naq.queries)-1 {
			if naq.check(res) {
				return res, nil
			}
		}
	}
	return 0, nil
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
	if naq.debugs != nil {
		var msg []string
		var flag = true
		msg = append(msg, "not and check result: true")
		for i, c := range naq.checkers {
			if c == nil {
				msg = append(msg, fmt.Sprintf("check[%d] is nil", i))
				continue
			}
			if c.Check(id) {
				flag = false
			}
			msg = append(msg, c.DebugInfo()+"\t check result: "+strconv.FormatBool(c.Check(id)))
		}
		if flag {
			naq.debugs.Node[id] = append(naq.debugs.Node[id], msg)
		} else {
			msg[0] = "or check result: false"
			naq.debugs.Node[id] = append(naq.debugs.Node[id], msg)
		}
		return flag
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

func (naq *NotAndQuery) Marshal(idx *index.Indexer) map[string]interface{} {
	var queryInfo, checkInfo []map[string]interface{}
	res := make(map[string]interface{}, len(naq.queries))
	for _, v := range naq.queries {
		queryInfo = append(queryInfo, v.Marshal(idx))
	}
	if len(naq.checkers) != 0 {
		for _, v := range naq.checkers {
			checkInfo = append(checkInfo, v.Marshal(idx))
		}
		res["not_check"] = checkInfo
	}
	res["not"] = queryInfo
	return res
}

func (naq *NotAndQuery) Unmarshal(idx *index.Indexer, res map[string]interface{}, e operation.Operation) Query {
	and, ok := res["not"]
	if !ok {
		return nil
	}
	notCheck, ok := res["not_check"]
	r := and.([]map[string]interface{})
	var q []Query
	var c []check.Checker
	for i, v := range naq.queries {
		q = append(q, v.Unmarshal(idx, r[i], nil))
	}
	if !ok {
		return NewNotAndQuery(q, nil, 1)
	}
	checks := notCheck.([]map[string]interface{})
	for i, v := range naq.checkers {
		c = append(c, v.Unmarshal(idx, checks[i], e))
	}
	return NewNotAndQuery(q, c, 1)
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
		}
	}
}
