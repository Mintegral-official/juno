package index

import (
	"fmt"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewStorageIndexImpl(t *testing.T) {
    Convey("Get", t, func() {
		s := NewStorageIndexImpl()
    	So(s.Get("fieldName", 1), ShouldBeNil)
		So(s.Del("fieldName", 1), ShouldBeFalse)
		So(s.Iterator("fieldName"), ShouldBeNil)
    	So(s.Add("fieldName", 1, 1), ShouldBeNil)
    	So(s.Del("fieldName", 1), ShouldBeTrue)
    	So(s.Iterator("fieldName"), ShouldNotBeNil)
	})
}

func TestStorageIndexImpl(t *testing.T) {
	s := NewStorageIndexImpl()
	s.data.Store("fieldName1", NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.DocIdFunc))
	s.data.Store("fieldName2", nil)
	s.data.Store("fieldName4", NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.DocIdFunc))
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
			fmt.Println(a.Next())
			c++
		}
		So(c, ShouldEqual, 4)
		So(s.Del("fieldName1", document.DocId(1)), ShouldBeTrue)

		a = s.Iterator("fieldName1")
		c = 0
		for a.HasNext() {
			fmt.Println(a.Next())
			c++
		}
		So(c, ShouldEqual, 3)
		So(s.Del("XXX", document.DocId(1)), ShouldBeFalse)
		So(s.Get("fieldName1", document.DocId(1)), ShouldEqual, helpers.DocumentError)
		So(s.Get("fieldName1", document.DocId(2)), ShouldNotBeNil)
		fmt.Println("*******")
		fmt.Println(s.Get("fieldName1", document.DocId(2)))
		So(s.Get("fieldName2", document.DocId(2)), ShouldEqual, helpers.ParseError)
	})
}
