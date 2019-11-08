package query

import "github.com/Mintegral-official/juno/document"

type RangeQuery struct {
}

func (r *RangeQuery) HasNext() bool {
	panic("implement me")
}

func (r *RangeQuery) Next() document.DocId {
	panic("implement me")
}

func (t *RangeQuery) String() string {
	return ""
}
