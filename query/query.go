package query

import (
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
)

type Query interface {
	Next() (document.DocId, error)
	Current() (document.DocId, error)
	GetGE(id document.DocId) (document.DocId, error)
	DebugInfo() *debug.Debug
}
