package query

import "github.com/Mintegral-official/juno/document"

type InQuery struct {
}

func NewInQuery() *InQuery {
	return &InQuery{}
}

func (i *InQuery) HasNext() bool {
	panic("implement me")
}

func (i *InQuery) Next() document.DocId {
	panic("implement me")
}

func (t *InQuery) String() string {
	return ""
}
