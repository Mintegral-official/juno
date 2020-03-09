package check

import (
	"fmt"
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/index"
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

func (a *AndChecker) Unmarshal(idx *index.Indexer, res map[string]interface{}) Checker {
	value, ok := res["and_check"]
	if !ok {
		return nil
	}
	var checks []Checker
	uq := &unmarshal{}
	for _, v := range value.([]map[string]interface{}) {
		if c := uq.Unmarshal(idx, v); c != nil {
			checks = append(checks, c)
		}
	}
	return NewAndChecker(checks)
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
