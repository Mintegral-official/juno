package query

import (
	"github.com/MintegralTech/juno/check"
	"github.com/MintegralTech/juno/index"
	"github.com/MintegralTech/juno/marshal"
)

type UnmarshalV2 struct {
}

func (u *UnmarshalV2) UnmarshalV2(idx index.Index, marshalInfo *marshal.MarshalInfo) interface{} {
	if marshalInfo == nil {
		return nil
	}

	if marshalInfo.Operation == "and" {
		var andQuery = &AndQuery{}
		return andQuery.UnmarshalV2(idx, marshalInfo)
	}
	if marshalInfo.Operation == "or" {
		var orQuery = &OrQuery{}
		return orQuery.UnmarshalV2(idx, marshalInfo)
	}
	if marshalInfo.Operation == "not" {
		var notAndQuery = &NotAndQuery{}
		return notAndQuery.UnmarshalV2(idx, marshalInfo)
	}
	if marshalInfo.Operation == "=" {
		var termQuery = &TermQuery{}
		return termQuery.UnmarshalV2(idx, marshalInfo)
	}
	if marshalInfo.Operation == "in_check" {
		var inCheck = &check.InChecker{}
		return inCheck.UnmarshalV2(idx, marshalInfo)
	}
	if marshalInfo.Operation == "not_check" {
		var notCheck = &check.NotChecker{}
		return notCheck.UnmarshalV2(idx, marshalInfo)
	}
	if marshalInfo.Operation == "or_check" {
		var orCheck = &check.OrChecker{}
		return orCheck.UnmarshalV2(idx, marshalInfo)
	}
	if marshalInfo.Operation == "and_check" {
		var andCheck = &check.AndChecker{}
		return andCheck.UnmarshalV2(idx, marshalInfo)
	}
	if marshalInfo.Operation == "not_and_check" {
		var notAndCheck = &check.NotAndChecker{}
		return notAndCheck.UnmarshalV2(idx, marshalInfo)
	}
	if marshalInfo.Op >= 0 {
		var checkImpl = &check.CheckerImpl{}
		return checkImpl.UnmarshalV2(idx, marshalInfo)
	}
	return nil
}
