package check

import (
	"github.com/Mintegral-official/juno/document"
)

type NotChecker struct {
	c []Checker
}

func NewNotChecker(c []Checker) *NotChecker {
	if c == nil {
		return nil
	}
	return &NotChecker{
		c: c,
	}
}

func (n *NotChecker) Check(id document.DocId) bool {
	if n == nil {
		return true
	}
	for _, cValue := range n.c {
		if !cValue.Check(id) {
			return true
		}
	}
	return false
}
