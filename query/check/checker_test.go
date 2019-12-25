package check

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/query/operation"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInChecker_Check(t *testing.T) {
	sl, _ := datastruct.NewSkipList(datastruct.DefaultMaxLevel, helpers.DocIdFunc)

	sl.Add(document.DocId(1), 1)
	sl.Add(document.DocId(3), 6)
	sl.Add(document.DocId(6), 5)
	sl.Add(document.DocId(10), 10)

	sl1, _ := datastruct.NewSkipList(datastruct.DefaultMaxLevel, helpers.DocIdFunc)

	sl1.Add(document.DocId(1), 1)
	sl1.Add(document.DocId(4), 1)
	sl1.Add(document.DocId(6), 8)
	sl1.Add(document.DocId(9), 1)

	Convey("checker", t, func() {
		c := NewChecker(sl.Iterator(), 10, operation.EQ)
		So(c.Check(3), ShouldBeFalse)
		So(c.Check(10), ShouldBeTrue)
	})

	Convey("and checker", t, func() {
		c := NewChecker(sl.Iterator(), 3, operation.GE)
		d := NewChecker(sl1.Iterator(), 10, operation.LT)
		a := NewAndChecker([]Checker{
			c, d,
		})
		So(a.Check(3), ShouldBeFalse)
		So(a.Check(6), ShouldBeTrue)
		So(a.Check(10), ShouldBeFalse)
		//	So(a.Check(6), ShouldBeTrue)
	})

	Convey("or checker", t, func() {
		c := NewChecker(sl.Iterator(), 6, operation.EQ)
		d := NewChecker(sl1.Iterator(), 10, operation.EQ)
		o := NewOrChecker([]Checker{
			c, d,
		})
		So(o.Check(3), ShouldBeTrue)
		So(o.Check(6), ShouldBeFalse)
		So(o.Check(10), ShouldBeFalse)
		So(o.Check(6), ShouldBeFalse)
	})

	Convey("in checker", t, func() {
		c := NewChecker(sl.Iterator(), 6, operation.EQ)
		d := NewChecker(sl.Iterator(), 10, operation.EQ)
		o := NewInInChecker([]Checker{
			c, d,
		})
		So(o.Check(3), ShouldBeTrue)
		So(o.Check(6), ShouldBeFalse)
		So(o.Check(10), ShouldBeTrue)
		So(o.Check(6), ShouldBeFalse)
	})

	Convey("not checker", t, func() {
		c := NewChecker(sl.Iterator(), 6, operation.NE)
		d := NewChecker(sl.Iterator(), 10, operation.NE)
		o := NewNotChecker([]Checker{
			c, d,
		})
		So(o.Check(3), ShouldBeTrue)
		So(o.Check(6), ShouldBeFalse)
		So(o.Check(10), ShouldBeTrue)
		So(o.Check(6), ShouldBeTrue)
	})
}
