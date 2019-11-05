package index

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSimpleInvertedIndex_Add(t *testing.T) {
	s := NewSimpleInvertedIndex()
	Convey("Add", t, func() {
        So(s.Add("fileName", 1), ShouldBeNil)
	})
}

func TestSimpleInvertedIndex_Del(t *testing.T) {
	s := NewSimpleInvertedIndex()
	Convey("Del", t, func() {
		So(s.Del("filename", 1), ShouldBeFalse)
	})
}

func TestSimpleInvertedIndex_Iterator(t *testing.T) {
	s := NewSimpleInvertedIndex()
	Convey("Iterator", t, func() {
		So(s.Iterator("filename"), ShouldBeNil)
	})
}
