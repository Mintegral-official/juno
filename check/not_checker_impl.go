package check

import (
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
)

type NotChecker struct {
	si       datastruct.Iterator
	value    interface{}
	e        operation.Operation
	aDebug   *debug.Debug
	transfer bool
}

func NewNotChecker(si datastruct.Iterator, value interface{}, e operation.Operation, transfer bool) *NotChecker {
	return &NotChecker{
		si:       si,
		value:    value,
		e:        e,
		transfer: transfer,
	}
}

func (nc *NotChecker) DebugInfo() *debug.Debug {
	if nc.aDebug != nil {
		return nc.aDebug
	}
	return nil
}

func (nc *NotChecker) SetDebug(level int) {
	if nc.aDebug == nil {
		nc.aDebug = debug.NewDebug(level, "not checker")
	}
}

func (nc *NotChecker) Check(id document.DocId) bool {
	if nc == nil {
		return true
	}

	element := nc.si.GetGE(id)
	if element == nil {
		if nc.aDebug != nil {
			nc.aDebug.AddDebugMsg(fmt.Sprintf("docId: %d, Field: %s, Value: %v",
				id, nc.si.(*datastruct.SkipListIterator).FieldName, nc.value))
		}
		return false
	}
	key, v := element.Key(), element.Value()
	if key != id || v == nil {
		if nc.aDebug != nil {
			nc.aDebug.AddDebugMsg(fmt.Sprintf("docID: %d, GetGE ID %d,: Field: %s, FieldValue: %v, Value: %v",
				id, key, nc.si.(*datastruct.SkipListIterator).FieldName, v, nc.value))
		}
		return false
	}
	if nc.e == nil {
		if nc.transfer {
			o := operation.Operations{FieldValue: nc.value}
			if nc.aDebug != nil {
				nc.aDebug.AddDebugMsg(fmt.Sprintf("docID: %d, field: %s, fieldValue: %v, value: %v",
					id, nc.si.(*datastruct.SkipListIterator).FieldName, v, nc.value))
			}
			return !o.In(v)
		}
		o := operation.Operations{FieldValue: v}
		if nc.aDebug != nil {
			nc.aDebug.AddDebugMsg(fmt.Sprintf("docID: %d, field: %s, fieldValue: %v, value: %v",
				id, nc.si.(*datastruct.SkipListIterator).FieldName, v, nc.value))
		}
		return !o.In(nc.value)
	}
	if nc.transfer {
		nc.e.SetValue(nc.value)
		if nc.aDebug != nil {
			nc.aDebug.AddDebugMsg(fmt.Sprintf("docID: %d, field: %s, fieldValue: %v, value: %v",
				id, nc.si.(*datastruct.SkipListIterator).FieldName, v, nc.value))
		}
		return !nc.e.In(v)
	}
	nc.e.SetValue(v)
	if nc.aDebug != nil {
		nc.aDebug.AddDebugMsg(fmt.Sprintf("docID: %d, field: %s, fieldValue: %v, value: %v",
			id, nc.si.(*datastruct.SkipListIterator).FieldName, v, nc.value))
	}
	return !nc.e.In(nc.value)
}

func (nc *NotChecker) Marshal() map[string]interface{} {
	res := make(map[string]interface{}, 1)
	var tmp []interface{}
	tmp = append(tmp, nc.si.(*datastruct.SkipListIterator).FieldName)
	tmp = append(tmp, nc.value)
	if nc.e != nil {
		tmp = append(tmp, 1)
	} else {
		tmp = append(tmp, 0)
	}
	tmp = append(tmp, nc.transfer)
	res["not_check"] = tmp
	return res
}

func (nc *NotChecker) Unmarshal(idx *index.Indexer, res map[string]interface{}, e operation.Operation) Checker {
	v, ok := res["not_check"]
	if !ok {
		return nil
	}
	value := v.([]interface{})
	if value[2] == 1 {
		return NewNotChecker(idx.GetStorageIndex().Iterator(value[0].(string)), value[1], e, value[3].(bool))
	}
	return NewNotChecker(idx.GetStorageIndex().Iterator(value[0].(string)), value[1], nil, value[3].(bool))

}
