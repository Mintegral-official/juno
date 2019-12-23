package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/index"
	"testing"
)

func TestQuery(t *testing.T) {
	idx := index.NewIndex("index")
	_ = idx.Add(doc1)
	_ = idx.Add(doc2)
	_ = idx.Add(doc3)
	storage := idx.GetStorageIndex()

	f1 := storage.Iterator("field1")
	for f1.HasNext() {
		fmt.Println(f1.Current())
		f1.Next()
	}
	fmt.Println("****")
	f2 := storage.Iterator("field2")
	for f2.HasNext() {
		fmt.Println(f2.Current())
		f2.Next()
	}

	q := NewNotAndQuery([]Query{
		NewTermQuery(storage.Iterator("field1").(*datastruct.SkipListIterator)),
		NewTermQuery(storage.Iterator("field2").(*datastruct.SkipListIterator)),
		//NewTermQuery(nil),
	}, nil)

	if cur, err := q.Current(); err == nil {
		fmt.Println(cur, err)
	}
	cur, err := q.Next()
	for err == nil {
		fmt.Println(cur, err)
		cur, err = q.Next()
	}

}
