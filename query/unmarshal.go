package query

import (
	"github.com/MintegralTech/juno/index"
	"github.com/MintegralTech/juno/operation"
)

type UnmarshalQuery struct {
}

func (u *UnmarshalQuery) Unmarshal(idx *index.Indexer, res map[string]interface{}, e operation.Operation) Query {
	if _, ok := res["and"]; ok {
		var andQuery = &AndQuery{}
		return andQuery.Unmarshal(idx, res, e)
	} else if _, ok := res["or"]; ok {
		var orQuery = &OrQuery{}
		return orQuery.Unmarshal(idx, res, e)
	} else if _, ok := res["not"]; ok {
		var notAndQuery = &NotAndQuery{}
		return notAndQuery.Unmarshal(idx, res, e)
	} else if _, ok := res["="]; ok {
		var termQuery = &TermQuery{}
		return termQuery.Unmarshal(idx, res, e)
	}
	return nil
}
