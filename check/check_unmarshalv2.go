package check

import (
	"github.com/MintegralTech/juno/index"
	"github.com/MintegralTech/juno/marshal"
)

type unmarshalV2 struct {
}

func (u *unmarshalV2) UnmarshalV2(idx index.Index, marshalInfo *marshal.MarshalInfo) Checker {
	if marshalInfo.Operation == "in_check" {
		var inCheck = &InChecker{}
		return inCheck.UnmarshalV2(idx, marshalInfo)
	}
	if marshalInfo.Operation == "not_check" {
		var notCheck = &NotChecker{}
		return notCheck.UnmarshalV2(idx, marshalInfo)
	}
	if marshalInfo.Operation == "or_check" {
		var orCheck = &OrChecker{}
		return orCheck.UnmarshalV2(idx, marshalInfo)
	}
	if marshalInfo.Operation == "and_check" {
		var andCheck = &AndChecker{}
		return andCheck.UnmarshalV2(idx, marshalInfo)
	}
	if marshalInfo.Operation == "not_and_check" {
		var notAndCheck = &NotAndChecker{}
		return notAndCheck.UnmarshalV2(idx, marshalInfo)
	}
	if marshalInfo.Op >= 0 {
		var checkImpl = &CheckerImpl{}
		return checkImpl.UnmarshalV2(idx, marshalInfo)
	}
	return nil
}
