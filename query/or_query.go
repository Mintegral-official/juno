package query

import "github.com/Mintegral-official/juno/document"

type OrQuery struct {
}

func NewOrQuery() *OrQuery {
	return &OrQuery{}
}

func (o *OrQuery) Next() (document.DocId, error) {
	panic("implement me")
}

func (o *OrQuery) GetGE(id document.DocId) (document.DocId, error) {
	panic("implement me")
}

func (t *OrQuery) String() string {
	return ""
}
