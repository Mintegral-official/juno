package check

import (
	"fmt"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
	"strconv"
	"strings"
)

type OrChecker struct {
	c      []Checker
	aDebug *debug.Debugs
}

func NewOrChecker(c []Checker, isDebug ...int) *OrChecker {
	if c == nil {
		return nil
	}
	o := &OrChecker{}
	if len(isDebug) == 1 && isDebug[0] == 1 {
		o.aDebug = debug.NewDebugs(debug.NewDebug("OrCheck"))
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
		for i, cValue := range o.c {
			if cValue == nil {
				msg = append(msg, fmt.Sprintf("check[%d] is nil", i))
				continue
			}
			if c, ok := cValue.(*OrChecker); ok {
				msg = append(msg, c.DebugInfo())
			} else if c, ok := cValue.(*AndChecker); ok {
				msg = append(msg, c.DebugInfo())
			} else {
				msg = append(msg, strconv.FormatBool(cValue.Check(id)))
			}
		}
		o.aDebug.DebugInfo.AddDebugMsg("[" + strings.Join(msg, ",") + "]")
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

func (o *OrChecker) Marshal(idx *index.Indexer) map[string]interface{} {
	res := make(map[string]interface{}, 1)
	var tmp []map[string]interface{}
	for _, v := range o.c {
		tmp = append(tmp, v.Marshal(idx))
	}
	res["or_check"] = tmp
	return res
}

func (o *OrChecker) Unmarshal(idx *index.Indexer, res map[string]interface{}, e operation.Operation) Checker {
	v, ok := res["or_check"]
	if !ok {
		return nil
	}
	value := v.([]interface{})
	var c []Checker
	for i, v := range o.c {
		c = append(c, v.Unmarshal(idx, value[i].(map[string]interface{}), e))
	}
	return NewOrChecker(c)
}

func (o *OrChecker) DebugInfo() string {
	return o.aDebug.DebugInfo.String()
}

func (o *OrChecker) SetDebug() {
	o.aDebug = debug.NewDebugs(debug.NewDebug("OrCheck"))
	for _, v := range o.c {
		switch v.(type) {
		case *AndChecker:
			v.(*AndChecker).aDebug = debug.NewDebugs(debug.NewDebug("AndCheck"))
		case *OrChecker:
			v.(*OrChecker).aDebug = debug.NewDebugs(debug.NewDebug("OrCheck"))
		}
	}
}

func (o *OrChecker) UnsetDebug() {
	o.aDebug = nil
	for _, v := range o.c {
		switch v.(type) {
		case *AndChecker:
			v.(*AndChecker).aDebug = nil
		case *OrChecker:
			v.(*OrChecker).aDebug = nil
		}
	}
}
