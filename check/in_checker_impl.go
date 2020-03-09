package check

import (
	"fmt"
	"github.com/MintegralTech/juno/datastruct"
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/helpers"
	"github.com/MintegralTech/juno/index"
	"github.com/MintegralTech/juno/operation"
)

type InChecker struct {
	si       datastruct.Iterator
	value    interface{}
	e        operation.Operation
	aDebug   *debug.Debug
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

func (i *InChecker) DebugInfo() *debug.Debug {
	if i.aDebug != nil {
		i.aDebug.FieldName = i.si.(*datastruct.SkipListIterator).FieldName
		return i.aDebug
	}
	return nil
}

func (i *InChecker) SetDebug(level int) {
	if i.aDebug == nil {
		i.aDebug = debug.NewDebug(level, "In checker")
	}
}

func (i *InChecker) Check(id document.DocId) bool {
	if i == nil {
		return true
	}
	element := i.si.GetGE(id)
	if element == nil {
		if i.aDebug != nil {
			i.aDebug.AddDebugMsg(fmt.Sprintf("docId: %d, value: %v, reason: %v", id, i.value, helpers.ElementNotfound))
		}
		return false
	}
	key, v := element.Key(), element.Value()
	if key != id || v == nil {
		if i.aDebug != nil {
			i.aDebug.AddDebugMsg(fmt.Sprintf("docID: %d, GetGE[ID: %d, value: %v], value: %v",
				id, key, v, i.value))
		}
		return false
	}
	var f bool
	if i.e == nil {
		if i.transfer {
			o := operation.Operations{FieldValue: i.value}
			f = o.In(v)
			if i.aDebug != nil && f == false {
				i.aDebug.AddDebugMsg(fmt.Sprintf("docID: %d, GetGE[ID: %d, value: %v], value: %v",
					id, key, v, i.value))
			}
			return f
		}
		o := operation.Operations{FieldValue: v}
		f = o.In(i.value)
		if i.aDebug != nil && f == false {
			i.aDebug.AddDebugMsg(fmt.Sprintf("docID: %d, GetGE[ID: %d, value: %v], value: %v",
				id, key, v, i.value))
		}
		return f
	}
	if i.transfer {
		i.e.SetValue(i.value)
		f = i.e.In(v)
		if i.aDebug != nil && f == false {
			i.aDebug.AddDebugMsg(fmt.Sprintf("docID: %d, GetGE[ID: %d, value: %v], value: %v",
				id, key, v, i.value))
		}
		return f
	}
	i.e.SetValue(v)
	f = i.e.In(i.value)
	if i.aDebug != nil && f == false {
		i.aDebug.AddDebugMsg(fmt.Sprintf("docID: %d, GetGE[ID: %d, value: %v], value: %v",
			id, key, v, i.value))
	}
	return f
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
