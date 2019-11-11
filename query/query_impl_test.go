package query

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewQueryImpl(t *testing.T) {
	a := NewQuery(NewAndQuery(NewTermQuery()))
	Convey("Next", t, func() {
		v, e := a.Next()
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})
	Convey("GetGE", t, func() {
		v, e := a.GetGE(0)
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})

}
