package query

import "github.com/Mintegral-official/juno/document"

type OrQuery struct {
}

func NewOrQuery() *OrQuery {
	return &OrQuery{}
}

func (o *OrQuery) HasNext() bool {
	panic("implement me")
}

func (o *OrQuery) Next() document.DocId {
	panic("implement me")
}

func (t *OrQuery) String() string {
	return ""
}
