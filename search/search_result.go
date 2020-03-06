package search

import (
	"fmt"
	"github.com/Mintegral-official/juno/check"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
	"github.com/Mintegral-official/juno/query"
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
		s.Docs = append(s.Docs, id)
		id, err = query.Current()
		query.Next()
	}
	s.Time = time.Since(now)
	s.IndexDebug = iIndexer.DebugInfo()
	s.QueryDebug = query.DebugInfo()
}

func (s *Searcher) Debug(idx *index.Indexer, q map[string]interface{}, e operation.Operation, ids []document.DocId) {
	uq := query.UnmarshalQuery{}
	s.Search(idx, uq.Unmarshal(idx, q, e))
	queryMarshal := q
	var res = make(map[document.DocId]map[string]interface{}, len(ids))
	for _, id := range ids {
		for k, v := range queryMarshal {
			switch k {
			case "and", "or", "not", "and_check", "or_check", "not_and_check":
				debugInfo(v, idx, id, e)
			case "=":
				var termQuery = &query.TermQuery{}
				if res, err := termQuery.Unmarshal(idx, map[string]interface{}{k: v.([]string)}, e).GetGE(id); err != nil {
					queryMarshal[k] = append(v.([]string), "id not found")
				} else if res == id {
					queryMarshal[k] = append(v.([]string), "id found")
				}
			case "check":
				var c = &check.CheckerImpl{}
				queryMarshal[k] = append(v.([]interface{}), fmt.Sprintf("check result %s",
					strconv.FormatBool(c.Unmarshal(idx, map[string]interface{}{k: v}, e).Check(id))))
			case "in_check":
				var c = &check.InChecker{}
				queryMarshal[k] = append(v.([]interface{}), fmt.Sprintf("check result %s",
					strconv.FormatBool(c.Unmarshal(idx, map[string]interface{}{k: v}, e).Check(id))))
			case "not_check":
				var c = &check.NotChecker{}
				queryMarshal[k] = append(v.([]interface{}), fmt.Sprintf("check result %s",
					strconv.FormatBool(c.Unmarshal(idx, map[string]interface{}{k: v}, e).Check(id))))
			}
		}
		res[id] = queryMarshal
	}
	s.FilterInfo = res
}

func debugInfo(res interface{}, idx *index.Indexer, id document.DocId, e operation.Operation) {
	for _, value := range res.([]map[string]interface{}) {
		for k, v := range value {
			switch k {
			case "and", "or", "not", "and_check", "or_check", "not_and_check":
				debugInfo(v, idx, id, e)
			case "=":
				var termQuery = &query.TermQuery{}
				if res, err := termQuery.Unmarshal(idx, map[string]interface{}{k: v.([]string)}, e).GetGE(id); err != nil {
					value[k] = append(v.([]string), "id not found")
				} else if res == id {
					value[k] = append(v.([]string), "id found")
				}
			case "check":
				var chk = &check.CheckerImpl{}
				value[k] = append(v.([]interface{}), fmt.Sprintf("check result %s",
					strconv.FormatBool(chk.Unmarshal(idx, map[string]interface{}{k: v}, e).Check(id))))
			case "in_check":
				var chk = &check.InChecker{}
				value[k] = append(v.([]interface{}), fmt.Sprintf("check result %s",
					strconv.FormatBool(chk.Unmarshal(idx, map[string]interface{}{k: v}, e).Check(id))))
			case "not_check":
				var chk = &check.CheckerImpl{}
				value[k] = append(v.([]interface{}), fmt.Sprintf("check result %s",
					strconv.FormatBool(chk.Unmarshal(idx, map[string]interface{}{k: v}, e).Check(id))))
			}
		}
	}
}
