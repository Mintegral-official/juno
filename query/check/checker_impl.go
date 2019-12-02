package check

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/query/operation"
)

type CheckerImpl struct {
	si    *datastruct.SkipListIterator
	value interface{}
	op    operation.OP
}

func NewCheckerImpl(si *datastruct.SkipListIterator, value interface{}, op operation.OP) *CheckerImpl {
	return &CheckerImpl{
		si:    si,
		value: value,
		op:    op,
	}
}

func (c *CheckerImpl) Check(id document.DocId) bool {
	if c == nil {
		return true
	}
	return UtilCheck(id, c.op, c.value)
}
