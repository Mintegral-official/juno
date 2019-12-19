package check

import (
	"github.com/Mintegral-official/juno/document"
)

type InInChecker struct {
	c []Checker
}

func NewInInChecker(c []Checker) *InInChecker {
	if c == nil {
		return nil
	}
	return &InInChecker{
		c: c,
	}
}

func (i *InInChecker) Check(id document.DocId) bool {
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
