package datastruct

import (
	"github.com/MintegralTech/juno/document"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSliceIterator(t *testing.T) {
	sl := NewSlice()
	sl.Add(1, nil)
	sl.Add(3, nil)
	Convey("New Slice Iterator", t, func() {
		iter := sl.Iterator()

		v := iter.Current()
		So(v, ShouldNotBeNil)
		So(v.Key(), ShouldEqual, 1)

		So(iter.HasNext(), ShouldBeTrue)
		iter.Next()
		v = iter.Current()
		So(v, ShouldNotBeNil)
		So(v.Key(), ShouldEqual, 3)

		So(iter.HasNext(), ShouldBeTrue)
		iter.Next()
		So(iter.Current(), ShouldBeNil)
	})

	s := NewSlice()
	for i := 0; i < 100; i++ {
		s.Add(document.DocId(i), nil)
	}
	for i := 101; i < 150; i += 3 {
		s.Add(document.DocId(i), nil)
	}

	Convey("Next", t, func() {
		iter := s.Iterator()
		So(iter.HasNext(), ShouldBeTrue)
		v := iter.Current()
		So(v, ShouldNotBeNil)
		So(v.Key(), ShouldEqual, 0)

		iter.Next()
		elem := iter.Current()
		So(elem, ShouldNotBeNil)
		So(elem, ShouldNotBeNil)
		So(elem.Key(), ShouldEqual, 1)

		v = iter.GetGE(5)
		So(v, ShouldNotBeNil)
		So(v.Key(), ShouldEqual, 5)
		elem = iter.Current()
		So(elem, ShouldNotBeNil)
		So(elem, ShouldNotBeNil)
		So(elem.Key(), ShouldEqual, 5)

		iter.Next()
		elem = iter.Current()
		So(elem, ShouldNotBeNil)
		So(elem, ShouldNotBeNil)
		So(elem.Key(), ShouldEqual, 6)

		v = iter.GetGE(102)
		So(v, ShouldNotBeNil)
		So(v.Key(), ShouldEqual, 104)
		elem = iter.Current()
		So(elem, ShouldNotBeNil)
		So(elem, ShouldNotBeNil)
		So(elem.Key(), ShouldEqual, 104)
		So(iter.HasNext(), ShouldBeTrue)

		v = iter.GetGE(147)
		So(v, ShouldNotBeNil)
		So(v.Key(), ShouldEqual, 149)
		elem = iter.Current()
		So(elem, ShouldNotBeNil)
		So(elem, ShouldNotBeNil)
		So(elem.Key(), ShouldEqual, 149)
		So(iter.HasNext(), ShouldBeTrue)

		v = iter.GetGE(160)
		So(v, ShouldBeNil)
		elem = iter.Current()
		So(elem, ShouldBeNil)
		So(iter.HasNext(), ShouldBeFalse)
	})
}

func TestSliceIterator_GetGE(t *testing.T) {
	s := NewSlice()
	for i := 0; i < 100; i++ {
		s.Add(document.DocId(i), nil)
	}
	a := s.Iterator()

	Convey("getGE", t, func() {
		v := a.GetGE(99)
		So(v.key, ShouldEqual, 99)

		v = a.GetGE(99)
		So(v.key, ShouldEqual, 99)

		v = a.GetGE(99)
		So(v.key, ShouldEqual, 99)
	})

	aa := NewSlice()
	for i := 0; i < 1000; i++ {
		aa.Add(document.DocId(i), nil)
	}
	s1 := aa.Iterator()
	Convey("del", t, func() {
		v := s1.GetGE(10)
		So(v.key, ShouldEqual, 10)
		v = s1.GetGE(324)
		So(v.key, ShouldEqual, 324)
		So(aa.Len(), ShouldEqual, 1000)
		v = s1.GetGE(10)
		So(v.key, ShouldEqual, 324)
		v = s1.GetGE(324)
		So(v.key, ShouldEqual, 324)
	})

}
