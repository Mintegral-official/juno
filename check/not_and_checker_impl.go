package check

import (
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
)

type NotAndChecker struct {
	c      []Checker
	aDebug *debug.Debug
}

func NewNotAndChecker(c []Checker, isDebug ...int) *NotAndChecker {
	if c == nil {
		return nil
	}
	na := &NotAndChecker{}
	if len(isDebug) == 1 && isDebug[0] == 1 {
		na.aDebug = debug.NewDebug("OrCheck")
	}
	na.c = c
	return na
}

func (na *NotAndChecker) Check(id document.DocId) bool {
	if na == nil {
		return true
	}
	for i := range na.c {
		if i == 0 {
			if !na.c[i].Check(id) {
				return false
			}
		} else {
			if na.c[i].Check(id) {
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

func (na *NotAndChecker) DebugInfo() string {
	return na.aDebug.String()
}

func (na *NotAndChecker) SetDebug() {
	na.aDebug = debug.NewDebug("OrCheck")
	for _, v := range na.c {
		switch v.(type) {
		case *AndChecker:
			v.(*AndChecker).aDebug = debug.NewDebug("AndCheck")
		case *OrChecker:
			v.(*OrChecker).aDebug = debug.NewDebug("OrCheck")
		}
	}
}
