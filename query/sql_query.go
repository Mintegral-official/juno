package query

import (
	"github.com/MintegralTech/juno/check"
	"github.com/MintegralTech/juno/datastruct"
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/index"
	"github.com/MintegralTech/juno/operation"
	"strconv"
	"strings"
)

type SqlQuery struct {
	Node       *datastruct.TreeNode
	Stack      *datastruct.Stack
	Expression *Expression
	e          operation.Operation
	transfer   bool
	debugs     *debug.Debugs
}

func NewSqlQuery(str string, e operation.Operation, transfer bool) (s *SqlQuery) {
	s = &SqlQuery{
		Node:       &datastruct.TreeNode{},
		Stack:      datastruct.NewStack(),
		Expression: NewExpression(str),
		e:          e,
		transfer:   transfer,
	}
	return s
}

func (sq *SqlQuery) exp2Tree() *datastruct.TreeNode {
	exp := sq.Expression.ToPostfix(sq.Expression.string2Strings())
	if sq.debugs != nil {
		sq.debugs.DebugInfo.AddDebugMsg("the exp has ", strconv.Itoa(len(exp)), " conditions")
	}
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
	node := sq.exp2Tree().To()
	for !sq.Stack.Empty() || !node.Empty() {
		if !node.Empty() {
			if node.Peek() != "&" && node.Peek() != "|" {
				if strings.Contains(node.Peek().(string), "@") {
					sq.Stack.Push(sq.parseIn(node.Pop().(string), idx))
				}
				if strings.Contains(node.Peek().(string), "#") {
					sq.Stack.Push(sq.parseNotIn(node.Pop().(string), idx))
					//tmp = 1
				}
				if strings.Contains(node.Peek().(string), "=") &&
					!strings.Contains(node.Peek().(string), ">") &&
					!strings.Contains(node.Peek().(string), "<") &&
					!strings.Contains(node.Peek().(string), "!") {
					sq.Stack.Push(sq.parseEQ(node.Pop().(string), idx))
				}
				if strings.Contains(node.Peek().(string), "!=") {
					sq.Stack.Push(sq.parseNE(node.Pop().(string), idx))
					//	tmp=1
				}
				if strings.Contains(node.Peek().(string), "<") &&
					!strings.Contains(node.Peek().(string), "=") {
					sq.Stack.Push(sq.parseLT(node.Pop().(string), idx))
				}
				if strings.Contains(node.Peek().(string), "<=") {
					sq.Stack.Push(sq.parseLE(node.Pop().(string), idx))
				}
				if strings.Contains(node.Peek().(string), ">") &&
					!strings.Contains(node.Peek().(string), "=") {
					sq.Stack.Push(sq.parseGT(node.Pop().(string), idx))
				}
				if strings.Contains(node.Peek().(string), ">=") {
					sq.Stack.Push(sq.parseGE(node.Pop().(string), idx))
				}
			} else if node.Peek() == "&" {
				//	if tmp == 1 {
				//		sq.Stack.Push(NewNotAndQuery([]Query{sq.Stack.Pop().(Query), sq.Stack.Pop().(Query)}, nil))
				//		tmp = 0
				//	} else {
				sq.Stack.Push(NewAndQuery([]Query{sq.Stack.Pop().(Query), sq.Stack.Pop().(Query)}, nil))
				//	}
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

func (sq *SqlQuery) parseIn(str string, idx *index.Indexer) Query {
	strSlice, storageIdx := strings.Split(str, "@"), idx.GetStorageIndex()
	s := strings.Split(strings.Trim(strings.Trim(strSlice[1], "["), "]"), ",")
	var (
		value = changeType(idx, strSlice[0], s...)
		c     = make([]check.Checker, 1)
	)
	c = append(c, check.NewInChecker(storageIdx.Iterator(strSlice[0]), value, sq.e, sq.transfer))
	return NewAndQuery([]Query{NewTermQuery(storageIdx.Iterator(strSlice[0])),}, c, )
}

func (sq *SqlQuery) parseNotIn(str string, idx *index.Indexer) Query {
	strSlice, storageIdx := strings.Split(str, "#"), idx.GetStorageIndex()
	s := strings.Split(strings.Trim(strings.Trim(strSlice[1], "["), "]"), ",")
	var (
		value = changeType(idx, strSlice[0], s...)
		c     = make([]check.Checker, 1)
	)
	c = append(c, check.NewNotChecker(storageIdx.Iterator(strSlice[0]), value, sq.e, sq.transfer))
	return NewAndQuery([]Query{NewTermQuery(storageIdx.Iterator(strSlice[0])),}, c, )
}

func (sq *SqlQuery) parseEQ(str string, idx *index.Indexer) Query {
	strSlice, invert := strings.Split(str, "="), idx.GetInvertedIndex()
	return NewTermQuery(invert.Iterator(strSlice[0], strSlice[1]))
}

func (sq *SqlQuery) parseNE(str string, idx *index.Indexer) Query {
	strSlice, storageIdx := strings.Split(str, "!="), idx.GetStorageIndex()
	var (
		value = changeType(idx, strSlice[0], strSlice[1])
		c     = make([]check.Checker, 1)
	)
	c = append(c, check.NewChecker(storageIdx.Iterator(strSlice[0]), value[0], operation.NE, sq.e, sq.transfer))
	return NewAndQuery([]Query{NewTermQuery(storageIdx.Iterator(strSlice[0])),}, c, )
}

func (sq *SqlQuery) parseLT(str string, idx *index.Indexer) Query {
	strSlice, storageIdx := strings.Split(str, "<"), idx.GetStorageIndex()
	var (
		value = changeType(idx, strSlice[0], strSlice[1])
		c     = make([]check.Checker, 1)
	)
	c = append(c, check.NewChecker(storageIdx.Iterator(strSlice[0]), value[0], operation.LT, sq.e, sq.transfer))
	return NewAndQuery([]Query{NewTermQuery(storageIdx.Iterator(strSlice[0])),}, c, )
}

func (sq *SqlQuery) parseLE(str string, idx *index.Indexer) Query {
	strSlice, storageIdx := strings.Split(str, "<="), idx.GetStorageIndex()
	var (
		value = changeType(idx, strSlice[0], strSlice[1])
		c     = make([]check.Checker, 1)
	)
	c = append(c, check.NewChecker(storageIdx.Iterator(strSlice[0]), value[0], operation.LE, sq.e, sq.transfer))
	return NewAndQuery([]Query{NewTermQuery(storageIdx.Iterator(strSlice[0])),}, c, )
}

func (sq *SqlQuery) parseGT(str string, idx *index.Indexer) Query {
	strSlice, storageIdx := strings.Split(str, ">"), idx.GetStorageIndex()
	var (
		value = changeType(idx, strSlice[0], strSlice[1])
		c     = make([]check.Checker, 1)
	)
	c = append(c, check.NewChecker(storageIdx.Iterator(strSlice[0]), value[0], operation.GT, sq.e, sq.transfer))
	return NewAndQuery([]Query{NewTermQuery(storageIdx.Iterator(strSlice[0])),}, c, )
}

func (sq *SqlQuery) parseGE(str string, idx *index.Indexer) Query {
	strSlice, storageIdx := strings.Split(str, ">="), idx.GetStorageIndex()
	var (
		value = changeType(idx, strSlice[0], strSlice[1])
		c     = make([]check.Checker, 1)
	)
	c = append(c, check.NewChecker(storageIdx.Iterator(strSlice[0]), value[0], operation.GE, sq.e, sq.transfer))
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
	case document.MapFieldType:
		fallthrough
	default:
		for _, v := range value {
			res = append(res, v)
		}
		return res
	}
}
