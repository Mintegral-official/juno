package check

import (
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/query/operation"
)

type Checker interface {
	Check(id document.DocId) bool
}

func UtilCheck(cValue interface{}, op operation.OP, value interface{}) bool {
	o := operation.Operations{FieldValue: cValue}
	switch op {
	case operation.EQ:
		return o.Equal(value)
	case operation.NE:
		return !o.Equal(value)
	case operation.LE:
		return o.Equal(value) || o.Less(value)
	case operation.GE:
		return !o.Less(value)
	case operation.LT:
		return o.Less(value)
	case operation.GT:
		return !o.Equal(value) && !o.Less(value)
	default:
		return false
	}
}
