package check

import (
	"github.com/Mintegral-official/juno/document"
)

type InCheckerImpl struct {
	c []Checker
}

func NewInCheckerImpl(c []Checker) *InCheckerImpl {
	if c == nil {
		return nil
	}
	return &InCheckerImpl{
		c: c,
	}
}

func (i *InCheckerImpl) Check(id document.DocId) bool {
	if i == nil {
		return true
	}
	for _, cValue := range i.c {
		if cValue.Check(id) {
			return true
		}
	}
	return false
}
