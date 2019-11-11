package query

import (
	"github.com/Mintegral-official/juno/document"
)

type OrQuery struct {
	querys []Query
	curIdx int
}

func NewOrQuery(querys ...Query) *OrQuery {
	if querys == nil {
		return nil
	}
	return &OrQuery{
		querys: querys,
	}
}

func (o *OrQuery) Next() (document.DocId, error) {
	return 0, nil
}

func (o *OrQuery) GetGE(id document.DocId) (document.DocId, error) {
	return 0, nil
}

func (t *OrQuery) String() string {
	return ""
}
