package check

import (
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
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

func (c *CheckerImpl) Unmarshal(idx *index.Indexer, res map[string]interface{}, e operation.Operation) Checker {
	v, ok := res["check"]
	if !ok {
		return nil
	}
	value := v.([]interface{})
	if value[3] == 1 {
		return NewChecker(idx.GetStorageIndex().Iterator(value[0].(string)), value[1], value[2].(operation.OP), e, value[4].(bool))
	}
	return NewChecker(idx.GetStorageIndex().Iterator(value[0].(string)), value[1], value[2].(operation.OP), nil, value[4].(bool))
}
