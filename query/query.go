package query

import "github.com/Mintegral-official/juno/document"

type Query interface {
	HasNext() bool
	Next() document.DocId
}
