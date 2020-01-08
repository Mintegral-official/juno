package query

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/query/check"
	"github.com/Mintegral-official/juno/query/operation"
	"strconv"
	"strings"
)

type SqlQuery struct {
	Node       *datastruct.TreeNode
	Stack      *datastruct.Stack
	Expression *Expression
}

func NewSqlQuery(str string) *SqlQuery {
	return &SqlQuery{
		Node:       &datastruct.TreeNode{},
		Stack:      datastruct.NewStack(),
		Expression: NewExpression(str),
	}
}

func (sq *SqlQuery) exp2Tree() *datastruct.TreeNode {
	exp := sq.Expression.ToPostfix(sq.Expression.string2Strings())
	for _, v := range exp {
		if v == "&" || v == "|" {
			if sq.Stack.Empty() || sq.Stack.Len() < 2 {
				panic("the expression is wrong")
			}
			root := &datastruct.TreeNode{
				Data: v,
			}
			root.Left = sq.Stack.Pop().(*datastruct.TreeNode)
			root.Right = sq.Stack.Pop().(*datastruct.TreeNode)
			sq.Stack.Push(root)
		} else {
			sq.Stack.Push(&datastruct.TreeNode{Data: v,})
		}
	}
	if sq.Stack.Empty() || sq.Stack.Len() > 1 {
		panic("the expression is wrong")
	}
	return sq.Stack.Pop().(*datastruct.TreeNode)
}

func (sq *SqlQuery) LRD(idx *index.Indexer) Query {
	node, tmp := sq.exp2Tree().To(), 0
	for !sq.Stack.Empty() || !node.Empty() {
		if !node.Empty() {
			if node.Peek() != "&" && node.Peek() != "|" {
				if strings.Contains(node.Peek().(string), "@") {
					sq.Stack.Push(parseIn(node.Pop().(string), idx))
				}
				if strings.Contains(node.Peek().(string), "#") {
					sq.Stack.Push(parseNotIn(node.Pop().(string), idx))
					//tmp = 1
				}
				if strings.Contains(node.Peek().(string), "=") &&
					!strings.Contains(node.Peek().(string), ">") &&
					!strings.Contains(node.Peek().(string), "<") &&
					!strings.Contains(node.Peek().(string), "!") {
					sq.Stack.Push(parseEQ(node.Pop().(string), idx))
				}
				if strings.Contains(node.Peek().(string), "!=") {
					sq.Stack.Push(parseNE(node.Pop().(string), idx))
					tmp=1
				}
				if strings.Contains(node.Peek().(string), "<") &&
					!strings.Contains(node.Peek().(string), "=") {
					sq.Stack.Push(parseLT(node.Pop().(string), idx))
				}
				if strings.Contains(node.Peek().(string), "<=") {
					sq.Stack.Push(parseLE(node.Pop().(string), idx))
				}
				if strings.Contains(node.Peek().(string), ">") &&
					!strings.Contains(node.Peek().(string), "=") {
					sq.Stack.Push(parseGT(node.Pop().(string), idx))
				}
				if strings.Contains(node.Peek().(string), ">=") {
					sq.Stack.Push(parseGE(node.Pop().(string), idx))
				}
			} else if node.Peek() == "&" {
				if tmp == 1 {
					sq.Stack.Push(NewNotAndQuery([]Query{sq.Stack.Pop().(Query), sq.Stack.Pop().(Query)}, nil))
					tmp = 0
				} else {
					sq.Stack.Push(NewAndQuery([]Query{sq.Stack.Pop().(Query), sq.Stack.Pop().(Query)}, nil))
				}
				node.Pop()
			} else if node.Peek() == "|" {
				sq.Stack.Push(NewOrQuery([]Query{sq.Stack.Pop().(Query), sq.Stack.Pop().(Query)}, nil))
				node.Pop()
			}
		} else {
			return sq.Stack.Pop().(Query)
		}
	}
	return sq.Stack.Pop().(Query)
}

func parseIn(str string, idx *index.Indexer) Query {
	strSlice, storageIdx := strings.Split(str, "@"), idx.GetStorageIndex()
	s := strings.Split(strings.Trim(strings.Trim(strSlice[1], "["), "]"), ",")
	var (
		value = changeType(idx, strSlice[0], s...)
		c     = make([]check.Checker, 1)
	)
	c = append(c, check.NewInChecker(storageIdx.Iterator(strSlice[0]), value...))
	return NewAndQuery([]Query{NewTermQuery(storageIdx.Iterator(strSlice[0])),}, c, )
}

