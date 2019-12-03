package index

import (
	"github.com/Mintegral-official/juno/document"
	"time"
)

type SearchResult struct {
	Docs []document.DocId
	Time time.Duration
}
