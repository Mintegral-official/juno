package query

import (
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
)

type NotAndQuery struct {
	querys   []Query
	checkers []Checker
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
	for {
		target, err := n.querys[0].Current()
		if err != nil {
			return 0, helpers.NoMoreData
		}
		if len(n.querys) == 1 {
			_, _ = n.querys[0].Next()
			if n.check(target) {
				return target, nil
			}
		}
		for i := 1; i < len(n.querys); i++ {
			cur, err := n.querys[i].GetGE(target)
			if (helpers.Compare(target, cur) != 0 || err != nil) && i == len(n.querys)-1 {
				_, _ = n.querys[0].Next()
				if n.check(target) {
					return target, nil
				}
			}
		}
		target, err = n.querys[0].Next()
	}
}

func (n *NotAndQuery) GetGE(id document.DocId) (document.DocId, error) {
	for {
		target, err := n.querys[0].GetGE(id)
		if err != nil {
			return 0, helpers.NoMoreData
		}
		if len(n.querys) == 1 {
			for !n.check(target) {
				target, err = n.querys[0].Next()
			}
			return target, nil
		}
		for i := 1; i < len(n.querys); i++ {
			if _, err := n.querys[i].Current(); err != nil {
				for !n.check(target) {
					target, err = n.querys[0].Next()
				}
				return target, nil
			}
			cur, err := n.querys[i].GetGE(target)
			if (helpers.Compare(target, cur) != 0 || err != nil) && i == len(n.querys)-1 {
				if n.check(target) {
					return target, nil
				}
			}
		}
		_, _ = n.querys[0].Next()
	}
}

func (n *NotAndQuery) Current() (document.DocId, error) {
	res, err := n.querys[0].Current()
	if err != nil {
		return 0, err
	}
	for i := 1; i < len(n.querys); i++ {
		tar, err := n.querys[i].Current()
		if err != nil {
			return 0, err
		}
		if tar != res {
			return res, nil
		}
	}
	return res, err
}

func (n *NotAndQuery) String() string {
	//panic("implement me")
	return ""
}

func (n *NotAndQuery) check(id document.DocId) bool {
	return true
}
