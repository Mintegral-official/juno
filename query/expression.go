package query

import (
	"github.com/Mintegral-official/juno/document"
)

type OP int64

const (
	InvalidDocid = 0xffffffff
)

const (
	EQ  = iota // 相等
	NE         // 不等
	LE         // 小于等于
	GE         // 大于等于
	LT         // 小于
	GT         // 大于
	AND        // 与
	OR         // 或
	NOT        // 非
	IN         // 范围
)

type Expression interface {
	HasNext() bool
	Next() document.DocId
	FindGE(t document.DocId) document.DocId
}

func NewExpression(name string) Expression {
	return nil
}
