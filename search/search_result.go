package search

import (
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/query"
	//	cmap "github.com/easierway/concurrent_map"
	"time"
)

type Searcher struct {
	Docs       []document.DocId
	Time       time.Duration
	IndexDebug string
	QueryDebug string
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
		if v, ok := iIndexer.GetCampaignMap().Load(id); ok && !iIndexer.GetBitMap().IsExist(v.(document.DocId)) {
			continue
		}
		s.Docs = append(s.Docs, id)
		id, err = query.Next()
	}
	s.Time = time.Since(now)
	s.IndexDebug = iIndexer.String()
	s.QueryDebug = query.String()
}
