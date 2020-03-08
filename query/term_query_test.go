package query

import (
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewTermQuery(t *testing.T) {
	Convey("term query nil", t, func() {
		tq := NewTermQuery(nil)
		So(tq, ShouldBeNil)
	})
}

func TestTermQuery_Next(t *testing.T) {
	ss := index.NewIndex("")
	s := ss.GetInvertedIndex()
	Convey("term query next", t, func() {
		So(s.Add("fieldName\0071", document.DocId(1)), ShouldBeNil)
		So(s.Add("fieldName\0071", document.DocId(5)), ShouldBeNil)
		So(s.Add("fieldName\0071", document.DocId(6)), ShouldBeNil)
		So(s.Add("fieldName\0071", document.DocId(7)), ShouldBeNil)
		So(s.Add("fieldName\0074", document.DocId(2)), ShouldBeNil)
		tq := NewTermQuery(s.Iterator("fieldName", "1"))
		id, err := tq.Current()
		So(id, ShouldEqual, document.DocId(1))
		So(err, ShouldBeNil)
		expectCase := []document.DocId{1, 5, 6, 7}
		for _, v := range expectCase {
			id, err := tq.Current()
			So(id, ShouldEqual, v)
			So(err, ShouldBeNil)
			tq.Next()
		}
		id, err = tq.Current()
		So(id, ShouldEqual, 0)
		So(err, ShouldNotBeNil)
	})
}

func TestTermQuery_GetGE(t *testing.T) {
	ss := index.NewIndex("")
	s := ss.GetInvertedIndex()
	Convey("term query next", t, func() {
		So(s.Add("fieldName\0071", document.DocId(1)), ShouldBeNil)
		So(s.Add("fieldName\0071", document.DocId(5)), ShouldBeNil)
		So(s.Add("fieldName\0071", document.DocId(6)), ShouldBeNil)
		So(s.Add("fieldName\0071", document.DocId(7)), ShouldBeNil)
		So(s.Add("fieldName\0074", document.DocId(2)), ShouldBeNil)
		tq := NewTermQuery(s.Iterator("fieldName", "1"))
		expectCase := []document.DocId{1, 1, 5, 5, 5, 6, 7}
		getGeCase := []document.DocId{1, 1, 2, 2, 5, 6, 7}
		for i, v := range expectCase {
			id, err := tq.GetGE(getGeCase[i])
			So(id, ShouldEqual, v)
			So(err, ShouldBeNil)
		}
		id, err := tq.GetGE(8)
		So(id, ShouldEqual, 0)
		So(err, ShouldNotBeNil)
	})
}