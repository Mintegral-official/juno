package check

import (
	"github.com/Mintegral-official/juno/document"
)

type AndChecker struct {
	c []Checker
}

func NewAndChecker(c []Checker) *AndChecker {
	if c == nil {
		return nil
	}
	return &AndChecker{
		c: c,
	}
}

func (a *AndChecker) Check(id document.DocId) bool {
	if a == nil {
		return true
	}
	for _, cValue := range a.c {
		if cValue == nil {
			continue
		}
		if !cValue.Check(id) {
			return false
		}
	}
	return true
}
