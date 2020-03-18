package check

import (
	"fmt"
	"github.com/MintegralTech/juno/datastruct"
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/helpers"
	"github.com/MintegralTech/juno/index"
	"github.com/MintegralTech/juno/marshal"
	"github.com/MintegralTech/juno/operation"
	"github.com/MintegralTech/juno/register"
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
		nc.aDebug.FieldName = nc.si.(*datastruct.SkipListIterator).FieldName
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
			nc.aDebug.AddDebugMsg(fmt.Sprintf("docId: %d, Value: %v, reason: %v",
				id, nc.value, helpers.ElementNotfound))
		}
		return false
	}
	key, v := element.Key(), element.Value()
	if key != id || v == nil {
		if nc.aDebug != nil {
			nc.aDebug.AddDebugMsg(fmt.Sprintf("docID: %d, GetGE[ID: %d, value: %v], value: %v",
				id, key, v, nc.value))
		}
		return false
	}
	var f bool
	if nc.e == nil {
		if nc.transfer {
			o := operation.Operations{FieldValue: nc.value}
			f = !o.In(v)
			if nc.aDebug != nil && f == false {
				nc.aDebug.AddDebugMsg(fmt.Sprintf("docID: %d, GetGE[ID: %d, value: %v], value: %v",
					id, key, v, nc.value))
			}
			return f
		}
		o := operation.Operations{FieldValue: v}
		f = !o.In(nc.value)
		if nc.aDebug != nil && f == false {
			nc.aDebug.AddDebugMsg(fmt.Sprintf("docID: %d, GetGE[ID: %d, value: %v], value: %v",
				id, key, v, nc.value))
		}
		return !o.In(nc.value)
	}
	if nc.transfer {
		nc.e.SetValue(nc.value)
		f = !nc.e.In(v)
		if nc.aDebug != nil && f == false {
			nc.aDebug.AddDebugMsg(fmt.Sprintf("docID: %d, GetGE[ID: %d, value: %v], value: %v",
				id, key, v, nc.value))
		}
		return f
	}
	nc.e.SetValue(v)
	f = !nc.e.In(nc.value)
	if nc.aDebug != nil && f == false {
		nc.aDebug.AddDebugMsg(fmt.Sprintf("docID: %d, GetGE[ID: %d, value: %v], value: %v",
			id, key, v, nc.value))
	}
	return f
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

func (nc *NotChecker) MarshalV2() *marshal.MarshalInfo {
	if nc == nil {
		return nil
	}
	info := &marshal.MarshalInfo{
		Name:       nc.si.(*datastruct.SkipListIterator).FieldName,
		QueryValue: nc.value,
		Operation:  "not_check",
		Transfer:   nc.transfer,
		Nodes:      nil,
	}
	if nc.e != nil {
		info.SelfOperation = true
	} else {
		info.SelfOperation = false
	}
	return info
}

func (nc *NotChecker) UnmarshalV2(idx index.Index, info *marshal.MarshalInfo) Checker {
	if info == nil {
		return nil
	}
	if info.SelfOperation {
		return NewInChecker(idx.GetStorageIndex().Iterator(info.Name), info.QueryValue, register.FieldMap[info.Name],
			info.Transfer)
	}
	return NewInChecker(idx.GetStorageIndex().Iterator(info.Name), info.QueryValue, nil, info.Transfer)
}

func (nc *NotChecker) Unmarshal(idx index.Index, res map[string]interface{}) Checker {
	v, ok := res["not_check"]
	if !ok {
		return nil
	}
	value := v.([]interface{})
	if value[2] == 1 {
		return NewNotChecker(idx.GetStorageIndex().Iterator(value[0].(string)),
			value[1], register.FieldMap[value[0].(string)], value[3].(bool))
	}
	return NewNotChecker(idx.GetStorageIndex().Iterator(value[0].(string)), value[1], nil, value[3].(bool))

}
