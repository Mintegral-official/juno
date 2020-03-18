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

type CheckerImpl struct {
	si       datastruct.Iterator
	value    interface{}
	op       operation.OP
	e        operation.Operation
	aDebug   *debug.Debug
	transfer bool
}

func NewChecker(si datastruct.Iterator, value interface{}, op operation.OP, e operation.Operation, transfer bool) *CheckerImpl {
	return &CheckerImpl{
		si:       si,
		value:    value,
		op:       op,
		e:        e,
		transfer: transfer,
	}
}

func (c *CheckerImpl) DebugInfo() *debug.Debug {
	if c.aDebug != nil {
		c.aDebug.FieldName = c.si.(*datastruct.SkipListIterator).FieldName
		return c.aDebug
	}
	return nil
}

func (c *CheckerImpl) SetDebug(level int) {
	if c.aDebug == nil {
		c.aDebug = debug.NewDebug(level, "checker")
	}
}

func (c *CheckerImpl) Check(id document.DocId) bool {
	if c == nil {
		return true
	}
	element := c.si.GetGE(id)
	if element == nil {
		if c.aDebug != nil {
			c.aDebug.AddDebugMsg(fmt.Sprintf("docID: %d, value:%v, operation: %s, reason: %s",
				id, c.value, OpMap[c.op], helpers.ElementNotfound))
		}
		return false
	}
	key, v := element.Key(), element.Value()
	if key != id || v == nil {
		if c.aDebug != nil {
			c.aDebug.AddDebugMsg(fmt.Sprintf("docID: %d, GetGE[ID: %d, value: %v], value: %v, operation: %s",
				id, key, v, c.value, OpMap[c.op]))
		}
		return false
	}

	var f bool
	if c.transfer {
		f = UtilCheck(c.value, c.op, v, c.e)
		if c.aDebug != nil && f == false {
			c.aDebug.AddDebugMsg(fmt.Sprintf("docID: %d, GetGE[ID: %d, value: %v], value: %v, operation: %s",
				id, key, v, c.value, OpMap[c.op]))
		}
		return f
	}
	f = UtilCheck(v, c.op, c.value, c.e)
	if c.aDebug != nil && f == false {
		c.aDebug.AddDebugMsg(fmt.Sprintf("docID: %d, GetGE[ID: %d, value: %v], value: %v, operation: %s",
			id, key, v, c.value, OpMap[c.op]))
	}
	return f
}

func (c *CheckerImpl) Marshal() map[string]interface{} {
	res := make(map[string]interface{}, 1)
	var tmp []interface{}
	tmp = append(tmp, c.si.(*datastruct.SkipListIterator).FieldName)
	tmp = append(tmp, c.value)
	tmp = append(tmp, c.op)
	if c.e != nil {
		tmp = append(tmp, 1)
	} else {
		tmp = append(tmp, 0)
	}
	tmp = append(tmp, c.transfer)
	tmp = append(tmp, OpMap[c.op])
	res["check"] = tmp
	return res
}

func (c *CheckerImpl) MarshalV2() *marshal.MarshalInfo {
	if c == nil {
		return nil
	}
	info := &marshal.MarshalInfo{
		Name:       c.si.(*datastruct.SkipListIterator).FieldName,
		QueryValue: c.value,
		Operation:  OpMap[c.op]+"_check",
		Op:         c.op,
		Transfer:   c.transfer,
		Nodes:      nil,
	}
	if c.e != nil {
		info.SelfOperation = true
	} else {
		info.SelfOperation = false
	}
	return info
}

func (c *CheckerImpl) UnmarshalV2(idx index.Index, info *marshal.MarshalInfo) Checker {
	if info == nil {
		return nil
	}
	if info.SelfOperation {
		return NewChecker(idx.GetStorageIndex().Iterator(info.Name), info.QueryValue, info.Op,
			register.FieldMap[info.Name], info.Transfer)
	}
	return NewChecker(idx.GetStorageIndex().Iterator(info.Name), info.QueryValue, info.Op,
		nil, info.Transfer)
}

func (c *CheckerImpl) Unmarshal(idx index.Index, res map[string]interface{}) Checker {
	v, ok := res["check"]
	if !ok {
		return nil
	}
	value := v.([]interface{})
	if value[3] == 1 {
		return NewChecker(idx.GetStorageIndex().Iterator(value[0].(string)),
			value[1], value[2].(operation.OP), register.FieldMap[value[0].(string)], value[4].(bool))
	}
	return NewChecker(idx.GetStorageIndex().Iterator(value[0].(string)),
		value[1], value[2].(operation.OP), nil, value[4].(bool))
}
