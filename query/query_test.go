package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"testing"
)

func TestQuery(t *testing.T) {
	var doc1 = &document.DocInfo{
		Id: 0,
		Fields: []*document.Field{
			{
				Name:      "field1",
				IndexType: 1,
				Value:     1,
			},
			{
				Name:      "field2",
				IndexType: 0,
				Value:     2,
			},
			{
				Name:      "field1",
				IndexType: 2,
				Value:     1,
			},
		},
	}
	var doc2 = &document.DocInfo{
		Id: 1,
		Fields: []*document.Field{
			{
				Name:      "field1",
				IndexType: 0,
				Value:     1,
			},
			{
				Name:      "field2",
				IndexType: 1,
				Value:     2,
			},
			{
				Name:      "field1",
				IndexType: 0,
				Value:     1,
			},
		},
	}
	var doc3 = &document.DocInfo{
		Id: 2,
		Fields: []*document.Field{
			{
				Name:      "field1",
				IndexType: 0,
				Value:     1,
			},
			{
				Name:      "field2",
				IndexType: 0,
				Value:     2,
			},
			{
				Name:      "field1",
				IndexType: 1,
				Value:     1,
			},
		},
	}
	idx := index.NewIndex("index")
	_ = idx.Add(doc1)
	_ = idx.Add(doc2)
	_ = idx.Add(doc3)
    q := NewNotAndQuery([]Query{
    	NewTermQuery(idx.GetStorageIndex().Iterator("field1").(*datastruct.SkipListIterator)),
    	NewTermQuery(idx.GetStorageIndex().Iterator("field2").(*datastruct.SkipListIterator)),
    	NewTermQuery(nil),
	}, nil)

    if cur, err := q.Current(); err != nil {
    	fmt.Println(cur)
	} else {
		fmt.Println(err)
	}
	for cur, err := q.Next(); err != nil; {
		fmt.Println(cur)
	}
}
