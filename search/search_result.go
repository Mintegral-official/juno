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
	var (
		v  interface{}
		ok bool
	)
	now := time.Now()
	id, err := query.Next()
	for err == nil {
		v, ok = iIndexer.GetCampaignMap().Get(index.DocId(id))
		if ok && !iIndexer.GetBitMap().IsExist(v.(document.DocId)) {
			continue
		}
		s.Docs = append(s.Docs, id)
		id, err = query.Next()
	}
	s.Time = time.Since(now)
	s.IndexDebug = iIndexer.DebugInfo()
	s.QueryDebug = query.DebugInfo()
}
