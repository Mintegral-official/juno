package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/document"
)
import "github.com/pkg/errors"

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
	panic("implement me")
}

func (t *AndQuery) String() string {
	return ""
}
