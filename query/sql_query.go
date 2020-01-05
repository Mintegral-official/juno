package query

import (
	"fmt"
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

func (sq *SqlQuery) LRD(impl *index.Indexer) Query {
	node, tmp := sq.exp2Tree().To(), 0
	for !sq.Stack.Empty() || !node.Empty() {
		if !node.Empty() {
			if node.Peek() != "&" && node.Peek() != "|" {
				if strings.Contains(node.Peek().(string), "@") {
					sq.Stack.Push(parseIn(node.Pop().(string), impl))
				}
				if strings.Contains(node.Peek().(string), "#") {
					sq.Stack.Push(parseNotIn(node.Pop().(string), impl))
					tmp = 1
				}
				if strings.Contains(node.Peek().(string), "=") &&
					!strings.Contains(node.Peek().(string), ">") &&
					!strings.Contains(node.Peek().(string), "<") &&
					!strings.Contains(node.Peek().(string), "!") {
					sq.Stack.Push(parseEQ(node.Pop().(string), impl))
				}
				if strings.Contains(node.Peek().(string), "!=") {
					sq.Stack.Push(parseNE(node.Pop().(string), impl))
				}
				if strings.Contains(node.Peek().(string), "<") &&
					!strings.Contains(node.Peek().(string), "=") {
					sq.Stack.Push(parseLT(node.Pop().(string), impl))
				}
				if strings.Contains(node.Peek().(string), "<=") {
					sq.Stack.Push(parseLE(node.Pop().(string), impl))
				}
				if strings.Contains(node.Peek().(string), ">") &&
					!strings.Contains(node.Peek().(string), "=") {
					sq.Stack.Push(parseGT(node.Pop().(string), impl))
				}
				if strings.Contains(node.Peek().(string), ">=") {
					sq.Stack.Push(parseGE(node.Pop().(string), impl))
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

func parseIn(str string, impl *index.Indexer) Query {
	strSlice, storageIdx := strings.Split(str, "@"), impl.GetStorageIndex()
	s := strings.Split(strings.Trim(strings.Trim(strSlice[1], "["), "]"), ",")
	var (
		typ  = impl.GetDataType(strSlice[0])
		iter = storageIdx.Iterator(strSlice[0])
		c    = make([]check.Checker, 1)
	)
	for i := 1; i < len(s); i++ {
		switch typ {
		case document.BoolFieldType:
			b, _ := strconv.ParseBool(s[i])
			c = append(c, check.NewInChecker(iter, b))
		case document.Int8FieldType:
			b, _ := strconv.ParseInt(s[i], 10, 8)
			c = append(c, check.NewInChecker(iter, int8(b)))
		case document.Int16FieldType:
			b, _ := strconv.ParseInt(s[i], 10, 16)
			c = append(c, check.NewInChecker(iter, int16(b)))
		case document.Int32FieldType:
			b, _ := strconv.ParseInt(s[i], 10, 32)
			c = append(c, check.NewInChecker(iter, int32(b)))
		case document.IntFieldType:
			b, _ := strconv.Atoi(s[i])
			c = append(c, check.NewInChecker(iter, b))
		case document.Int64FieldType:
			b, _ := strconv.ParseInt(s[i], 10, 64)
			c = append(c, check.NewInChecker(iter, b))
		case document.Float32FieldType:
			b, _ := strconv.ParseFloat(s[i], 32)
			c = append(c, check.NewInChecker(iter, float32(b)))
		case document.Float64FieldType:
			b, _ := strconv.ParseInt(s[i], 10, 64)
			c = append(c, check.NewInChecker(iter, b))
		case document.StringFieldType:
			c = append(c, check.NewInChecker(iter, s[i]))
		default:
			panic(fmt.Sprintf("the data[%v] type[%T] is wrong.", s[0], typ))
		}
	}
	return NewAndQuery([]Query{NewTermQuery(iter),}, c, )
}

func parseNotIn(str string, impl *index.Indexer) Query {
	strSlice, storageIdx := strings.Split(str, "#"), impl.GetStorageIndex()
	s := strings.Split(strings.Trim(strings.Trim(strSlice[1], "["), "]"), ",")
	var (
		typ  = impl.GetDataType(strSlice[0])
		iter = storageIdx.Iterator(strSlice[0])
		c    = make([]check.Checker, 1)
	)
	for i := 1; i < len(s); i++ {
		switch typ {
		case document.BoolFieldType:
			b, _ := strconv.ParseBool(s[i])
			c = append(c, check.NewNotChecker(iter, b))
		case document.Int8FieldType:
			b, _ := strconv.ParseInt(s[i], 10, 8)
			c = append(c, check.NewNotChecker(iter, int8(b)))
		case document.Int16FieldType:
			b, _ := strconv.ParseInt(s[i], 10, 16)
			c = append(c, check.NewNotChecker(iter, int16(b)))
		case document.Int32FieldType:
			b, _ := strconv.ParseInt(s[i], 10, 32)
			c = append(c, check.NewNotChecker(iter, int32(b)))
		case document.IntFieldType:
			b, _ := strconv.Atoi(s[i])
			c = append(c, check.NewNotChecker(iter, b))
		case document.Int64FieldType:
			b, _ := strconv.ParseInt(s[i], 10, 64)
			c = append(c, check.NewNotChecker(iter, b))
		case document.Float32FieldType:
			b, _ := strconv.ParseFloat(s[i], 32)
			c = append(c, check.NewNotChecker(iter, float32(b)))
		case document.Float64FieldType:
			b, _ := strconv.ParseInt(s[i], 10, 64)
			c = append(c, check.NewNotChecker(iter, b))
		case document.StringFieldType:
			c = append(c, check.NewNotChecker(iter, strSlice[i]))
		default:
			panic(fmt.Sprintf("the data[%v] type[%T] is wrong.", strSlice[0], typ))
		}
	}
	return NewAndQuery([]Query{NewTermQuery(iter),}, c, )

}

func parseEQ(str string, impl *index.Indexer) Query {
	strSlice, invert := strings.Split(str, "="), impl.GetInvertedIndex()
	return NewTermQuery(invert.Iterator(strSlice[0], strSlice[1]))
}

func parseNE(str string, impl *index.Indexer) Query {
	strSlice, storageIdx := strings.Split(str, "!="), impl.GetStorageIndex()
	var (
		typ  = impl.GetDataType(strSlice[0])
		iter = storageIdx.Iterator(strSlice[0])
		c    = make([]check.Checker, 1)
	)
	for i := 1; i < len(strSlice); i++ {
		switch typ {
		case document.BoolFieldType:
			b, _ := strconv.ParseBool(strSlice[i])
			c = append(c, check.NewNotChecker(iter, b))
		case document.Int8FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 8)
			c = append(c, check.NewNotChecker(iter, int8(b)))
		case document.Int16FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 16)
			c = append(c, check.NewNotChecker(iter, int16(b)))
		case document.Int32FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 32)
			c = append(c, check.NewNotChecker(iter, int32(b)))
		case document.IntFieldType:
			b, _ := strconv.Atoi(strSlice[i])
			c = append(c, check.NewNotChecker(iter, b))
		case document.Int64FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 64)
			c = append(c, check.NewNotChecker(iter, b))
		case document.Float32FieldType:
			b, _ := strconv.ParseFloat(strSlice[i], 32)
			c = append(c, check.NewNotChecker(iter, float32(b)))
		case document.Float64FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 64)
			c = append(c, check.NewNotChecker(iter, b))
		case document.StringFieldType:
			c = append(c, check.NewNotChecker(iter, strSlice[i]))
		default:
			panic(fmt.Sprintf("the data[%v] type[%T] is wrong.", strSlice[0], typ))
		}
	}
	return NewAndQuery([]Query{NewTermQuery(iter),}, c, )
}

func parseLT(str string, impl *index.Indexer) Query {
	strSlice, storageIdx := strings.Split(str, "<"), impl.GetStorageIndex()
	var (
		typ  = impl.GetDataType(strSlice[0])
		iter = storageIdx.Iterator(strSlice[0])
		c    = make([]check.Checker, 1)
	)
	for i := 1; i < len(strSlice); i++ {
		switch typ {
		case document.BoolFieldType:
			b, _ := strconv.ParseBool(strSlice[i])
			c = append(c, check.NewChecker(iter, b, operation.LT))
		case document.Int8FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 8)
			c = append(c, check.NewChecker(iter, int8(b), operation.LT))
		case document.Int16FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 16)
			c = append(c, check.NewChecker(iter, int16(b), operation.LT))
		case document.Int32FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 32)
			c = append(c, check.NewChecker(iter, int32(b), operation.LT))
		case document.IntFieldType:
			b, _ := strconv.Atoi(strSlice[i])
			c = append(c, check.NewChecker(iter, b, operation.LT))
		case document.Int64FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 64)
			c = append(c, check.NewChecker(iter, b, operation.LT))
		case document.Float32FieldType:
			b, _ := strconv.ParseFloat(strSlice[i], 32)
			c = append(c, check.NewChecker(iter, float32(b), operation.LT))
		case document.Float64FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 64)
			c = append(c, check.NewChecker(iter, b, operation.LT))
		case document.StringFieldType:
			c = append(c, check.NewChecker(iter, strSlice[i], operation.LT))
		default:
			panic(fmt.Sprintf("the data[%v] type[%T] is wrong.", strSlice[0], typ))
		}
	}
	return NewAndQuery([]Query{NewTermQuery(iter),}, c, )
}

