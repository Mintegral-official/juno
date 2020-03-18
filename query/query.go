package query

import (
	"github.com/MintegralTech/juno/debug"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/index"
)

type Query interface {
	Next()
	Current() (document.DocId, error)
	GetGE(id document.DocId) (document.DocId, error)
	Marshal() map[string]interface{}
	Unmarshal(idx index.Index, res map[string]interface{}) Query
	DebugInfo() *debug.Debug
	SetDebug(level int)
	SetLabel(label string)
}
