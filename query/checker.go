package query

import "github.com/Mintegral-official/juno/document"

type Checker interface {
	Check(id document.DocId) bool
}
