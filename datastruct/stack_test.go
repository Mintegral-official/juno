package datastruct

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewStack(t *testing.T) {
	s := NewStack()
	Convey("stack", t, func() {
		So(s, ShouldNotBeNil)
		So(s.Empty(), ShouldBeTrue)
		So(s.Len(), ShouldEqual, 0)
		So(s.Peek(), ShouldBeNil)
		So(s.Pop(), ShouldBeNil)
		s.Push("world")
		s.Push("hello")
		So(s.Empty(), ShouldBeFalse)
		So(s.Len(), ShouldEqual, 2)
		So(s.Peek(), ShouldEqual, "hello")
		So(s.Pop(), ShouldEqual, "hello")
		So(s.Empty(), ShouldBeFalse)
		So(s.Len(), ShouldEqual, 1)
		So(s.Peek(), ShouldEqual, "world")
		So(s.Pop(), ShouldEqual, "world")
		So(s.Empty(), ShouldBeTrue)
		So(s.Len(), ShouldEqual, 0)
		So(s.Peek(), ShouldBeNil)
		So(s.Pop(), ShouldBeNil)
	})
}
