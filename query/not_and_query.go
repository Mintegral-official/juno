package query

import "github.com/Mintegral-official/juno/document"

type NotAndQuery struct {
	*AndQuery
	checker []Checker
}

func NewNotAndQuery(querys []Query, checks []Checker) *NotAndQuery {
	if querys == nil {
		return nil
	}
	return &NotAndQuery{
		AndQuery: NewAndQuery(querys, checks),
		checker:  checks,
	}
}

func (n NotAndQuery) Next() (document.DocId, error) {
	//	panic("implement me")
	return 0, nil
}

func (n NotAndQuery) GetGE(id document.DocId) (document.DocId, error) {
	// panic("implement me")
	return 0, nil
}

func (n NotAndQuery) Current() (document.DocId, error) {
  return 0, nil
}

func (n NotAndQuery) String() string {
	panic("implement me")
}
