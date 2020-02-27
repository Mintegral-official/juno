package check

import (
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
	"strconv"
)

type NotChecker struct {
	si       datastruct.Iterator
	value    interface{}
	e        operation.Operation
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

func (nc *NotChecker) DebugInfo() string {
	tmp := false
	if nc.e != nil {
		tmp = true
	}
	return "FieldName: " + nc.si.(*datastruct.SkipListIterator).FieldName + "\t" +
		"value: " + fmt.Sprintf("%v", nc.value) + "\t" +
		"OP: not" + "\t" +
		"defined operation: " + strconv.FormatBool(tmp) + "\t" +
		"transfer: " + strconv.FormatBool(nc.transfer)

}

func (nc *NotChecker) Check(id document.DocId) bool {
	if nc == nil {
		return true
	}

	element := nc.si.GetGE(id)
	if element == nil {
		return false
	}
	key, v := element.Key(), element.Value()
	if key != id || v == nil {
		return false
	}
	if nc.e == nil {
		if nc.transfer {
			o := operation.Operations{FieldValue: nc.value}
			return !o.In(v)
		}
		o := operation.Operations{FieldValue: v}
		return !o.In(nc.value)
	}
	if nc.transfer {
		nc.e.SetValue(nc.value)
		return !nc.e.In(v)
	}
	nc.e.SetValue(v)
	return !nc.e.In(nc.value)
}

func (nc *NotChecker) Marshal(idx *index.Indexer) map[string]interface{} {
	storageIdx := idx.GetStorageIndex().(*index.StorageIndexer)
	if len(storageIdx.GetFieldName()) == 0 {
		return nil
	}
	fieldName := storageIdx.GetFieldName()
	res := make(map[string]interface{}, 1)
	var tmp []interface{}
	tmp = append(tmp, fieldName[0])
	tmp = append(tmp, nc.value)
	if nc.e != nil {
		tmp = append(tmp, 1)
	} else {
		tmp = append(tmp, 0)
	}
	tmp = append(tmp, nc.transfer)
	res["not_check"] = tmp
	fieldName = append(fieldName[:0], fieldName[1:]...)
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
