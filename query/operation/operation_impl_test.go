package operation

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewOperationImpl(t *testing.T) {
	i := make([]interface{}, 0)
	Convey("Equal", t, func() {
		impl := NewOperationImpl(10)
		So(impl.Equal(10), ShouldBeTrue)
		So(impl.Equal(11), ShouldBeFalse)
		So(impl.Equal(9), ShouldBeFalse)

		impl = NewOperationImpl(10)
		So(!impl.Equal(10), ShouldBeFalse)
		So(!impl.Equal(11), ShouldBeTrue)
		So(!impl.Equal(9), ShouldBeTrue)

		impl = NewOperationImpl(10)
		So(impl.Less(10) || impl.Equal(10), ShouldBeTrue)
		So(impl.Less(11) || impl.Equal(100), ShouldBeTrue)

		impl = NewOperationImpl(10)
		So(!impl.Less(10), ShouldBeTrue)
		So(!impl.Less(11), ShouldBeFalse)
		So(!impl.Less(9), ShouldBeTrue)

		impl = NewOperationImpl(10, )
		So(impl.Less(9), ShouldBeFalse)
		So(impl.Less(10), ShouldBeFalse)
		So(impl.Less(11), ShouldBeTrue)

		impl = NewOperationImpl(10)
		So(!impl.Less(10) && ! impl.Equal(10), ShouldBeFalse)
		So(!impl.Less(10) && ! impl.Equal(9), ShouldBeTrue)

		impl = NewOperationImpl(10)
		So(impl.Less(11) && impl.Equal(10), ShouldBeTrue)
		So(impl.Less(10) && impl.Equal(10), ShouldBeFalse)
		So(impl.Less(11) && impl.Less(100), ShouldBeTrue)
		i = make([]interface{}, 0)
		i = append(i, 1)
		i = append(i, 2)
		i = append(i, 10)
		i = append(i, 16)
		So(impl.Less(11) && impl.In(i), ShouldBeTrue)
		So(impl.Equal(10) && impl.In(i), ShouldBeTrue)

		impl = NewOperationImpl(10)
		So(impl.Less(11) || impl.Equal(10), ShouldBeTrue)
		So(impl.Less(10) || impl.Equal(10), ShouldBeTrue)
		i = make([]interface{}, 0)
		i = append(i, 1)
		i = append(i, 2)
		i = append(i, 10)
		i = append(i, 16)
		So(impl.In(i), ShouldBeTrue)

		impl = NewOperationImpl(10)
		i = make([]interface{}, 0)
		i = append(i, 1)
		i = append(i, 2)
		i = append(i, 16)
		So(!impl.In(i), ShouldBeTrue)
		So(!impl.Less(9) && ! impl.Equal(1), ShouldBeTrue)

		impl = NewOperationImpl(10)
		i = make([]interface{}, 0)
		i = append(i, 1)
		i = append(i, 10)
		i = append(i, 16)
		So(impl.In(i), ShouldBeTrue)
		So(impl.Equal(1) || impl.Equal(10) || impl.Equal(13), ShouldBeTrue)
		So(impl.Less(1) || impl.Equal(10) || impl.Less(8), ShouldBeTrue)
		So(impl.Less(1) || impl.Equal(1132) || impl.Less(11), ShouldBeTrue)

	})
}
