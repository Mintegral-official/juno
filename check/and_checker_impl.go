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

type AndChecker struct {
	c      []Checker
	aDebug *debug.Debugs
}

func NewAndChecker(c []Checker, isDebug ...int) *AndChecker {
	if c == nil {
		return nil
	}
	a := &AndChecker{}
	if len(isDebug) == 1 && isDebug[0] == 1 {
		a.aDebug = debug.NewDebugs(debug.NewDebug("AndCheck"))
	}
	a.c = c
	return a
}

func (a *AndChecker) Check(id document.DocId) bool {
	if a == nil {
		return true
	}
	if a.aDebug != nil {
		if a.aDebug != nil {
			var msg []string
			for i, cValue := range a.c {
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
			a.aDebug.DebugInfo.AddDebugMsg("[" + strings.Join(msg, ",") + "]")
		}
	}
	for _, cValue := range a.c {
		if cValue == nil {
			continue
		}
		if !cValue.Check(id) {
			return false
		}
	}
	return true
}

func (a *AndChecker) Marshal(idx *index.Indexer) map[string]interface{} {
	res := make(map[string]interface{}, 1)
	var tmp []map[string]interface{}
	for _, v := range a.c {
		tmp = append(tmp, v.Marshal(idx))
	}
	res["and_check"] = tmp
	return res
}

func (a *AndChecker) Unmarshal(idx *index.Indexer, res map[string]interface{}, e operation.Operation) Checker {
	v, ok := res["and_check"]
	if !ok {
		return nil
	}
	value := v.([]interface{})
	var c []Checker
	for i, v := range a.c {
		c = append(c, v.Unmarshal(idx, value[i].(map[string]interface{}), e))
	}
	return NewAndChecker(c, 1)
}

func (a *AndChecker) DebugInfo() string {
	if a.aDebug != nil {
		return a.aDebug.DebugInfo.String()
	}
	return ""
}

func (a *AndChecker) SetDebug() {
	a.aDebug = debug.NewDebugs(debug.NewDebug("AndCheck"))
	for _, v := range a.c {
		switch v.(type) {
		case *AndChecker:
			v.(*AndChecker).aDebug = debug.NewDebugs(debug.NewDebug("AndCheck"))
		case *OrChecker:
			v.(*OrChecker).aDebug = debug.NewDebugs(debug.NewDebug("OrCheck"))
		}
	}
}

func (a *AndChecker) UnsetDebug() {
	a.aDebug = nil
	for _, v := range a.c {
		switch v.(type) {
		case *AndChecker:
			v.(*AndChecker).aDebug = nil
		case *OrChecker:
			v.(*OrChecker).aDebug = nil
		}
	}
}
