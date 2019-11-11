package query

import "github.com/Mintegral-official/juno/document"

type Query interface {
	Next() (document.DocId, error)
	GetGE(id document.DocId) (document.DocId, error)
	String() string
}
