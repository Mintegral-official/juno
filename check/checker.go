package check

import (
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/index"
	"github.com/MintegralTech/juno/operation"
)

type Checker interface {
	Check(id document.DocId) bool
	DebugInfo() *debug.Debug
	SetDebug(level int)
	Marshal() map[string]interface{}
	Unmarshal(idx *index.Indexer, res map[string]interface{}, e operation.Operation) Checker
}

var OpMap = map[operation.OP]string{
	operation.EQ:  "=",   // 相等
	operation.NE:  "!=",  // 不等
	operation.LE:  "<=",  // 小于等于
	operation.GE:  ">=",  // 大于等于
	operation.LT:  "<",   // 小于
	operation.GT:  ">",   // 大于
	operation.AND: "and", // 与
	operation.OR:  "or",  // 或
	operation.NOT: "not", // 非
	operation.IN:  "in",  // 范围
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
