package query

import "github.com/Mintegral-official/juno/document"

type NotAndQuery struct {
}

func NewNotAndQuery() *NotAndQuery {
	return &NotAndQuery{}
}

func (n NotAndQuery) Next() (document.DocId, error) {
	//	panic("implement me")
	return 0, nil
}

func (n NotAndQuery) GetGE(id document.DocId) (document.DocId, error) {
	// panic("implement me")
	return 0, nil
}

func (n NotAndQuery) String() string {
	panic("implement me")
}
