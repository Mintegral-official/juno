package check

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/query/operation"
)

type NotChecker struct {
	si    datastruct.Iterator
	value []interface{}
	e     operation.Operation
}

func NewNotChecker(si datastruct.Iterator, value []interface{}, e operation.Operation) *NotChecker {
	return &NotChecker{
		si:    si,
		value: value,
		e:     e,
	}
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
		o := operation.Operations{FieldValue: v}
		return !o.In(nc.value)
	}
	nc.e.SetValue(v)
	return nc.e.In(nc.value)
}
