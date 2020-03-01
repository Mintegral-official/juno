package check

import (
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
	"strconv"
)

type InChecker struct {
	si       datastruct.Iterator
	value    interface{}
	e        operation.Operation
	transfer bool
}

func NewInChecker(si datastruct.Iterator, value interface{}, e operation.Operation, transfer bool) *InChecker {
	return &InChecker{
		si:       si,
		value:    value,
		e:        e,
		transfer: transfer,
	}
}

func (i *InChecker) DebugInfo() string {
	tmp := false
	if i.e != nil {
		tmp = true
	}
	return "FieldName: " + i.si.(*datastruct.SkipListIterator).FieldName + "\t" +
		"value: " + fmt.Sprintf("%v", i.value) + "\t" +
		"OP: in" + "\t" +
		"defined operation: " + strconv.FormatBool(tmp) + "\t" +
		"transfer: " + strconv.FormatBool(i.transfer)

}

func (i *InChecker) Check(id document.DocId) bool {
	if i == nil {
		return true
	}
	element := i.si.GetGE(id)
	if element == nil {
		return false
	}
	key, v := element.Key(), element.Value()
	if key != id || v == nil {
		return false
	}
	if i.e == nil {
		if i.transfer {
			o := operation.Operations{FieldValue: i.value}
			return o.In(v)
		}
		o := operation.Operations{FieldValue: v}
		return o.In(i.value)
	}
	if i.transfer {
		i.e.SetValue(i.value)
		return i.e.In(v)
	}
	i.e.SetValue(v)
	return i.e.In(i.value)
}

func (i *InChecker) Marshal() map[string]interface{} {
	res := make(map[string]interface{}, 1)
	var tmp []interface{}
	tmp = append(tmp, i.si.(*datastruct.SkipListIterator).FieldName)
	tmp = append(tmp, i.value)
	if i.e != nil {
		tmp = append(tmp, 1)
	} else {
		tmp = append(tmp, 0)
	}
	tmp = append(tmp, i.transfer)
	res["in_check"] = tmp
	return res
}

func (i *InChecker) Unmarshal(idx *index.Indexer, res map[string]interface{}, e operation.Operation) Checker {
	v, ok := res["in_check"]
	if !ok {
		return nil
	}
	value := v.([]interface{})
	if value[2] == 1 {
		return NewInChecker(idx.GetStorageIndex().Iterator(value[0].(string)), value[1], e, value[3].(bool))
	}
	return NewInChecker(idx.GetStorageIndex().Iterator(value[0].(string)), value[1], nil, value[3].(bool))

}
