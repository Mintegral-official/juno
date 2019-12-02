package check

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
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
	iter := c.si
	if UtilCheck(iter.Current().(*datastruct.Element).Value(), c.op, c.value) {
		for iter.HasNext() {
			element := iter.Next()
			if helpers.Compare(id, element.(*datastruct.Element).Key()) == 0 {
				return true
			} else if helpers.Compare(id, element.(*datastruct.Element).Key()) > 0 {
				return false
			}
		}
	}
	return false
}
