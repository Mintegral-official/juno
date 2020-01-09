package check

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/query/operation"
)

type CheckerImpl struct {
	si    datastruct.Iterator
	value interface{}
	op    operation.OP
	e     operation.Operation
}

func NewChecker(si datastruct.Iterator, value interface{}, op operation.OP, e operation.Operation) *CheckerImpl {
	return &CheckerImpl{
		si:    si,
		value: value,
		op:    op,
		e:     e,
	}
}

func (c *CheckerImpl) Check(id document.DocId) bool {
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
	if key != id {
		return false
	}
	v = iter.Current().(*datastruct.Element).Value()
	if v == nil {
		return false
	}
	return UtilCheck(v, c.op, c.value, c.e)
}
