package check

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/query/operation"
)

type InChecker struct {
	si    datastruct.Iterator
	value interface{}
	op    operation.OP
}

func NewChecker(si datastruct.Iterator, value interface{}, op operation.OP) *InChecker {
	return &InChecker{
		si:    si,
		value: value,
		op:    op,
	}
}

func (c *InChecker) Check(id document.DocId) bool {
	if c == nil {
		return true
	}
	iter := c.si
	v := iter.Current().(*datastruct.Element).Value()
	if v == nil {
		return false
	}

	element := iter.GetGE(id)
	if element == nil {
		return false
	}
	key := element.(*datastruct.Element).Key()
	if key == 0 {
		return false
	}

	if key != id {
		return false
	}
	v = iter.Current().(*datastruct.Element).Value()
	return UtilCheck(v, c.op, c.value)
}
