package query

import "github.com/Mintegral-official/juno/document"

type QueryImpl struct {
	query Query
}

func (q *QueryImpl) Next() (document.DocId, error) {
	return q.query.Next()
}

func (q *QueryImpl) GetGE(id document.DocId) (document.DocId, error) {
	return q.query.GetGE(id)
}

func (q *QueryImpl) String() string {
	return ""
}

func NewQuery(query Query) *QueryImpl {
	return &QueryImpl{
		query: query,
	}
}
