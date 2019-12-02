package check

import (
	"github.com/Mintegral-official/juno/document"
)

type AndCheckerImpl struct {
	c []Checker
}

func NewAndCheckerImpl(c []Checker) *AndCheckerImpl {
	if c == nil {
		return nil
	}
	return &AndCheckerImpl{
		c: c,
	}
}

func (a *AndCheckerImpl) Check(id document.DocId) bool {
	if a == nil {
		return true
	}
	for _, cValue := range a.c {
		if !cValue.Check(id) {
			return false
		}
	}
	return true
}
