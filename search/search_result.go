package search

import (
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/query"
	"time"
)

type Result struct {
	Docs       []document.DocId
	Time       time.Duration
	QueryDebug string
}

func NewResult() *Result {
	return &Result{
		Docs: []document.DocId{},
	}
}

func (r *Result) Search(iIndexer *index.Indexer, query query.Query) *Result {
	if query == nil {
		return nil
	}
	now := time.Now()
	if _, err := query.Current(); err != nil {

		return r
	}
	id, err := query.Next()
	for err == nil {
		if !iIndexer.GetBitMap().IsExist(uint64(iIndexer.GetCampaignMap()[id])) {
			continue
		}
		r.Docs = append(r.Docs, id)
		id, err = query.Next()
	}
	r.Time = time.Since(now)
	r.QueryDebug = query.String()
	return r
}
