package query

import (
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
)

type Query interface {
	// for query
	Next()
	Current() (document.DocId, error)
	GetGE(id document.DocId) (document.DocId, error)

	// for debug
	SetDebug(isDebug int) // level: 1 normal(走到末尾的链添或者checker失败加debug信息)  2 detail 所有的都要加Debug (暂不实现)

	Marshal() map[string]interface{}
	Unmarshal(idx *index.Indexer, res map[string]interface{}, e operation.Operation) Query
}
