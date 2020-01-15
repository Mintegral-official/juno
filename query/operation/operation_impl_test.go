package operation

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewOperations(t *testing.T) {
	i := make([]int, 0)
	Convey("operation interface test", t, func() {
		impl := NewOperations(10)
		So(impl.Equal(10), ShouldBeTrue)
		So(impl.Equal(11), ShouldBeFalse)
		So(impl.Equal(9), ShouldBeFalse)

		impl = NewOperations(10)
		So(!impl.Equal(10), ShouldBeFalse)
		So(!impl.Equal(11), ShouldBeTrue)
		So(!impl.Equal(9), ShouldBeTrue)

		impl = NewOperations(10)
		So(impl.Less(10) || impl.Equal(10), ShouldBeTrue)
		So(impl.Less(11) || impl.Equal(100), ShouldBeTrue)

		impl = NewOperations(10)
		So(!impl.Less(10), ShouldBeTrue)
		So(!impl.Less(11), ShouldBeFalse)
		So(!impl.Less(9), ShouldBeTrue)

		impl = NewOperations(10, )
		So(impl.Less(9), ShouldBeFalse)
		So(impl.Less(10), ShouldBeFalse)
		So(impl.Less(11), ShouldBeTrue)

		impl = NewOperations(10)
		So(!impl.Less(10) && ! impl.Equal(10), ShouldBeFalse)
		So(!impl.Less(10) && ! impl.Equal(9), ShouldBeTrue)

		impl = NewOperations(10)
		So(impl.Less(11) && impl.Equal(10), ShouldBeTrue)
		So(impl.Less(10) && impl.Equal(10), ShouldBeFalse)
		So(impl.Less(11) && impl.Less(100), ShouldBeTrue)
		i = make([]int, 0)
		i = append(i, 1)
		i = append(i, 2)
		i = append(i, 10)
		i = append(i, 16)
		So(impl.Less(11) && impl.In(i), ShouldBeTrue)
		So(impl.Equal(10) && impl.In(i), ShouldBeTrue)

		impl = NewOperations(10)
		So(impl.Less(11) || impl.Equal(10), ShouldBeTrue)
		So(impl.Less(10) || impl.Equal(10), ShouldBeTrue)
		i = make([]int, 0)
		i = append(i, 1)
		i = append(i, 2)
		i = append(i, 10)
		i = append(i, 16)
		So(impl.In(i), ShouldBeTrue)

		impl = NewOperations(10)
		i = make([]int, 0)
		i = append(i, 1)
		i = append(i, 2)
		i = append(i, 16)
		So(!impl.In(i), ShouldBeTrue)
		So(!impl.Less(9) && ! impl.Equal(1), ShouldBeTrue)

		impl = NewOperations(10)
		i = make([]int, 0)
		i = append(i, 1)
		i = append(i, 10)
		i = append(i, 16)
		So(impl.In(i), ShouldBeTrue)
		So(impl.Equal(1) || impl.Equal(10) || impl.Equal(13), ShouldBeTrue)
		So(impl.Less(1) || impl.Equal(10) || impl.Less(8), ShouldBeTrue)
		So(impl.Less(1) || impl.Equal(1132) || impl.Less(11), ShouldBeTrue)

	})
}
