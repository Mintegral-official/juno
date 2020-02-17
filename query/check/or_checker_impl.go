package check

import (
	"github.com/Mintegral-official/juno/document"
)

type OrChecker struct {
	c []Checker
}

func NewOrChecker(c []Checker) *OrChecker {
	if c == nil {
		return nil
	}
	return &OrChecker{
		c: c,
	}
}

func (o *OrChecker) Check(id document.DocId) bool {
	if o == nil {
		return true
	}
	for _, cValue := range o.c {
		if cValue == nil {
			continue
		}
		if cValue.Check(id) {
			return true
		}
	}
	return false
}
