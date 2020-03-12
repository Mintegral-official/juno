package query

import (
	"github.com/MintegralTech/juno/check"
	"github.com/MintegralTech/juno/index"
)

type Unmarshal struct {
}

func (u *Unmarshal) Unmarshal(idx index.Index, res map[string]interface{}) interface{} {
	if _, ok := res["and"]; ok {
		var andQuery = &AndQuery{}
		return andQuery.Unmarshal(idx, res)
	}
	if _, ok := res["or"]; ok {
		var orQuery = &OrQuery{}
		return orQuery.Unmarshal(idx, res)
	}
	if _, ok := res["not"]; ok {
		var notAndQuery = &NotAndQuery{}
		return notAndQuery.Unmarshal(idx, res)
	}
	if _, ok := res["="]; ok {
		var termQuery = &TermQuery{}
		return termQuery.Unmarshal(idx, res)
	}
	if _, ok := res["=_check"]; ok {
		var checkImpl = &check.CheckerImpl{}
		return checkImpl.Unmarshal(idx, res)
	}
	if _, ok := res["in_check"]; ok {
		var inCheck = &check.InChecker{}
		return inCheck.Unmarshal(idx, res)
	}
	if _, ok := res["not_check"]; ok {
		var notCheck = &check.NotChecker{}
		return notCheck.Unmarshal(idx, res)
	}
	if _, ok := res["or_check"]; ok {
		var orCheck = &check.OrChecker{}
		return orCheck.Unmarshal(idx, res)
	}
	if _, ok := res["and_check"]; ok {
		var andCheck = &check.AndChecker{}
		return andCheck.Unmarshal(idx, res)
	}
	if _, ok := res["not_and_check"]; ok {
		var notAndCheck = &check.NotAndChecker{}
		return notAndCheck.Unmarshal(idx, res)
	}
	return nil
}
