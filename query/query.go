package query

type Query struct {
	Exp Expression
}

func (q *Query) HasNext() bool {
	return q.Exp.HasNext()
}

func (q *Query) Next() DocId {
	return q.Exp.Next()
}

func NewQuery(s string) *Query {
	return nil
}
