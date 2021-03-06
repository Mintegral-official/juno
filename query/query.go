package query

import (
	"github.com/Mintegral-official/juno/debug"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
)

type Query interface {
	Next()
	Current() (document.DocId, error)
	GetGE(id document.DocId) (document.DocId, error)
	Marshal() map[string]interface{}
	Unmarshal(idx *index.Indexer, res map[string]interface{}, e operation.Operation) Query
	DebugInfo() *debug.Debug
	SetDebug(level int)
}
