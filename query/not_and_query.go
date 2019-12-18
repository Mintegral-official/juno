package query

import (
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/query/check"
)

type NotAndQuery struct {
	querySlice []Query
	checkers   []check.Checker
	curIdx     int
}

func NewNotAndQuery(querys []Query, checkers []check.Checker) *NotAndQuery {
	if querys == nil {
		return nil
	}
	return &NotAndQuery{
		checkers:   checkers,
		querySlice: querys,
	}
}

func (naq *NotAndQuery) Next() (document.DocId, error) {
	for {
		target, err := naq.querySlice[0].Current()
		if err != nil {
			return 0, helpers.NoMoreData
		}
		if len(naq.querySlice) == 1 {
			_, _ = naq.querySlice[0].Next()
			if naq.check(target) {
				return target, nil
			}
		}
		for i := 1; i < len(naq.querySlice); i++ {
			cur, err := naq.querySlice[i].GetGE(target)
			if (helpers.Compare(target, cur) != 0 || err != nil) && i == len(naq.querySlice)-1 {
				_, _ = naq.querySlice[0].Next()
				if naq.check(target) {
					return target, nil
				}
			}
		}
		target, err = naq.querySlice[0].Next()
	}
}

func (naq *NotAndQuery) GetGE(id document.DocId) (document.DocId, error) {
	for {
		target, err := naq.querySlice[0].GetGE(id)
		if err != nil {
			return 0, helpers.NoMoreData
		}
		if len(naq.querySlice) == 1 {
			for !naq.check(target) {
				target, err = naq.querySlice[0].Next()
			}
			return target, nil
		}
		for i := 1; i < len(naq.querySlice); i++ {
			if _, err := naq.querySlice[i].Current(); err != nil {
				for !naq.check(target) {
					target, err = naq.querySlice[0].Next()
				}
				return target, nil
			}
			cur, err := naq.querySlice[i].GetGE(target)
			if (helpers.Compare(target, cur) != 0 || err != nil) && i == len(naq.querySlice)-1 {
				if naq.check(target) {
					return target, nil
				}
			}
		}
		_, _ = naq.querySlice[0].Next()
	}
}

func (naq *NotAndQuery) Current() (document.DocId, error) {
	res, err := naq.querySlice[0].Current()
	if err != nil {
		return 0, err
	}
	for i := 1; i < len(naq.querySlice); i++ {
		tar, err := naq.querySlice[i].Current()
		if err != nil {
			return 0, err
		}
		if tar != res {
			return res, nil
		}
	}
	return res, err
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
