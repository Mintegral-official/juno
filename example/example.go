package main

import (
	"fmt"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/index"
	"github.com/MintegralTech/juno/query"
	"github.com/MintegralTech/juno/search"
)

func main() {
	idx := index.NewIndex("default")
	_ = idx.Add(&document.DocInfo{
		Id: 1,
		Fields: []*document.Field{
			{Name: "field1", IndexType: document.InvertedIndexType, Value: int64(1), ValueType: document.IntFieldType},
			{Name: "field2", IndexType: document.InvertedIndexType, Value: "abc", ValueType: document.StringFieldType},
		},
	})

	s := search.NewSearcher()
	s.Search(idx, query.NewTermQuery(idx.GetInvertedIndex().Iterator("field1", "1")))
	fmt.Println(s.Docs)
}
