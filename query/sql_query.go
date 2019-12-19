package query

import (
	"github.com/Mintegral-official/juno/datastruct"
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



//func (sq *SqlQuery) LRD(impl *index.Indexer) Query {
//	node, tmp := sq.exp2Tree(), 0
//	for !sq.Stack.Empty() || node != nil {
//		if node != nil {
//			if node.Data != "&" && node.Data != "|" {
//				if strings.Contains(node.Data.(string), "@") {
//					sq.Stack.Push(parseIn(node.Data.(string), impl))
//				} else if strings.Contains(node.Data.(string), "#") {
//					sq.Stack.Push(parseNotIn(node.Data.(string), impl))
//					tmp = 1
//				} else if strings.Contains(node.Data.(string), "=") {
//					sq.Stack.Push(parseEq(node.Data.(string), impl))
//				} else if strings.Contains(node.Data.(string), "!=") {
//					sq.Stack.Push(parseNE(node.Data.(string), impl))
//				} else if strings.Contains(node.Data.(string), "<") {
//					sq.Stack.Push(parseLT(node.Data.(string), impl))
//				} else if strings.Contains(node.Data.(string), "<=") {
//					sq.Stack.Push(parseLE(node.Data.(string), impl))
//				} else if strings.Contains(node.Data.(string), ">") {
//					sq.Stack.Push(parseGT(node.Data.(string), impl))
//				} else if strings.Contains(node.Data.(string), ">=") {
//					sq.Stack.Push(parseGE(node.Data.(string), impl))
//				}
//			} else if node.Data == "&" {
//				if tmp == 1 {
//					sq.Stack.Push(NewNotAndQuery([]Query{sq.Stack.Pop().(Query), sq.Stack.Pop().(Query)}, nil))
//					tmp = 0
//				} else {
//					sq.Stack.Push(NewAndQuery([]Query{sq.Stack.Pop().(Query), sq.Stack.Pop().(Query)}, nil))
//				}
//			} else if node.Data == "|" {
//				sq.Stack.Push(NewOrQuery([]Query{sq.Stack.Pop().(Query), sq.Stack.Pop().(Query)}, nil))
//			}
//		} else {
//			return sq.Stack.Pop().(Query)
//		}
//	}
//	return sq.Stack.Pop().(Query)
//}
//
//func parseIn(str string, impl *index.Indexer) Query {
//	strSlice, invert := strings.Split(str, "@"), impl.GetInvertedIndex()
//	values := strings.Split(strings.Trim(strings.Trim(strSlice[1], "["), "]"), ",")
//	var querys []Query
//	for _, v := range values {
//		querys = append(querys, NewTermQuery(invert.Iterator(strSlice[0]+"_"+v).(*datastruct.SkipListIterator)))
//	}
//	return NewOrQuery(querys, nil)
//}
//
//func parseNotIn(str string, impl *index.Indexer) Query {
//	strSlice, invert := strings.Split(str, "#"), impl.GetInvertedIndex()
//	values := strings.Split(strings.Trim(strings.Trim(strSlice[1], "["), "]"), ",")
//	var querys []Query
//	for _, v := range values {
//		querys = append(querys, NewTermQuery(invert.Iterator(strSlice[0]+"_"+v).(*datastruct.SkipListIterator)))
//	}
//	return NewOrQuery(querys, nil)
//}
//
//func parseEq(str string, impl *index.Indexer) Query {
//	strSlice, invert := strings.Split(str, "="), impl.GetInvertedIndex()
//	return NewTermQuery(invert.Iterator(strSlice[0] + "_" + strSlice[1]).(*datastruct.SkipListIterator))
//}
//
//func parseNE(str string, impl *index.Indexer) Query {
//	strSlice, storageIdx := strings.Split(str, "!="), impl.GetStorageIndex()
//	return NewAndQuery([]Query{
//		NewTermQuery(storageIdx.Iterator(strSlice[0]).(*datastruct.SkipListIterator)),
//	}, []check.Checker{
//		check.NewInChecker(storageIdx.Iterator(strSlice[0]).(*datastruct.SkipListIterator), strSlice[1], operation.NE),
//	}, )
//}
//
//func parseLT(str string, impl *index.Indexer) Query {
//	strSlice, storageIdx := strings.Split(str, "<"), impl.GetStorageIndex()
//	return NewAndQuery([]Query{
//		NewTermQuery(storageIdx.Iterator(strSlice[0]).(*datastruct.SkipListIterator)),
//	}, []check.Checker{
//		check.NewInChecker(storageIdx.Iterator(strSlice[0]).(*datastruct.SkipListIterator), strSlice[1], operation.LT),
//	}, )
//}
//
//func parseLE(str string, impl *index.Indexer) Query {
//	strSlice, storageIdx := strings.Split(str, "<="), impl.GetStorageIndex()
//	return NewAndQuery([]Query{
//		NewTermQuery(storageIdx.Iterator(strSlice[0]).(*datastruct.SkipListIterator)),
//	}, []check.Checker{
//		check.NewInChecker(storageIdx.Iterator(strSlice[0]).(*datastruct.SkipListIterator), strSlice[1], operation.LE),
//	}, )
//}
//
//func parseGT(str string, impl *index.Indexer) Query {
//	strSlice, storageIdx := strings.Split(str, ">"), impl.GetStorageIndex()
//	return NewAndQuery([]Query{
//		NewTermQuery(storageIdx.Iterator(strSlice[0]).(*datastruct.SkipListIterator)),
//	}, []check.Checker{
//		check.NewInChecker(storageIdx.Iterator(strSlice[0]).(*datastruct.SkipListIterator), strSlice[1], operation.GT),
//	}, )
//}
//
//func parseGE(str string, impl *index.Indexer) Query {
//	strSlice, storageIdx := strings.Split(str, ">="), impl.GetStorageIndex()
//	return NewAndQuery([]Query{
//		NewTermQuery(storageIdx.Iterator(strSlice[0]).(*datastruct.SkipListIterator)),
//	}, []check.Checker{
//		check.NewInChecker(storageIdx.Iterator(strSlice[0]).(*datastruct.SkipListIterator), strSlice[1], operation.GE),
//	}, )
//}
