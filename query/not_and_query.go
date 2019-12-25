package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/query/check"
)

type NotAndQuery struct {
	queries  []Query
	checkers []check.Checker
	curIdx   int
	aDebug   *debug.Debug
}

func NewNotAndQuery(queries []Query, checkers []check.Checker) *NotAndQuery {
	naq := &NotAndQuery{
		aDebug: &debug.Debug{
			Name: "NewNotAndQuery",
			Msg:  []string{},
		},
	}
	if len(queries) == 0 {
		return naq
	}
	naq.queries = queries
	naq.checkers = checkers
	return naq
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
			naq.aDebug.AddDebug(fmt.Sprintf("docID[%d] is not suitable", target))
		}
		for i := 1; i < len(naq.queries); i++ {
			cur, err := naq.queries[i].GetGE(target)
			if (helpers.Compare(target, cur) != 0 || err != nil) && i == len(naq.queries)-1 {
				_, _ = naq.queries[0].Next()
				for !naq.check(target) {
					naq.aDebug.AddDebug(fmt.Sprintf("docID[%d] is not suitable", target))
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
				naq.aDebug.AddDebug(fmt.Sprintf("docID[%d] is not suitable", target))
				target, err = naq.queries[0].Next()
			}
			return target, nil
		}
		for i := 1; i < len(naq.queries); i++ {
			if _, err := naq.queries[i].Current(); err != nil {
				for !naq.check(target) {
					naq.aDebug.AddDebug(fmt.Sprintf("docID[%d] is not suitable", target))
					target, err = naq.queries[0].Next()
				}
				return target, nil
			}
			cur, err := naq.queries[i].GetGE(target)
			if (helpers.Compare(target, cur) != 0 || err != nil) && i == len(naq.queries)-1 {
				if naq.check(target) {
					return target, nil
				}
				naq.aDebug.AddDebug(fmt.Sprintf("docID[%d] is not suitable", target))
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
			naq.aDebug.AddDebug(fmt.Sprintf("docID[%d] is not suitable", res))
		}
	}
	naq.aDebug.AddDebug(fmt.Sprintf("current data[%d] is not suitable", res))
	return 0, errors.New("current data is not suitable")
}

func (naq *NotAndQuery) String() string {

	for i := 0; i < len(naq.queries); i++ {
		naq.aDebug.AddDebug(naq.queries[i].String())
	}
	if res, err := json.Marshal(naq.aDebug); err == nil {
		return string(res)
	} else {
		return err.Error()
	}
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
