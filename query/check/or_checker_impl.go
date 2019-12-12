package check

import (
	"github.com/Mintegral-official/juno/document"
)

type OrCheckerImpl struct {
	c []Checker
}

func NewOrCheckerImpl(c []Checker) *OrCheckerImpl {
	if c == nil {
		return nil
	}
	return &OrCheckerImpl{
		c: c,
	}
}

func (o *OrCheckerImpl) Check(id document.DocId) bool {
	if o == nil {
		return true
	}
	for _, cValue := range o.c {
		if cValue.Check(id) {
			return true
		}
	}
	return false
}
