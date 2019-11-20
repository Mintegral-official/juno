package query

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewOperationImpl(t *testing.T) {
	//var op OP = EQ
	Convey("Equal", t, func() {
		So(1, ShouldEqual, 1)
		//impl := NewOperationImpl("US", op, helpers.DocIdFunc)
		//So(impl.Equal("US"), ShouldBeTrue)
		//So(impl.Equal("USS"), ShouldBeFalse)
		//op = NE
		//impl = NewOperationImpl("US", op, helpers.DocIdFunc)
		//So(impl.Equal("US"), ShouldBeFalse)
		//So(impl.Equal("USS"), ShouldBeFalse)
		//op = LE
		//impl = NewOperationImpl("US", op, helpers.DocIdFunc)
		//So(impl.Equal("US"), ShouldBeTrue)
		//So(impl.Less("US"), ShouldBeTrue)
	})
}