package datastruct

import "github.com/Mintegral-official/juno/document"

type Iterator interface {
	HasNext() bool
	Next()
	Current() *Element
	GetGE(id document.DocId) *Element
}
