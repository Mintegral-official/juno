package check

import (
	"fmt"
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/index"
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

func (o *OrChecker) Unmarshal(idx *index.Indexer, res map[string]interface{}) Checker {
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
