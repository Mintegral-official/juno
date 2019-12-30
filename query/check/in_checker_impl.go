package check

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/query/operation"
)

type InChecker struct {
	si    datastruct.Iterator
	value []interface{}
}

func NewInChecker(si datastruct.Iterator, value ...interface{}) *InChecker {
	return &InChecker{
		si:    si,
		value: value,
	}
}

func (i *InChecker) Check(id document.DocId) bool {
	if i == nil {
		return true
	}
	iter := i.si
	v := iter.Current().(*datastruct.Element).Value()
	if v == nil {
		return false
	}

	element := iter.GetGE(id)
	if element == nil {
		return false
	}
	key := element.(*datastruct.Element).Key()
	if key != id {
		return false
	}
	v = iter.Current().(*datastruct.Element).Value()
	if v == nil {
		return false
	}
	o := operation.Operations{FieldValue: v}
	return o.In(i.value)
}
