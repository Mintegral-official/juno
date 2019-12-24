package query

import (
	"errors"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/query/check"
)

type NotAndQuery struct {
	queries  []Query
	checkers []check.Checker
	curIdx   int
}

func NewNotAndQuery(queries []Query, checkers []check.Checker) *NotAndQuery {
	if len(queries) == 0 {
		return &NotAndQuery{}
	}
	return &NotAndQuery{
		checkers: checkers,
		queries:  queries,
	}
}

func (naq *NotAndQuery) Next() (document.DocId, error) {
	for {
		target, err := naq.queries[0].Current()
		if err != nil {
			return 0, helpers.NoMoreData
		}
		if len(naq.queries) == 1 {
			_, _ = naq.queries[0].Next()
			if naq.check(target) {
				return target, nil
			}
		}
		for i := 1; i < len(naq.queries); i++ {
			cur, err := naq.queries[i].GetGE(target)
			if (helpers.Compare(target, cur) != 0 || err != nil) && i == len(naq.queries)-1 {
				_, _ = naq.queries[0].Next()
				for !naq.check(target) {
					target, err = naq.queries[0].Current()
					if err != nil {
						return 0, err
					}
					_, _ = naq.queries[0].Next()
				}
				return target, nil
			}
		}
		target, err = naq.queries[0].Next()
	}
}

func (naq *NotAndQuery) GetGE(id document.DocId) (document.DocId, error) {
	for {
		target, err := naq.queries[0].GetGE(id)
		if err != nil {
			return 0, helpers.NoMoreData
		}
		if len(naq.queries) == 1 {
			for !naq.check(target) {
				target, err = naq.queries[0].Next()
			}
			return target, nil
		}
		for i := 1; i < len(naq.queries); i++ {
			if _, err := naq.queries[i].Current(); err != nil {
				for !naq.check(target) {
					target, err = naq.queries[0].Next()
				}
				return target, nil
			}
			cur, err := naq.queries[i].GetGE(target)
			if (helpers.Compare(target, cur) != 0 || err != nil) && i == len(naq.queries)-1 {
				if naq.check(target) {
					return target, nil
				}
			}
		}
		_, _ = naq.queries[0].Next()
	}
}

func (naq *NotAndQuery) Current() (document.DocId, error) {
	res, err := naq.queries[0].Current()
	if err != nil {
		return 0, err
	}
	for i := 1; i < len(naq.queries); i++ {
		tar, err := naq.queries[i].GetGE(res)
		_, _ = naq.queries[0].Next()
		if err != nil {
			continue
		}
		if tar == res {
			return 0, errors.New("this target is not result")
		} else if i == len(naq.queries)-1 {
			if naq.check(res) {
				return res, nil
			}
		}
	}
	return 0, errors.New("current data is not filter")
}

func (naq *NotAndQuery) String() string {
	//panic("implement me")
	return ""
}

func (naq *NotAndQuery) check(id document.DocId) bool {
	if naq.checkers == nil {
		return true
	}
	for i := 1; i < len(naq.checkers); i++ {
		if naq.checkers[i].Check(id) {
			return false
		}
	}
	return true
}
