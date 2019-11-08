package query

import "github.com/Mintegral-official/juno/document"

type AndQuery struct {
}

func NewAndQuery() *AndQuery {
	return &AndQuery{}
}

func (a *AndQuery) HasNext() bool {
	panic("implement me")
}

func (a *AndQuery) Next() document.DocId {
	panic("implement me")
}

func (t *AndQuery) String() string {
	return ""
}
