package check

import (
	"fmt"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
	"strconv"
)

type OrChecker struct {
	c      []Checker
	aDebug *debug.Debug
}

func NewOrChecker(c []Checker, isDebug ...int) *OrChecker {
	if c == nil {
		return nil
	}
	o := &OrChecker{}
	if len(isDebug) == 1 && isDebug[0] == 1 {
		o.aDebug = debug.NewDebug("OrCheck")
	}
	o.c = c
	return o
}

func (o *OrChecker) Check(id document.DocId) bool {
	if o == nil {
		return true
	}
	if o.aDebug != nil {
		var msg []string
		var flag = false
		msg = append(msg, "or checker: false")
		for i, c := range o.c {
			if c == nil {
				msg = append(msg, fmt.Sprintf("check[%d] is nil", i))
				continue
			}
			if c.Check(id) {
				flag = true
			}
			msg = append(msg, c.DebugInfo()+"\tis checked: "+strconv.FormatBool(c.Check(id)))
		}
		if !flag {
			o.aDebug.Node[id] = append(o.aDebug.Node[id], msg)
		} else {
			msg[0] = "and check result: true"
			o.aDebug.Node[id] = append(o.aDebug.Node[id], msg)
		}
		return flag
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
	for i, v := range o.c {
		c = append(c, v.Unmarshal(idx, value[i], e))
	}
	return NewOrChecker(c)
}

func (o *OrChecker) DebugInfo() string {
	return o.aDebug.String()
}

func (o *OrChecker) SetDebug() {
	o.aDebug = debug.NewDebug("OrCheck")
	for _, v := range o.c {
		switch v.(type) {
		case *AndChecker:
			v.(*AndChecker).aDebug = debug.NewDebug("AndCheck")
		case *OrChecker:
			v.(*OrChecker).aDebug = debug.NewDebug("OrCheck")
		}
	}
}
