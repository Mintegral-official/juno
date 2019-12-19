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

func Search(iIndexer *index.Indexer, query query.Query) *Result {
	if query == nil {
		return nil
	}
	s, now := &Result{Docs: []document.DocId{}}, time.Now()
	if _, err := query.Current(); err != nil {
		return s
	}
	id, err := query.Next()
	for err == nil {
		if !iIndexer.GetBitMap().IsExist(uint64(iIndexer.GetCampaignMap()[id])) {
			continue
		}
		s.Docs = append(s.Docs, id)
		id, err = query.Next()
	}
	s.Time = time.Since(now)
	return s
}
