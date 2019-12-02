package check

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/query/operation"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCheckerImpl_Check(t *testing.T) {
	sl, _ := datastruct.NewSkipList(datastruct.DefaultMaxLevel, helpers.DocIdFunc)

	sl.Add(document.DocId(1), [1]byte{})
	sl.Add(document.DocId(3), [1]byte{})
	sl.Add(document.DocId(6), [1]byte{})
	sl.Add(document.DocId(10), [1]byte{})

	sl1, _ := datastruct.NewSkipList(datastruct.DefaultMaxLevel, helpers.DocIdFunc)

	sl1.Add(document.DocId(1), [1]byte{})
	sl1.Add(document.DocId(4), [1]byte{})
	sl1.Add(document.DocId(6), [1]byte{})
	sl1.Add(document.DocId(9), [1]byte{})

	Convey("checker", t, func() {
		c := NewCheckerImpl(sl.Iterator(), document.DocId(10), operation.EQ)
		So(c.Check(3), ShouldBeFalse)
		So(c.Check(10), ShouldBeTrue)
	})

	Convey("Andchecker", t, func() {
		c := NewCheckerImpl(sl.Iterator(), document.DocId(6), operation.EQ)
		d := NewCheckerImpl(sl1.Iterator(), document.DocId(10), operation.LT)
		a := NewAndCheckerImpl([]Checker{
			c, d,
		})
		So(a.Check(3), ShouldBeFalse)
		So(a.Check(6), ShouldBeTrue)
		So(a.Check(10), ShouldBeFalse)
		So(a.Check(6), ShouldBeTrue)
	})

	Convey("Orchecker", t, func() {
		c := NewCheckerImpl(sl.Iterator(), document.DocId(6), operation.EQ)
		d := NewCheckerImpl(sl1.Iterator(), document.DocId(10), operation.EQ)
		o := NewOrCheckerImpl([]Checker{
			c, d,
		})
		So(o.Check(3), ShouldBeFalse)
		So(o.Check(6), ShouldBeTrue)
		So(o.Check(10), ShouldBeTrue)
		So(o.Check(6), ShouldBeTrue)
	})

	Convey("Inchecker", t, func() {
		c := NewCheckerImpl(sl.Iterator(), document.DocId(6), operation.EQ)
		d := NewCheckerImpl(sl.Iterator(), document.DocId(10), operation.EQ)
		o := NewInCheckerImpl([]Checker{
			c, d,
		})
		So(o.Check(3), ShouldBeFalse)
		So(o.Check(6), ShouldBeTrue)
		So(o.Check(10), ShouldBeTrue)
		So(o.Check(6), ShouldBeTrue)
	})

	Convey("Notchecker", t, func() {
		c := NewCheckerImpl(sl.Iterator(), document.DocId(6), operation.NE)
		d := NewCheckerImpl(sl.Iterator(), document.DocId(10), operation.NE)
		o := NewNotCheckerImpl([]Checker{
			c, d,
		})
		So(o.Check(3), ShouldBeTrue)
		So(o.Check(6), ShouldBeFalse)
		So(o.Check(10), ShouldBeFalse)
		So(o.Check(6), ShouldBeFalse)
	})
}
