package check

import (
	"fmt"
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/index"
	"github.com/MintegralTech/juno/marshal"
)

type OrChecker struct {
	c      []Checker
	aDebug *debug.Debug
}

func NewOrChecker(c []Checker) *OrChecker {
	if c == nil {
		return nil
	}
	return &OrChecker{
		c: c,
	}
}

func (o *OrChecker) Check(id document.DocId) bool {
	if o == nil {
		return true
	}
	for _, cValue := range o.c {
		if cValue == nil {
			continue
		}
		if cValue.Check(id) {
			return true
		}
	}
	if o.aDebug != nil {
		o.aDebug.AddDebugMsg(fmt.Sprintf("%d in orChecker check result: false", id))
	}
	return false
}

func (o *OrChecker) Marshal() map[string]interface{} {
	res := make(map[string]interface{}, 1)
	var tmp []map[string]interface{}
	for _, v := range o.c {
		tmp = append(tmp, v.Marshal())
	}
	res["or_check"] = tmp
	return res
}

func (o *OrChecker) MarshalV2() *marshal.MarshalInfo {
	if o == nil {
		return nil
	}
	info := &marshal.MarshalInfo{
		Operation: "or_check",
		Nodes:     make([]*marshal.MarshalInfo, 0),
	}
	for _, v := range o.c {
		m := v.MarshalV2()
		if m != nil {
			info.Nodes = append(info.Nodes, m)
		}
	}
	return info
}

func (o *OrChecker) UnmarshalV2(idx index.Index, info *marshal.MarshalInfo) Checker {
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
	return NewOrChecker(c)
}

func (o *OrChecker) Unmarshal(idx index.Index, res map[string]interface{}) Checker {
	value, ok := res["or_check"]
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
	return NewOrChecker(checks)
}

func (o *OrChecker) DebugInfo() *debug.Debug {
	if o.aDebug != nil {
		for _, v := range o.c {
			if v.DebugInfo() != nil {
				o.aDebug.AddDebug(v.DebugInfo())
			}
		}
		return o.aDebug
	}
	return nil
}

func (o *OrChecker) SetDebug(level int) {
	o.aDebug = debug.NewDebug(level, "OrCheck")
	for _, v := range o.c {
		v.SetDebug(level)
	}
}