func parseLE(str string, impl *index.Indexer) Query {
	strSlice, storageIdx := strings.Split(str, "<="), impl.GetStorageIndex()
	var (
		typ  = impl.GetDataType(strSlice[0])
		iter = storageIdx.Iterator(strSlice[0])
		c    = make([]check.Checker, 1)
	)
	for i := 1; i < len(strSlice); i++ {
		switch typ {
		case document.BoolFieldType:
			b, _ := strconv.ParseBool(strSlice[i])
			c = append(c, check.NewChecker(iter, b, operation.LE))
		case document.Int8FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 8)
			c = append(c, check.NewChecker(iter, int8(b), operation.LE))
		case document.Int16FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 16)
			c = append(c, check.NewChecker(iter, int16(b), operation.LE))
		case document.Int32FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 32)
			c = append(c, check.NewChecker(iter, int32(b), operation.LE))
		case document.IntFieldType:
			b, _ := strconv.Atoi(strSlice[i])
			c = append(c, check.NewChecker(iter, b, operation.LE))
		case document.Int64FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 64)
			c = append(c, check.NewChecker(iter, b, operation.LE))
		case document.Float32FieldType:
			b, _ := strconv.ParseFloat(strSlice[i], 32)
			c = append(c, check.NewChecker(iter, float32(b), operation.LE))
		case document.Float64FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 64)
			c = append(c, check.NewChecker(iter, b, operation.LE))
		case document.StringFieldType:
			c = append(c, check.NewChecker(iter, strSlice[i], operation.LE))
		default:
			panic(fmt.Sprintf("the data[%v] type[%T] is wrong.", strSlice[0], typ))
		}
	}
	return NewAndQuery([]Query{NewTermQuery(iter),}, c, )
}

