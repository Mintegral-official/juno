package juno

import (
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/query"
)

type Query struct {
	Exp query.Expression
}

func (q *Query) HasNext() bool {
	return q.Exp.HasNext()
}

func (q *Query) Next() document.DocId {
	return q.Exp.Next()
}

func NewQuery(s string) *Query {
	return nil
}
