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
	v := iter.Current().(*datastruct.Element).Value()
	if v == nil {
		return false
	}

	element := iter.GetGE(id)
	if element == nil {
		return false
	}
	key := element.(*datastruct.Element).Key()
	if key == nil {
		return false
	}

	if k := key.(document.DocId); helpers.Compare(k, id) != 0 {
		return false
	}
	v = iter.Current().(*datastruct.Element).Value()
	return UtilCheck(v, c.op, c.value)
}
