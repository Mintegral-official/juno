package query

import (
	"github.com/Mintegral-official/juno/document"
)

type NotAndQuery struct {
	checkers []Checker
	querys   []Query
	curIdx   int
}

func NewNotAndQuery(querys []Query, checkers []Checker) *NotAndQuery {
	if querys == nil {
		return nil
	}

	return &NotAndQuery{
		checkers: checkers,
		querys:   querys,
	}
}

func (n *NotAndQuery) Next() (document.DocId, error) {
	return 0, nil
}

func (n *NotAndQuery) GetGE(id document.DocId) (document.DocId, error) {
	return 0, nil
}

func (n *NotAndQuery) Current() (document.DocId, error) {
	return 0, nil
}

func (n *NotAndQuery) String() string {
	//panic("implement me")
	return ""
}

func (n *NotAndQuery) check(id document.DocId) bool {
	return true
}
