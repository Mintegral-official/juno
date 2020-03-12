package search

import (
	"fmt"
	"github.com/MintegralTech/juno/check"
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/index"
	"github.com/MintegralTech/juno/query"
	"strconv"
	"time"
)

type Searcher struct {
	Docs       []document.DocId
	Time       time.Duration
	FilterInfo map[document.DocId]map[string]interface{}
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
	id, err := query.Current()
	query.Next()
	for err == nil {
		if vb, ok := iIndexer.GetBitMap().Get(index.DocId(id)); ok {
			s.Docs = append(s.Docs, vb.(document.DocId))
		}
		id, err = query.Current()
		query.Next()
	}
	s.Time = time.Since(now)
	s.IndexDebug = iIndexer.DebugInfo()
	s.QueryDebug = query.DebugInfo()
}

func (s *Searcher) Debug(idx *index.Indexer, q map[string]interface{}, ids []document.DocId) {
	uq := query.Unmarshal{}
	s.Search(idx, uq.Unmarshal(idx, q).(query.Query))
	queryMarshal := q
	var res = make(map[document.DocId]map[string]interface{}, len(ids))
	for _, id := range ids {
		tmp, ok := idx.GetCampaignMap().Get(index.DocId(id))
		if !ok {
			continue
		}
		for k, v := range queryMarshal {
			switch k {
			case "and", "or", "not", "and_check", "or_check", "not_and_check":
				debugInfo(v, idx, tmp.(document.DocId))
			case "=":
				var termQuery = &query.TermQuery{}
				if res, err := termQuery.Unmarshal(idx, map[string]interface{}{k: v.([]string)}).GetGE(tmp.(document.DocId));
					err != nil || res != id {
					queryMarshal[k] = append(v.([]string), "id not found")
				} else if res == id {
					queryMarshal[k] = append(v.([]string), "id found")
				}
			case "check":
				var c = &check.CheckerImpl{}
				queryMarshal[k] = append(v.([]interface{}), fmt.Sprintf("check result %s",
					strconv.FormatBool(c.Unmarshal(idx, map[string]interface{}{k: v}).Check(tmp.(document.DocId)))))
			case "in_check":
				var c = &check.InChecker{}
				queryMarshal[k] = append(v.([]interface{}), fmt.Sprintf("check result %s",
					strconv.FormatBool(c.Unmarshal(idx, map[string]interface{}{k: v}).Check(tmp.(document.DocId)))))
			case "not_check":
				var c = &check.NotChecker{}
				queryMarshal[k] = append(v.([]interface{}), fmt.Sprintf("check result %s",
					strconv.FormatBool(c.Unmarshal(idx, map[string]interface{}{k: v}).Check(tmp.(document.DocId)))))
			}
		}
		res[id] = queryMarshal
	}
	s.FilterInfo = res
}

func debugInfo(res interface{}, idx *index.Indexer, id document.DocId) {
	for _, value := range res.([]map[string]interface{}) {
		for k, v := range value {
			switch k {
			case "and", "or", "not", "and_check", "or_check", "not_and_check":
				debugInfo(v, idx, id)
			case "=":
				var termQuery = &query.TermQuery{}
				if res, err := termQuery.Unmarshal(idx, map[string]interface{}{k: v.([]string)}).GetGE(id);
					err != nil || res != id {
					value[k] = append(v.([]string), "id not found")
				} else if res == id {
					value[k] = append(v.([]string), "id found")
				}
			case "check":
				var chk = &check.CheckerImpl{}
				value[k] = append(v.([]interface{}), fmt.Sprintf("check result %s",
					strconv.FormatBool(chk.Unmarshal(idx, map[string]interface{}{k: v}).Check(id))))
			case "in_check":
				var chk = &check.InChecker{}
				value[k] = append(v.([]interface{}), fmt.Sprintf("check result %s",
					strconv.FormatBool(chk.Unmarshal(idx, map[string]interface{}{k: v}).Check(id))))
			case "not_check":
				var chk = &check.CheckerImpl{}
				value[k] = append(v.([]interface{}), fmt.Sprintf("check result %s",
					strconv.FormatBool(chk.Unmarshal(idx, map[string]interface{}{k: v}).Check(id))))
			}
		}
	}
}
