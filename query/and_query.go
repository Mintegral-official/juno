package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/document"
	"github.com/pkg/errors"
)

type AndQuery struct {
	querys []Query
	curIdx int
}

func NewAndQuery(querys ...Query) *AndQuery {
	if querys == nil {
		return nil
	}
	return &AndQuery{
		querys: querys,
	}
}

func (a *AndQuery) Next() (document.DocId, error) {
	lastIdx := a.curIdx
	curIdx := a.curIdx
	target, err := a.querys[curIdx].Next()
	if err != nil {
		return 0, errors.Wrap(err, "no more data")
	}
	for {
		curIdx = (curIdx + 1) % len(a.querys)
		cur, err := a.querys[curIdx].GetGE(target)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(target), curIdx))
		}
		if cur != target {
			lastIdx = curIdx
			target = cur
		}
		if (curIdx+1)%len(a.querys) == lastIdx {
			return target, nil
		}
	}
}

func (a *AndQuery) GetGE(id document.DocId) (document.DocId, error) {
	curIdx := a.curIdx
	res, err := a.querys[a.curIdx].GetGE(id)
	// fmt.Println(err)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(res), curIdx))
	}
	curIdx++
	for curIdx < len(a.querys) {
		cur, err := a.querys[a.curIdx].GetGE(id)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("not find [%d] in querys[%d]", int64(cur), curIdx))
		}
		if cur <= res {
			res = cur
		}
		curIdx++
	}
    return res, nil
}

func (t *AndQuery) String() string {
	return ""
}
