package query

import "github.com/Mintegral-official/juno/document"

type QueryImpl struct {
}

func (q QueryImpl) HasNext() bool {
	return false
}

func (q QueryImpl) Next() document.DocId {
	return InvalidDocid
}

func (q QueryImpl) String() string {
	return ""
}

func NewQueryImpl() Query {
	return &QueryImpl{}
}
