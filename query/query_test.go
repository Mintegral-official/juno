package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/index"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestQuery(t *testing.T) {
	convey.Convey("query", t, func() {
		idx := index.NewIndex("index")
		_ = idx.Add(doc1)
		_ = idx.Add(doc2)
		_ = idx.Add(doc3)
		storage := idx.GetStorageIndex()
		q := NewNotAndQuery([]Query{
			NewTermQuery(storage.Iterator("field1")),
			NewTermQuery(storage.Iterator("field2")),
			NewAndQuery([]Query{
				NewTermQuery(storage.Iterator("field1")),
				NewTermQuery(storage.Iterator("field2")),
			}, nil),
		}, nil)

		cur, err := q.Current()
		convey.So(cur, convey.ShouldEqual, 0)
		convey.So(err, convey.ShouldNotBeNil)
		cur, err = q.Next()
		convey.So(cur, convey.ShouldEqual, 0)
		convey.So(err, convey.ShouldNotBeNil)
		cur, err = q.Next()
		convey.So(cur, convey.ShouldEqual, 0)
		convey.So(err, convey.ShouldNotBeNil)
		fmt.Println(q.String())
	})
}

func TestNotAndQuery_Next(t *testing.T) {
	convey.Convey("query", t, func() {
		idx := index.NewIndex("index")
		_ = idx.Add(doc1)
		_ = idx.Add(doc2)
		_ = idx.Add(doc3)
		storage := idx.GetStorageIndex()
		q := NewNotAndQuery([]Query{
			NewTermQuery(storage.Iterator("field1")),
			NewTermQuery(storage.Iterator("field2")),
			NewTermQuery(storage.Iterator("field3")),
		}, nil)

		cur, err := q.Current()
		convey.So(cur, convey.ShouldEqual, 0)
		convey.So(err, convey.ShouldNotBeNil)
		cur, err = q.Next()
		convey.So(cur, convey.ShouldEqual, 0)
		convey.So(err, convey.ShouldNotBeNil)
		cur, err = q.Next()
		convey.So(cur, convey.ShouldEqual, 0)
		convey.So(err, convey.ShouldNotBeNil)
	})
}
