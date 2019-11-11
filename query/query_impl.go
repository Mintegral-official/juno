package query

import "github.com/Mintegral-official/juno/document"

type QueryImpl struct {
}

func (q QueryImpl) Next() (document.DocId, error) {
	return 0, nil
}

func (q QueryImpl) String() string {
	return ""
}

func NewQueryImpl() *QueryImpl {
	return &QueryImpl{}
}
