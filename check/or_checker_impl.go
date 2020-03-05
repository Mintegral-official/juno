package check

import (
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
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

func (o *OrChecker) Unmarshal(idx *index.Indexer, res map[string]interface{}, e operation.Operation) Checker {
	v, ok := res["or_check"]
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
	return NewOrChecker(c)
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
