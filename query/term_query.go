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

func (t *TermQuery) HasNext() bool {
	panic("implement me")
}

func (t *TermQuery) Next() document.DocId {
	panic("implement me")
}
