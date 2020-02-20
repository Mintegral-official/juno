package check

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/operation"
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
