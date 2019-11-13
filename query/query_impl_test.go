package query

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewQueryImpl(t *testing.T) {
	termQuery, err := NewTermQuery()
	if err != nil {
		fmt.Println(nil)
	}
	a := NewQuery(NewAndQuery( nil, termQuery))
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
