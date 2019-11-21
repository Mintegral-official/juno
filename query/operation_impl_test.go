package query

import (
	"github.com/Mintegral-official/juno/helpers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewOperationImpl(t *testing.T) {
	var op OP = EQ
	i := make([]interface{}, 0)
	Convey("Equal", t, func() {
		impl := NewOperationImpl(10, op, helpers.IntCompare)
		So(impl.Equal(10), ShouldBeTrue)
		So(impl.Equal(11), ShouldBeFalse)
		So(impl.Equal(9), ShouldBeFalse)

		op = NE
		impl = NewOperationImpl(10, op, helpers.IntCompare)
		So(!impl.Equal(10), ShouldBeFalse)
		So(!impl.Equal(11), ShouldBeTrue)
		So(!impl.Equal(9), ShouldBeTrue)

		op = LE
		impl = NewOperationImpl(10, op, helpers.IntCompare)
		So(impl.Less(10) || impl.Equal(10), ShouldBeTrue)
		So(impl.Less(11) || impl.Equal(100), ShouldBeTrue)

		op = GE
		impl = NewOperationImpl(10, op, helpers.IntCompare)
		So(!impl.Less(10), ShouldBeTrue)
		So(!impl.Less(11), ShouldBeFalse)
		So(!impl.Less(9), ShouldBeTrue)

		op = LT
		impl = NewOperationImpl(10, op, helpers.IntCompare)
		So(impl.Less(9), ShouldBeFalse)
		So(impl.Less(10), ShouldBeFalse)
		So(impl.Less(11), ShouldBeTrue)

		op = GT
		impl = NewOperationImpl(10, op, helpers.IntCompare)
		So(!impl.Less(10) && ! impl.Equal(10), ShouldBeFalse)
		So(!impl.Less(10) && ! impl.Equal(9), ShouldBeTrue)

		op = AND
		impl = NewOperationImpl(10, op, helpers.IntCompare)
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

		op = OR
		impl = NewOperationImpl(10, op, helpers.IntCompare)
		So(impl.Less(11) || impl.Equal(10), ShouldBeTrue)
		So(impl.Less(10) || impl.Equal(10), ShouldBeTrue)
		i = make([]interface{}, 0)
		i = append(i, 1)
		i = append(i, 2)
		i = append(i, 10)
		i = append(i, 16)
		So(impl.In(i), ShouldBeTrue)

        op = NOT
		impl = NewOperationImpl(10, op, helpers.IntCompare)
		i = make([]interface{}, 0)
		i = append(i, 1)
		i = append(i, 2)
		i = append(i, 16)
		So(!impl.In(i), ShouldBeTrue)
		So(!impl.Less(9) && ! impl.Equal(1), ShouldBeTrue)

		op = IN
		impl = NewOperationImpl(10, op, helpers.IntCompare)
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
