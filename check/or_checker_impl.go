package check

import (
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
)

type OrChecker struct {
	c []Checker
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
	value := v.([]map[string]interface{})
	var c []Checker
	for i, v := range o.c {
		c = append(c, v.Unmarshal(idx, value[i], e))
	}
	return NewOrChecker(c)
}
