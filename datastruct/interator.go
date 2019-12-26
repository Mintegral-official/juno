package datastruct

import "github.com/Mintegral-official/juno/document"

type Iterator interface {
	HasNext() bool
	Next()
	Current() interface{}
	GetGE(id document.DocId) interface{}
}
