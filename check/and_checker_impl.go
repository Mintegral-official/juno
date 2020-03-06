package check

import (
	"fmt"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
)

type AndChecker struct {
	c      []Checker
	aDebug *debug.Debug
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
	for i, cValue := range a.c {
		if cValue == nil {
			continue
		}
		if !cValue.Check(id) {
			if a.aDebug != nil {
				a.aDebug.AddDebugMsg(fmt.Sprintf("%d in andChecker[%d] check result: false", id, i))
			}
			return false
		}
	}
	return true
}

func (a *AndChecker) Marshal() map[string]interface{} {
	res := make(map[string]interface{}, 1)
	var tmp []map[string]interface{}
	for _, v := range a.c {
		tmp = append(tmp, v.Marshal())
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
	for _, v := range value {
		if _, ok := v["and_check"]; ok {
			var tmp = &AndChecker{}
			c = append(c, tmp.Unmarshal(idx, v, e))
		} else if _, ok := v["or_check"]; ok {
			var tmp = &OrChecker{}
			c = append(c, tmp.Unmarshal(idx, v, e))
		} else if _, ok := v["in_check"]; ok {
			var tmp = &InChecker{}
			c = append(c, tmp.Unmarshal(idx, v, e))
		} else if _, ok := v["not_check"]; ok {
			var tmp = &NotChecker{}
			c = append(c, tmp.Unmarshal(idx, v, e))
		} else if _, ok := v["check"]; ok {
			var tmp = CheckerImpl{}
			c = append(c, tmp.Unmarshal(idx, v, e))
		}
	}
	return NewAndChecker(c)
}

func (a *AndChecker) DebugInfo() *debug.Debug {
	if a.aDebug != nil {
		for _, v := range a.c {
			a.aDebug.AddDebug(v.DebugInfo())
		}
		return a.aDebug
	}
	return nil
}

func (a *AndChecker) SetDebug(level int) {
	a.aDebug = debug.NewDebug(level, "AndCheck")
	for _, v := range a.c {
		v.SetDebug(level)
	}
}
