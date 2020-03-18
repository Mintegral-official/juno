package check

import (
	"fmt"
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/index"
	"github.com/MintegralTech/juno/marshal"
)

type NotAndChecker struct {
	c      []Checker
	aDebug *debug.Debug
}

func NewNotAndChecker(c []Checker) *NotAndChecker {
	if c == nil {
		return nil
	}
	return &NotAndChecker{
		c: c,
	}
}

func (na *NotAndChecker) Check(id document.DocId) bool {
	if na == nil {
		return true
	}
	for i := range na.c {
		if i == 0 {
			if !na.c[i].Check(id) {
				if na.aDebug != nil {
					na.aDebug.AddDebugMsg(fmt.Sprintf("%d in notAndChecker[%d] check result: false", id, i))
				}
				return false
			}
		} else {
			if na.c[i].Check(id) {
				if na.aDebug != nil {
					na.aDebug.AddDebugMsg(fmt.Sprintf("%d in notAndChecker[%d] check result: false", id, i))
				}
				return false
			}
		}

	}
	return true
}

func (na *NotAndChecker) Marshal() map[string]interface{} {
	res := make(map[string]interface{}, 1)
	var tmp []map[string]interface{}
	for _, v := range na.c {
		tmp = append(tmp, v.Marshal())
	}
	res["not_and_check"] = tmp
	return res
}


func (na *NotAndChecker) MarshalV2() *marshal.MarshalInfo {
	if na == nil {
		return nil
	}
	info := &marshal.MarshalInfo{
		Operation: "not_and_check",
		Nodes:     make([]*marshal.MarshalInfo, 0),
	}
	for _, v := range na.c {
		m := v.MarshalV2()
		if m != nil {
			info.Nodes = append(info.Nodes, m)
		}
	}
	return info
}

func (na *NotAndChecker) UnmarshalV2(idx index.Index, info *marshal.MarshalInfo) Checker {
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
	return NewNotAndChecker(c)
}

func (na *NotAndChecker) Unmarshal(idx index.Index, res map[string]interface{}) Checker {
	value, ok := res["not_and_check"]
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
	return NewNotAndChecker(checks)
}

func (na *NotAndChecker) DebugInfo() *debug.Debug {
	if na.aDebug != nil {
		for _, v := range na.c {
			if v.DebugInfo() != nil {
				na.aDebug.AddDebug(v.DebugInfo())
			}
		}
		return na.aDebug
	}
	return nil
}

func (na *NotAndChecker) SetDebug(level int) {
	if na.aDebug == nil {
		na.aDebug = debug.NewDebug(level, "OrCheck")
	}
	for _, v := range na.c {
		v.SetDebug(level)
	}
}
