package check

import (
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
)

type Checker interface {
	Check(id document.DocId) bool
	DebugInfo() string
	Marshal(idx *index.Indexer) map[string]interface{}
	Unmarshal(idx *index.Indexer, res map[string]interface{}, e operation.Operation) Checker
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
