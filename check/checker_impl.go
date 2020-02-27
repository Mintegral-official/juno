package check

import (
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
	"strconv"
)

type CheckerImpl struct {
	si       datastruct.Iterator
	value    interface{}
	op       operation.OP
	e        operation.Operation
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

func (c *CheckerImpl) DebugInfo() string {
	tmp := false
	if c.e != nil {
		tmp = true
	}
	var opMap = map[operation.OP]string{
		operation.EQ:  "=",   // 相等
		operation.NE:  "!=",  // 不等
		operation.LE:  "<=",  // 小于等于
		operation.GE:  ">=",  // 大于等于
		operation.LT:  "<",   // 小于
		operation.GT:  ">",   // 大于
		operation.AND: "and", // 与
		operation.OR:  "or",  // 或
		operation.NOT: "not", // 非
		operation.IN:  "in",  // 范围
	}
	return "FieldName: " + c.si.(*datastruct.SkipListIterator).FieldName + "\t" +
		"value: " + fmt.Sprintf("%v", c.value) + "\t" +
		"OP: " + opMap[c.op] + "\t" +
		"defined operation: " + strconv.FormatBool(tmp) + "\t" +
		"transfer: " + strconv.FormatBool(c.transfer)

}

func (c *CheckerImpl) Check(id document.DocId) bool {
	if c == nil {
		return true
	}
	element := c.si.GetGE(id)
	if element == nil {
		return false
	}
	key, v := element.Key(), element.Value()
	if key != id || v == nil {
		return false
	}
	if c.transfer {
		return UtilCheck(c.value, c.op, v, c.e)
	}
	return UtilCheck(v, c.op, c.value, c.e)
}

func (c *CheckerImpl) Marshal(idx *index.Indexer) map[string]interface{} {
	storageIdx := idx.GetStorageIndex().(*index.StorageIndexer)
	if len(storageIdx.GetFieldName()) == 0 {
		return nil
	}
	fieldName := storageIdx.GetFieldName()
	res := make(map[string]interface{}, 1)
	var tmp []interface{}
	tmp = append(tmp, fieldName[0])
	tmp = append(tmp, c.value)
	tmp = append(tmp, c.op)
	if c.e != nil {
		tmp = append(tmp, 1)
	} else {
		tmp = append(tmp, 0)
	}
	tmp = append(tmp, c.transfer)
	res["=_check"] = tmp
	fieldName = append(fieldName[:0], fieldName[1:]...)
	return res
}

func (c *CheckerImpl) Unmarshal(idx *index.Indexer, res map[string]interface{}, e operation.Operation) Checker {
	v, ok := res["=_check"]
	if !ok {
		return nil
	}
	value := v.([]interface{})
	if value[3] == 1 {
		return NewChecker(idx.GetStorageIndex().Iterator(value[0].(string)), value[1], value[2].(operation.OP), e, value[4].(bool))
	}
	return NewChecker(idx.GetStorageIndex().Iterator(value[0].(string)), value[1], value[2].(operation.OP), nil, value[4].(bool))
}
