package datastruct

import "github.com/MintegralTech/juno/document"

type Iterator interface {
	HasNext() bool
	Next()
	Current() *Element
	GetGE(id document.DocId) *Element
}
