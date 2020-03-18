package check

import (
	"github.com/MintegralTech/juno/index"
)

type unmarshal struct {
}

func (u *unmarshal) Unmarshal(idx index.Index, res map[string]interface{}) Checker {
	if _, ok := res["check"]; ok {
		var checkImpl = &CheckerImpl{}
		return checkImpl.Unmarshal(idx, res)
	}
	if _, ok := res["in_check"]; ok {
		var inCheck = &InChecker{}
		return inCheck.Unmarshal(idx, res)
	}
	if _, ok := res["not_check"]; ok {
		var notCheck = &NotChecker{}
		return notCheck.Unmarshal(idx, res)
	}
	if _, ok := res["or_check"]; ok {
		var orCheck = &OrChecker{}
		return orCheck.Unmarshal(idx, res)
	}
	if _, ok := res["and_check"]; ok {
		var andCheck = &AndChecker{}
		return andCheck.Unmarshal(idx, res)
	}
	if _, ok := res["not_and_check"]; ok {
		var notAndCheck = &NotAndChecker{}
		return notAndCheck.Unmarshal(idx, res)
	}
	return nil
}