func parseGT(str string, impl *index.Indexer) Query {
	strSlice, storageIdx := strings.Split(str, ">"), impl.GetStorageIndex()
	var (
		typ  = impl.GetDataType(strSlice[0])
		iter = storageIdx.Iterator(strSlice[0])
		c    = make([]check.Checker, 1)
	)
	for i := 1; i < len(strSlice); i++ {
		switch typ {
		case document.BoolFieldType:
			b, _ := strconv.ParseBool(strSlice[i])
			c = append(c, check.NewChecker(iter, b, operation.GT))
		case document.Int8FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 8)
			c = append(c, check.NewChecker(iter, int8(b), operation.GT))
		case document.Int16FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 16)
			c = append(c, check.NewChecker(iter, int16(b), operation.GT))
		case document.Int32FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 32)
			c = append(c, check.NewChecker(iter, int32(b), operation.GT))
		case document.IntFieldType:
			b, _ := strconv.Atoi(strSlice[i])
			c = append(c, check.NewChecker(iter, b, operation.GT))
		case document.Int64FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 64)
			c = append(c, check.NewChecker(iter, b, operation.GT))
		case document.Float32FieldType:
			b, _ := strconv.ParseFloat(strSlice[i], 32)
			c = append(c, check.NewChecker(iter, float32(b), operation.GT))
		case document.Float64FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 64)
			c = append(c, check.NewChecker(iter, b, operation.GT))
		case document.StringFieldType:
			c = append(c, check.NewChecker(iter, strSlice[i], operation.GT))
		default:
			panic(fmt.Sprintf("the data[%v] type[%T] is wrong.", strSlice[0], typ))
		}
	}
	return NewAndQuery([]Query{NewTermQuery(iter),}, c, )
}

func parseGE(str string, impl *index.Indexer) Query {
	strSlice, storageIdx := strings.Split(str, ">="), impl.GetStorageIndex()
	var (
		typ  = impl.GetDataType(strSlice[0])
		iter = storageIdx.Iterator(strSlice[0])
		c    = make([]check.Checker, 1)
	)
	for i := 1; i < len(strSlice); i++ {
		switch typ {
		case document.BoolFieldType:
			b, _ := strconv.ParseBool(strSlice[i])
			c = append(c, check.NewChecker(iter, b, operation.GE))
		case document.Int8FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 8)
			c = append(c, check.NewChecker(iter, int8(b), operation.GE))
		case document.Int16FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 16)
			c = append(c, check.NewChecker(iter, int16(b), operation.GE))
		case document.Int32FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 32)
			c = append(c, check.NewChecker(iter, int32(b), operation.GE))
		case document.IntFieldType:
			b, _ := strconv.Atoi(strSlice[i])
			c = append(c, check.NewChecker(iter, b, operation.GE))
		case document.Int64FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 64)
			c = append(c, check.NewChecker(iter, b, operation.GE))
		case document.Float32FieldType:
			b, _ := strconv.ParseFloat(strSlice[i], 32)
			c = append(c, check.NewChecker(iter, float32(b), operation.GE))
		case document.Float64FieldType:
			b, _ := strconv.ParseInt(strSlice[i], 10, 64)
			c = append(c, check.NewChecker(iter, b, operation.GE))
		case document.StringFieldType:
			c = append(c, check.NewChecker(iter, strSlice[i], operation.GE))
		default:
			panic(fmt.Sprintf("the data[%v] type[%T] is wrong.", strSlice[0], typ))
		}
	}
	return NewAndQuery([]Query{NewTermQuery(iter),}, c, )
}
