package index

import (
	"github.com/Mintegral-official/juno/query"
)

type Index interface {
	Add(doc *DocInfo) error
	Del(doc *DocInfo) error
	Update(doc *DocInfo) error

	Dump(filename string)
	Load(filename string)

	Search(query *query.Query) *query.SearchResult
}
