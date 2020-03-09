package check

import (
	"fmt"
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/index"
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

func (na *NotAndChecker) Unmarshal(idx *index.Indexer, res map[string]interface{}) Checker {
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
