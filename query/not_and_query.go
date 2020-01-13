package query

import (
	"errors"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/query/check"
	"strconv"
)

type NotAndQuery struct {
	queries  []Query
	checkers []check.Checker
	curIdx   int
	debugs   *debug.Debugs
}

func NewNotAndQuery(queries []Query, checkers []check.Checker, isDebug ...int) *NotAndQuery {
	naq := &NotAndQuery{}
	if len(isDebug) == 1 && isDebug[0] == 1 {
		naq.debugs = debug.NewDebugs(debug.NewDebug("NotAndQuery"))
	}
	if len(queries) == 0 {
		return naq
	}
	naq.queries = queries
	naq.checkers = checkers
	return naq
}

func (naq *NotAndQuery) Next() (document.DocId, error) {
	if naq.debugs != nil {
		naq.debugs.NextNum++
	}
	for {
		target, err := naq.queries[0].Current()
		if err != nil {
			return 0, helpers.NoMoreData
		}
		if len(naq.queries) == 1 {
			_, _ = naq.queries[0].Next()
			if target != 0 && naq.check(target) {
				return target, nil
			}
			if naq.debugs != nil {
				naq.debugs.DebugInfo.AddDebugMsg(strconv.FormatInt(int64(target), 10) + "has been filtered out")
			}
		}
		for i := 1; i < len(naq.queries); i++ {
			cur, err := naq.queries[i].GetGE(target)
			if (target != cur || err != nil) && i == len(naq.queries)-1 {
				_, _ = naq.queries[0].Next()
				for !naq.check(target) {
					if naq.debugs != nil {
						naq.debugs.DebugInfo.AddDebugMsg(strconv.FormatInt(int64(target), 10) + "has been filtered out")
					}
					target, err = naq.queries[0].Current()
					if err != nil {
						return 0, err
					}
					_, _ = naq.queries[0].Next()
				}
				if target != 0 {
					return target, nil
				}
			}
		}
		target, err = naq.queries[0].Next()
	}
}

func (naq *NotAndQuery) GetGE(id document.DocId) (document.DocId, error) {
	if naq.debugs != nil {
		naq.debugs.GetNum++
	}
	for {
		target, err := naq.queries[0].GetGE(id)
		if err != nil {
			return 0, helpers.NoMoreData
		}
		if len(naq.queries) == 1 {
			for !naq.check(target) {
				if naq.debugs != nil {
					naq.debugs.DebugInfo.AddDebugMsg(strconv.FormatInt(int64(target), 10) + "has been filtered out")
				}
				target, err = naq.queries[0].Next()
			}
			return target, nil
		}
		for i := 1; i < len(naq.queries); i++ {
			if _, err := naq.queries[i].Current(); err != nil {
				for !naq.check(target) {
					if naq.debugs != nil {
						naq.debugs.DebugInfo.AddDebugMsg(strconv.FormatInt(int64(target), 10) + "has been filtered out")
					}
					target, err = naq.queries[0].Next()
				}
				return target, nil
			}
			cur, err := naq.queries[i].GetGE(target)
			if (target != cur || err != nil) && i == len(naq.queries)-1 {
				if target != 0 && naq.check(target) {
					return target, nil
				}
				if naq.debugs != nil {
					naq.debugs.DebugInfo.AddDebugMsg(strconv.FormatInt(int64(target), 10) + "has been filtered out")
				}
			}
		}
		_, _ = naq.queries[0].Next()
	}
}

func (naq *NotAndQuery) Current() (document.DocId, error) {
	if naq.debugs != nil {
		naq.debugs.CurNum++
	}
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
			if naq.debugs != nil {
				naq.debugs.DebugInfo.AddDebugMsg(strconv.FormatInt(int64(res), 10) + "has been filtered out")
			}
		}
	}
	if naq.debugs != nil {
		naq.debugs.DebugInfo.AddDebugMsg(strconv.FormatInt(int64(res), 10) + "has been filtered out")
	}
	return 0, nil
}

func (naq *NotAndQuery) DebugInfo() *debug.Debug {
	if naq.debugs != nil {
		naq.debugs.DebugInfo.AddDebugMsg("next has been called: " + strconv.Itoa(naq.debugs.NextNum))
		naq.debugs.DebugInfo.AddDebugMsg("get has been called: " + strconv.Itoa(naq.debugs.GetNum))
		naq.debugs.DebugInfo.AddDebugMsg("current has been called: " + strconv.Itoa(naq.debugs.CurNum))
		for i := 0; i < len(naq.queries); i++ {
			naq.debugs.DebugInfo.AddDebug(naq.queries[i].DebugInfo())
		}
		return naq.debugs.DebugInfo
	}
	return nil
}

func (naq *NotAndQuery) check(id document.DocId) bool {
	if len(naq.checkers) == 0 {
		return true
	}
	for _, v := range naq.checkers {
		if v == nil {
			continue
		}
		if v.Check(id) {
			return false
		}
	}
	return true
}
