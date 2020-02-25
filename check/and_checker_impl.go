package check

import (
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
)

type AndChecker struct {
	c []Checker
}

func NewAndChecker(c []Checker) *AndChecker {
	if c == nil {
		return nil
	}
	return &AndChecker{
		c: c,
	}
}

func (a *AndChecker) Check(id document.DocId) bool {
	if a == nil {
		return true
	}
	for _, cValue := range a.c {
		if cValue == nil {
			continue
		}
		if !cValue.Check(id) {
			return false
		}
	}
	return true
}

func (a *AndChecker) Marshal(idx *index.Indexer) map[string]interface{} {
	res := make(map[string]interface{}, 1)
	var tmp []map[string]interface{}
	for _, v := range a.c {
		tmp = append(tmp, v.Marshal(idx))
	}
	res["and_check"] = tmp
	return res
}

func (a *AndChecker) Unmarshal(idx *index.Indexer, res map[string]interface{}, e operation.Operation) Checker {
	v, ok := res["and_check"]
	if !ok {
		return nil
	}
	value := v.([]map[string]interface{})
	var c []Checker
	for i, v := range a.c {
		c = append(c, v.Unmarshal(idx, value[i], e))
	}
	return NewAndChecker(c)
}
