package helpers

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFunc_Compare(t *testing.T) {
	Convey("compare", t, func() {
		var a, b *rune
		var c, d rune = 2, 1
		a = &c
		b = &d
		So(intCompare(int(*a), int(*b)), ShouldEqual, 1)
		So(intCompare(*a, *b), ShouldEqual, 1)
		So(intCompare(c, d), ShouldEqual, 1)
		So(Compare(a, b), ShouldEqual, 1)
	})
}