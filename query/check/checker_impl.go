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
	element := c.si.GetGE(id)
	if element == nil {
		return false
	}
	key, v := element.Key(), element.Value()
	if key != id || v == nil {
		return false
	}
	return UtilCheck(v, c.op, c.value, c.e)
}
