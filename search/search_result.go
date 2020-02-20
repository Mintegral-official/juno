package search

import (
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/query"
	"time"
)

type Searcher struct {
	Docs       []document.DocId
	Time       time.Duration
	IndexDebug *debug.Debug
	QueryDebug *debug.Debug
}

func NewSearcher() *Searcher {
	return &Searcher{
		Docs: []document.DocId{},
	}
}

func (s *Searcher) Search(iIndexer *index.Indexer, query query.Query) {
	if query == nil {
		panic("the query should not be nil")
		return
	}
	now := time.Now()
	id, err := query.Next()
	for err == nil {
		if v, ok := iIndexer.GetCampaignMap().Get(index.DocId(id)); ok && !iIndexer.GetBitMap().IsExist(v.(document.DocId)) {
			continue
		}
		s.Docs = append(s.Docs, id)
		id, err = query.Next()
	}
	s.Time = time.Since(now)
	s.IndexDebug = iIndexer.DebugInfo()
	s.QueryDebug = query.DebugInfo()
}

func (s *Searcher) Debug(q query.Query, ids []document.DocId) map[document.DocId]error {
	res := make(map[document.DocId]error, len(ids))
	for _, id := range ids {
		res[id] = nil
	}
	id, err := q.Next()
	for _, v := range ids {
		if v == id {
			res[id] = err
		}
	}
	for err == nil {
		id, err = q.Next()
		for _, v := range ids {
			if v == id {
				res[id] = err
			}
		}
	}
	return res
}
