package search

import (
	"fmt"
	"github.com/MintegralTech/juno/check"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/index"
	"github.com/MintegralTech/juno/query"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type SearcherResult struct {
	Docs []document.DocId
	Time time.Duration
}

func Search(iIndexer index.Index, query query.Query) *SearcherResult {
	if query == nil {
		logrus.Warnf("query is nil")
		return nil
	}
	var s = &SearcherResult{}
	now := time.Now()
	id, err := query.Current()
	query.Next()
	for err == nil {
		if i, e := iIndexer.GetId(id); e == nil {
			s.Docs = append(s.Docs, i)
		}
		id, err = query.Current()
		query.Next()
	}
	s.Time = time.Since(now)
	return s
}

func Replay(idx index.Index, q map[string]interface{}, ids []document.DocId) map[document.DocId][]map[string]interface{} {
	uq, queryMarshal := query.Unmarshal{}, q
	var res = make(map[document.DocId][]map[string]interface{}, len(ids))
	for _, id := range ids {
		tmp, ok := idx.GetInnerId(id)
		if ok != nil {
			continue
		}
		var tmpRes []map[string]interface{}
		for k, v := range queryMarshal {
			switch k {
			case "and", "or", "not":
				value := v.([]map[string]interface{})
				q := uq.Unmarshal(idx, map[string]interface{}{k: value}).(query.Query)
				if label, ok := value[len(value)-1]["label"]; ok && label != "" {
					tmpRes = append(tmpRes, map[string]interface{}{"label": label})
					if res, err := q.GetGE(tmp); err != nil || res != id {
						tmpRes = append(tmpRes, map[string]interface{}{"result": "not match"})
					} else if res == id {
						tmpRes = append(tmpRes, map[string]interface{}{"result": "match"})
					}
				}
				replay(v, idx, tmp)
			case "and_check", "or_check", "not_and_check":
				replay(v, idx, tmp)
			case "=":
				if res, err := uq.Unmarshal(idx, map[string]interface{}{k: v.([]string)}).(query.Query).GetGE(tmp); err != nil || res != id {
					queryMarshal[k] = append(v.([]string), "id not found")
				} else if res == id {
					queryMarshal[k] = append(v.([]string), "id found")
				}
			case "check":
				queryMarshal[k] = append(v.([]interface{}), fmt.Sprintf("check result %s",
					strconv.FormatBool(uq.Unmarshal(idx, map[string]interface{}{k: v}).(check.Checker).Check(tmp))))
			case "in_check":
				queryMarshal[k] = append(v.([]interface{}), fmt.Sprintf("check result %s",
					strconv.FormatBool(uq.Unmarshal(idx, map[string]interface{}{k: v}).(check.Checker).Check(tmp))))
			case "not_check":
				queryMarshal[k] = append(v.([]interface{}), fmt.Sprintf("check result %s",
					strconv.FormatBool(uq.Unmarshal(idx, map[string]interface{}{k: v}).(check.Checker).Check(tmp))))
			}
		}
		tmpRes = append(tmpRes, map[string]interface{}{"detail": queryMarshal})
		res[id] = tmpRes
		tmpRes = []map[string]interface{}{}
	}
	return res
}

func replay(res interface{}, idx index.Index, id document.DocId) {
	uq := &query.Unmarshal{}
	for i, value := range res.([]map[string]interface{}) {
		for k, v := range value {
			switch k {
			case "and", "or", "not":
				var tmpRes []map[string]interface{}
				value := v.([]map[string]interface{})
				q := uq.Unmarshal(idx, map[string]interface{}{k: value}).(query.Query)
				if label, ok := value[len(value)-1]["label"]; ok && label != "" {
					tmpRes = append(tmpRes, map[string]interface{}{"label": label})
					if res, err := q.GetGE(id); err != nil || res != id {
						tmpRes = append(tmpRes, map[string]interface{}{"result": "not match"})
					} else if res == id {
						tmpRes = append(tmpRes, map[string]interface{}{"result": "match"})
					}
				}
				tmpRes = append(tmpRes, map[string]interface{}{"detail": map[string]interface{}{k: value}})
				res.([]map[string]interface{})[i] = map[string]interface{}{k: tmpRes}
				replay(v, idx, id)
			case "and_check", "or_check", "not_and_check":
				replay(v, idx, id)
			case "=":
				if res, err := uq.Unmarshal(idx, map[string]interface{}{k: v.([]string)}).(query.Query).GetGE(id); err != nil || res != id {
					value[k] = append(v.([]string), "id not found")
				} else if res == id {
					value[k] = append(v.([]string), "id found")
				}
			case "check":
				value[k] = append(v.([]interface{}), fmt.Sprintf("check result %s",
					strconv.FormatBool(uq.Unmarshal(idx, map[string]interface{}{k: v}).(check.Checker).Check(id))))
			case "in_check":
				value[k] = append(v.([]interface{}), fmt.Sprintf("check result %s",
					strconv.FormatBool(uq.Unmarshal(idx, map[string]interface{}{k: v}).(check.Checker).Check(id))))
			case "not_check":
				value[k] = append(v.([]interface{}), fmt.Sprintf("check result %s",
					strconv.FormatBool(uq.Unmarshal(idx, map[string]interface{}{k: v}).(check.Checker).Check(id))))
			}
		}
	}
}
