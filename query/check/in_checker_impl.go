package check

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/query/operation"
)

type InChecker struct {
	si    datastruct.Iterator
	value []interface{}
	e     operation.Operation
}

func NewInChecker(si datastruct.Iterator, value []interface{}, e operation.Operation) *InChecker {
	return &InChecker{
		si:    si,
		value: value,
		e:     e,
	}
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
		o := operation.Operations{FieldValue: v}
		return o.In(i.value)
	}
	i.e.SetValue(v)
	return i.e.In(i.value)
}
