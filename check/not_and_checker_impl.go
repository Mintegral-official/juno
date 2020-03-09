package check

import (
	"fmt"
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/index"
	"github.com/MintegralTech/juno/operation"
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

func (na *NotAndChecker) Unmarshal(idx *index.Indexer, res map[string]interface{}, e operation.Operation) Checker {
	v, ok := res["not_and_check"]
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
		} else if _, ok := v["nor_and_check"]; ok {
			var tmp = NotAndChecker{}
			c = append(c, tmp.Unmarshal(idx, v, e))
		}
	}
	return NewOrChecker(c)
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
