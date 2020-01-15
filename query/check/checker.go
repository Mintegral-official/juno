package check

import (
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/query/operation"
)

type Checker interface {
	Check(id document.DocId) bool
}

func UtilCheck(cValue interface{}, op operation.OP, value interface{}, e operation.Operation) bool {
	if e == nil {
		e = operation.NewOperations(cValue)
	} else {
		e.SetValue(cValue)
	}
	switch op {
	case operation.EQ:
		return e.Equal(value)
	case operation.NE:
		return !e.Equal(value)
	case operation.LE:
		return e.Equal(value) || e.Less(value)
	case operation.GE:
		return !e.Less(value)
	case operation.LT:
		return e.Less(value)
	case operation.GT:
		return !e.Equal(value) && !e.Less(value)
	default:
		return false
	}
}
