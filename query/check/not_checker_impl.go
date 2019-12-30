package check

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/query/operation"
)

type NotChecker struct {
	si    datastruct.Iterator
	value []interface{}
}

func NewNotChecker(si datastruct.Iterator, value ...interface{}) *NotChecker {
	return &NotChecker{
		si:    si,
		value: value,
	}
}

func (nc *NotChecker) Check(id document.DocId) bool {
	if nc == nil {
		return true
	}
	iter := nc.si
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
	return !o.In(nc.value)
}
