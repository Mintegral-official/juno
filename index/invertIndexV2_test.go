package index

import (
	"github.com/MintegralTech/juno/datastruct"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/helpers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInvertedIndexerV2_Add(t *testing.T) {
	s := NewInvertedIndexV2()
	Convey("Add", t, func() {
		So(s.Add("fileName", 1), ShouldBeNil)
	})
}

func TestInvertedIndexerV2_Iterator(t *testing.T) {
	s := NewInvertedIndexV2()
	Convey("Iterator", t, func() {
		So(s.Iterator("filename", "nil"), ShouldNotBeNil)
	})
}

func TestInvertedIndexerV2(t *testing.T) {
	s := NewInvertedIndexV2()
	sl1 := datastruct.NewSlice()
	s.data.Store("fieldName1", sl1)
	s.data.Store("fieldName2", nil)
	sl2 := datastruct.NewSlice()
	s.data.Store("fieldName4", sl2)
	Convey("Add", t, func() {
		So(s.Add("fieldName\0071", document.DocId(1)), ShouldBeNil)
		So(s.Add("fieldName\0071", document.DocId(5)), ShouldBeNil)
		So(s.Add("fieldName\0071", document.DocId(6)), ShouldBeNil)
		So(s.Add("fieldName\0071", document.DocId(7)), ShouldBeNil)
		So(s.Add("fieldName\0074", document.DocId(2)), ShouldBeNil)
		//So(s.Del("fieldName", 1), ShouldBeFalse)
		//So(s.Del("fieldName\0071", document.DocId(1)), ShouldBeTrue)
		a := s.Iterator("fieldName", "1")
		So(s.Iterator("fieldName", "1"), ShouldNotBeNil)
		So(s.Iterator("fieldName", "2"), ShouldNotBeNil)
		So(s.Iterator("fieldName", "4"), ShouldNotBeNil)
		So(s.Iterator("fieldName", "0"), ShouldNotBeNil)
		c := 0
		for a.HasNext() {
			a.Next()
			c++
		}
		So(c, ShouldEqual, 4)
		So(s.Count(), ShouldEqual, 5)
		So(s.Add("fieldName2", 11), ShouldEqual, helpers.ParseError)
		So(s.Add("fieldName", 111), ShouldBeNil)
		if v, ok := s.data.Load("fieldName"); ok {
			if v1, ok := v.(*datastruct.Slice); ok {
				So(v1.Len(), ShouldEqual, 1)
			}
		}
	})
}
