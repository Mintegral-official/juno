package check

import (
	"fmt"
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/index"
	"github.com/MintegralTech/juno/marshal"
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

func (a *AndChecker) MarshalV2() *marshal.MarshalInfo {
	if a == nil {
		return nil
	}
	info := &marshal.MarshalInfo{
		Operation: "or",
		Nodes:     make([]*marshal.MarshalInfo, 0),
	}
	for _, v := range a.c {
		m := v.MarshalV2()
		if m != nil {
			info.Nodes = append(info.Nodes, m)
		}
	}
	return info
}

func (a *AndChecker) UnmarshalV2(idx index.Index, info *marshal.MarshalInfo) Checker {
	if info == nil {
		return nil
	}
	var c []Checker
	uq := &unmarshalV2{}
	for _, v := range info.Nodes {
		m := uq.UnmarshalV2(idx, v)
		if m != nil {
			c = append(c, m.(Checker))
		}
	}
	return NewAndChecker(c)
}

func (a *AndChecker) Unmarshal(idx index.Index, res map[string]interface{}) Checker {
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
