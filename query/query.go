package query

import "github.com/Mintegral-official/juno/index"

type Query struct {
	Exp Expression
}

func (q *Query) HasNext() bool {
	return q.Exp.HasNext()
}

func (q *Query) Next() index.DocInfo {
	return q.Exp.Next()
}

func NewQuery(s string) *Query {
	return nil
}
