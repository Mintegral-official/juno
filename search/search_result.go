package search

import (
	"fmt"
	"github.com/Mintegral-official/juno/check"
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
	"github.com/Mintegral-official/juno/query"
	"strconv"
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
	for err != helpers.NoMoreData {
		if err == nil {
			if v, ok := iIndexer.GetCampaignMap().Get(index.DocId(id)); ok && !iIndexer.GetBitMap().IsExist(v.(document.DocId)) {
				continue
			}
			s.Docs = append(s.Docs, id)
		}
		id, err = query.Next()
	}
	s.Time = time.Since(now)
	s.IndexDebug = iIndexer.DebugInfo()
	s.QueryDebug = query.DebugInfo()
}

func (s *Searcher) Debug(idx *index.Indexer, q query.Query, e operation.Operation, ids []document.DocId) map[document.DocId]map[string]interface{} {
	queryMarshal := q.Marshal()
	var res = make(map[document.DocId]map[string]interface{}, len(ids))
	for _, id := range ids {
		for k, v := range queryMarshal {
			switch k {
			case "and", "or", "not", "and_check", "or_check":
				debugInfo(v, idx, id, e)
			case "=":
				invertValue := idx.GetValueById(id)[0]
				if _, ok := invertValue[v.([]string)[0]]; !ok || len(invertValue[v.([]string)[0]]) == 0 {
					queryMarshal[k] = append(v.([]string), "id not found")
				}
				for i, iv := range invertValue[v.([]string)[0]] {
					if iv == v.([]string)[1] {
						queryMarshal[k] = append(v.([]string), "id found")
					} else if i == len(invertValue[v.([]string)[0]])-1 {
						queryMarshal[k] = append(v.([]string), "id not found")
					}
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
	return res
}

func debugInfo(res interface{}, idx *index.Indexer, id document.DocId, e operation.Operation) {
	for _, value := range res.([]map[string]interface{}) {
		for k, v := range value {
			switch k {
			case "and", "or", "not", "and_check", "or_check":
				debugInfo(v, idx, id, e)
			case "=":
				invertValue := idx.GetValueById(id)[0]
				if _, ok := invertValue[v.([]string)[0]]; !ok || len(invertValue[v.([]string)[0]]) == 0 {
					value[k] = append(v.([]string), "id not found")
				}
				for i, iv := range invertValue[v.([]string)[0]] {
					if iv == v.([]string)[1] {
						value[k] = append(v.([]string), "id found")
					} else if i == len(invertValue[v.([]string)[0]])-1 {
						value[k] = append(v.([]string), "id not found")
					}
				}
			case "check":
				var c = &check.CheckerImpl{}
				value[k] = append(v.([]interface{}), fmt.Sprintf("check result %s",
					strconv.FormatBool(c.Unmarshal(idx, map[string]interface{}{k: v}, e).Check(id))))
			case "in_check":
				var c = &check.InChecker{}
				value[k] = append(v.([]interface{}), fmt.Sprintf("check result %s",
					strconv.FormatBool(c.Unmarshal(idx, map[string]interface{}{k: v}, e).Check(id))))
			case "not_check":
				var c = &check.NotChecker{}
				value[k] = append(v.([]interface{}), fmt.Sprintf("check result %s",
					strconv.FormatBool(c.Unmarshal(idx, map[string]interface{}{k: v}, e).Check(id))))
			}
		}
	}
}
