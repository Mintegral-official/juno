package index

import (
	"github.com/MintegralTech/juno/datastruct"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/helpers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewStorageIndexerAdd(t *testing.T) {
	Convey("storage method test", t, func() {
		s := NewStorageIndexer()
		So(s.Get("fieldName", 1), ShouldBeNil)
		So(s.Del("fieldName", 1), ShouldBeFalse)
		So(s.Iterator("fieldName"), ShouldNotBeNil)
		So(s.Add("fieldName", 1, 1), ShouldBeNil)
		So(s.Del("fieldName", 1), ShouldBeTrue)
		So(s.Iterator("fieldName"), ShouldNotBeNil)
	})
}

func TestStorageIndexer(t *testing.T) {
	s := NewStorageIndexer()
	sl1 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
	s.data.Store("fieldName1", sl1)
	s.data.Store("fieldName2", nil)
	sl2 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
	s.data.Store("fieldName4", sl2)
	Convey("ADD & GET &DEL & ITERATOR", t, func() {
		So(s.Add("fieldName1", document.DocId(1), nil), ShouldBeNil)
		So(s.Add("fieldName1", document.DocId(2), nil), ShouldBeNil)
		So(s.Add("fieldName1", document.DocId(3), nil), ShouldBeNil)
		So(s.Add("fieldName1", document.DocId(4), nil), ShouldBeNil)
		So(s.Add("fieldName2", document.DocId(222), nil), ShouldEqual, helpers.ParseError)
		So(s.Add("fieldName4", document.DocId(444), nil), ShouldBeNil)
		So(s.Add("fieldName", document.DocId(0), nil), ShouldEqual, nil)
		a := s.Iterator("fieldName1")
		c := 0
		for a.HasNext() {
			a.Next()
			c++
		}
		So(c, ShouldEqual, 4)
		So(s.Del("fieldName1", document.DocId(1)), ShouldBeTrue)

		a = s.Iterator("fieldName1")
		c = 0
		for a.HasNext() {
			if a.Current() != nil {
				c++
			}
			a.Next()
		}
		So(c, ShouldEqual, 3)
		So(s.Del("XXX", document.DocId(1)), ShouldBeFalse)
		So(s.Get("fieldName1", document.DocId(1)), ShouldEqual, helpers.DocumentError)
		So(s.Get("fieldName1", document.DocId(2)), ShouldNotBeNil)
		So(s.Get("fieldName2", document.DocId(2)), ShouldEqual, helpers.ParseError)
		So(s.Iterator("fieldName2"), ShouldNotBeNil)
		So(s.Iterator("fieldName4"), ShouldNotBeNil)
		So(s.Iterator("fieldName0"), ShouldNotBeNil)
		So(s.Count(), ShouldEqual, 4)
	})
}