func parseNotIn(str string, idx *index.Indexer) Query {
	strSlice, storageIdx := strings.Split(str, "#"), idx.GetStorageIndex()
	s := strings.Split(strings.Trim(strings.Trim(strSlice[1], "["), "]"), ",")
	var (
		value = changeType(idx, strSlice[0], s...)
		c     = make([]check.Checker, 1)
	)
	c = append(c, check.NewNotChecker(storageIdx.Iterator(strSlice[0]), value...))
	return NewAndQuery([]Query{NewTermQuery(storageIdx.Iterator(strSlice[0])),}, c, )

}

func parseEQ(str string, idx *index.Indexer) Query {
	strSlice, invert := strings.Split(str, "="), idx.GetInvertedIndex()
	return NewTermQuery(invert.Iterator(strSlice[0], strSlice[1]))
}

func parseNE(str string, idx *index.Indexer) Query {
	strSlice, storageIdx := strings.Split(str, "!="), idx.GetStorageIndex()
	var (
		value = changeType(idx, strSlice[0], strSlice[1])
		c     = make([]check.Checker, 1)
	)
	c = append(c, check.NewChecker(storageIdx.Iterator(strSlice[0]), value[0], operation.NE))
	return NewAndQuery([]Query{NewTermQuery(storageIdx.Iterator(strSlice[0])),}, c, )
}

func parseLT(str string, idx *index.Indexer) Query {
	strSlice, storageIdx := strings.Split(str, "<"), idx.GetStorageIndex()
	var (
		value = changeType(idx, strSlice[0], strSlice[1])
		c     = make([]check.Checker, 1)
	)
	c = append(c, check.NewChecker(storageIdx.Iterator(strSlice[0]), value[0], operation.LT))
	return NewAndQuery([]Query{NewTermQuery(storageIdx.Iterator(strSlice[0])),}, c, )
}

func parseLE(str string, idx *index.Indexer) Query {
	strSlice, storageIdx := strings.Split(str, "<="), idx.GetStorageIndex()
	var (
		value = changeType(idx, strSlice[0], strSlice[1])
		c     = make([]check.Checker, 1)
	)
	c = append(c, check.NewChecker(storageIdx.Iterator(strSlice[0]), value[0], operation.LE))
	return NewAndQuery([]Query{NewTermQuery(storageIdx.Iterator(strSlice[0])),}, c, )
}

func parseGT(str string, idx *index.Indexer) Query {
	strSlice, storageIdx := strings.Split(str, ">"), idx.GetStorageIndex()
	var (
		value = changeType(idx, strSlice[0], strSlice[1])
		c     = make([]check.Checker, 1)
	)
	c = append(c, check.NewChecker(storageIdx.Iterator(strSlice[0]), value[0], operation.GT))
	return NewAndQuery([]Query{NewTermQuery(storageIdx.Iterator(strSlice[0])),}, c, )
}

func parseGE(str string, idx *index.Indexer) Query {
	strSlice, storageIdx := strings.Split(str, ">="), idx.GetStorageIndex()
	var (
		value = changeType(idx, strSlice[0], strSlice[1])
		c     = make([]check.Checker, 1)
	)
	c = append(c, check.NewChecker(storageIdx.Iterator(strSlice[0]), value[0], operation.GE))
	return NewAndQuery([]Query{NewTermQuery(storageIdx.Iterator(strSlice[0])),}, c, )
}

func changeType(idx *index.Indexer, name string, value ...string) (res []interface{}) {
	switch typ := idx.GetDataType(name); typ {
	case document.BoolFieldType:
		for _, v := range value {
			if b, err := strconv.ParseBool(v); err == nil {
				res = append(res, b)
			} else {
				panic("the value type is not bool.")
			}
		}
		return res
	case document.IntFieldType:
		for _, v := range value {
			if b, err := strconv.ParseInt(v, 10, 64); err == nil {
				res = append(res, b)
			} else {
				panic("the value type is not int.")
			}
		}
		return res
	case document.FloatFieldType:
		for _, v := range value {
			if b, err := strconv.ParseFloat(v, 64); err == nil {
				res = append(res, b)
			} else {
				panic("the value type is not float.")
			}
		}
		return res
	case document.StringFieldType:
		for _, v := range value {
			res = append(res, v)
		}
		return res
	case document.SliceFieldType:
		fallthrough
	case document.SelfDefinedFieldType:
		fallthrough
	default:
		for _, v := range value {
			res = append(res, v)
		}
		return res
	}
}
