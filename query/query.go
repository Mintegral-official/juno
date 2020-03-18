package query

import (
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/index"
	"github.com/MintegralTech/juno/marshal"
)

type Query interface {
	Next()
	Current() (document.DocId, error)
	GetGE(id document.DocId) (document.DocId, error)

	Marshal() map[string]interface{}
	Unmarshal(idx index.Index, res map[string]interface{}) Query
	MarshalV2() *marshal.MarshalInfo
	UnmarshalV2(idx index.Index, info *marshal.MarshalInfo) Query
	DebugInfo() *debug.Debug
	SetDebug(level int)
	SetLabel(label string)
}
