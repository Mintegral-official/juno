package index

import (
	"github.com/Mintegral-official/juno/query"
)

type Index interface {
	Add(doc *DocInfo)
	Del(doc *DocInfo)
	UpDate(doc *DocInfo)

	Dump(filename string)
	Load()

	Search(query *query.Query)
}
