package index

import "github.com/Mintegral-official/juno/document"

type InvertedIterator interface {
	HasNext() bool
	Next() document.DocId
	GetGE(id document.DocId) document.DocId
}
