package query

import (
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
)

type TermQuery struct {
	index.InvertedIterator
}

func NewTermQuery() *TermQuery {
	return &TermQuery{}
}

func (t *TermQuery) Next() (document.DocId, error) {
	panic("implement me")
}

func (t *TermQuery) GetGE(id document.DocId) (document.DocId, error) {
	panic("implement me")
}

func (t *TermQuery) String() string {
	return ""
}
