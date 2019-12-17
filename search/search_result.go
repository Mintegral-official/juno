package search

import (
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/query"
	"time"
)

type Result struct {
	Docs []document.DocId
	Time time.Duration
}

func Search(ii *index.IndexImpl, query query.Query) *Result {
	if query == nil {
		return nil
	}
	s, now := &Result{Docs: []document.DocId{}}, time.Now()
	if _, err := query.Current(); err != nil {
		return s
	}
	id, err := query.Next()
	for err == nil {
		if !ii.GetBitMap().IsExist(uint64(ii.GetCampaignMap()[id])) {
			continue
		}
		s.Docs = append(s.Docs, id)
		id, err = query.Next()
	}
	s.Time = time.Since(now)
	return s
}
