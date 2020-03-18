package search

import (
	"fmt"
	"github.com/MintegralTech/juno/check"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/index"
	"github.com/MintegralTech/juno/marshal"
	"github.com/MintegralTech/juno/query"
	"github.com/sirupsen/logrus"
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

func Replay(idx index.Index, info *marshal.MarshalInfo, ids []document.DocId) map[document.DocId]*marshal.MarshalInfo {
	if info == nil {
		return nil
	}
	uq, marshalInfo := query.UnmarshalV2{}, info
	var res = make(map[document.DocId]*marshal.MarshalInfo, len(ids))
	for _, id := range ids {
		innerId, ok := idx.GetInnerId(id)
		if ok != nil {
			res[id] = nil
			continue
		}
		switch marshalInfo.Operation {
		case "and", "or", "not":
			q := uq.UnmarshalV2(idx, marshalInfo).(query.Query)
			if resId, err := q.GetGE(innerId); err != nil || resId != innerId {
				marshalInfo.Result = false
			} else if resId == innerId {
				marshalInfo.Result = true
			}
			for _, v := range marshalInfo.Nodes {
				if v != nil {
					Replay(idx, v, []document.DocId{id})
				}
			}
		case "=":
			q := uq.UnmarshalV2(idx, marshalInfo).(query.Query)
			marshalInfo.IndexValue = idx.GetIndexDebugInfoById(id).InvertIndex[marshalInfo.Name]
			if resId, err := q.GetGE(innerId); err != nil || resId != innerId {
				marshalInfo.Result = false
			} else if resId == innerId {
				marshalInfo.Result = true
			}
		case "and_check", "or_check", "not_and_check":
			c := uq.UnmarshalV2(idx, marshalInfo).(check.Checker)
			marshalInfo.Result = c.Check(innerId)
			for _, v := range marshalInfo.Nodes {
				if v != nil {
					Replay(idx, v, []document.DocId{id})
				}
			}
		case "in_check", "not_check":
			c := uq.UnmarshalV2(idx, marshalInfo).(check.Checker)
			marshalInfo.Result = c.Check(innerId)
			marshalInfo.IndexValue = idx.GetIndexDebugInfoById(id).StorageIndex[marshalInfo.Name]
		default:
			if marshalInfo.Operation != "" {
				fmt.Println(marshalInfo)
				c := uq.UnmarshalV2(idx, marshalInfo).(check.Checker)
				marshalInfo.Result = c.Check(innerId)
				marshalInfo.IndexValue = idx.GetIndexDebugInfoById(id).StorageIndex[marshalInfo.Name]
			}
		}
		res[id] = marshalInfo
	}
	return res
}
