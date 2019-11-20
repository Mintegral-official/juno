package query

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
)

type CheckerImpl struct {
	sl    datastruct.SkipList
	value interface{}
}

func NewCheckerImpl(sl datastruct.SkipList, value interface{}) *CheckerImpl {
	return &CheckerImpl{
		sl:    sl,
		value: value,
	}
}

func (c CheckerImpl) Check(id document.DocId) bool {
	if v, e := c.sl.Get(id); e == nil {
		return c.sl.Cmp().Compare(v, c.value) == 0
	}
	return false
}
