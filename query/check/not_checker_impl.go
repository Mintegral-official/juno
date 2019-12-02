package check

import (
	"github.com/Mintegral-official/juno/document"
)

type NotCheckerImpl struct {
	c []Checker
}

func NewNotCheckerImpl(c []Checker) *NotCheckerImpl {
	if c == nil {
		return nil
	}
	return &NotCheckerImpl{
		c: c,
	}
}

func (n *NotCheckerImpl) Check(id document.DocId) bool {
	if n == nil {
		return true
	}
	for _, cValue := range n.c {
		if !cValue.Check(id) {
			return false
		}
	}
	return true
}
