package check

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/query/operation"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInChecker_Check(t *testing.T) {
	sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl.Add(1, 1)
	sl.Add(3, 6)
	sl.Add(6, 5)
	sl.Add(10, 10)

	sl1 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl1.Add(1, 1)
	sl1.Add(4, 1)
	sl1.Add(6, 8)
	sl1.Add(9, 1)

	Convey("checker", t, func() {
		c := NewChecker(sl.Iterator(), 10, operation.EQ, nil)
		So(c.Check(3), ShouldBeFalse)
		So(c.Check(10), ShouldBeTrue)
	})

	Convey("and checker", t, func() {
		c := NewChecker(sl.Iterator(), 3, operation.GE, nil)
		d := NewChecker(sl1.Iterator(), 10, operation.LT, nil)
		a := NewAndChecker([]Checker{
			c, d,
		})
		So(a.Check(3), ShouldBeFalse)
		So(a.Check(6), ShouldBeTrue)
		So(a.Check(10), ShouldBeFalse)
		//	So(a.Check(6), ShouldBeTrue)
	})

	Convey("or checker", t, func() {
		c := NewChecker(sl.Iterator(), 6, operation.EQ, nil)
		d := NewChecker(sl1.Iterator(), 10, operation.EQ, nil)
		o := NewOrChecker([]Checker{
			c, d,
		})
		So(o.Check(3), ShouldBeTrue)
		So(o.Check(6), ShouldBeFalse)
		So(o.Check(10), ShouldBeFalse)
		So(o.Check(6), ShouldBeFalse)
	})

	Convey("in checker", t, func() {
		//c := NewChecker(sl.Iterator(), 6, operation.EQ)
		//d := NewChecker(sl.Iterator(), 10, operation.EQ)
		var a = []int{1, 6, 3, 10}
		c := make([]interface{}, len(a))
		for _, v := range a {
			c = append(c, v)
		}
		o := NewInChecker(sl.Iterator(), c, nil)
		So(o.Check(3), ShouldBeTrue)
		So(o.Check(6), ShouldBeFalse)
		So(o.Check(10), ShouldBeTrue)
		So(o.Check(6), ShouldBeFalse)
	})

	Convey("not checker", t, func() {
		//c := NewChecker(sl.Iterator(), 6, operation.NE)
		//d := NewChecker(sl.Iterator(), 10, operation.NE)
		var a = []int{1, 6, 3, 10}
		c := make([]interface{}, len(a))
		for _, v := range a {
			c = append(c, v)
		}
		o := NewNotChecker(sl.Iterator(), c, nil)
		So(o.Check(3), ShouldBeFalse)
		So(o.Check(6), ShouldBeTrue)
		So(o.Check(10), ShouldBeFalse)
		So(o.Check(6), ShouldBeFalse)
	})
}

func TestNewChecker(t *testing.T) {
	sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
	sl.Add(1, []int{1, 2, 3})
	sl.Add(3, []int{4, 5, 3})
	sl.Add(6, []int{6, 8, 1})
	sl.Add(10, []int{10})
	Convey("in checker 2", t, func() {
		var a = []int{1, 6, 3, 10}
		c := make([]interface{}, 1)
		c = append(c, a)
		o := NewInChecker(sl.Iterator(), c, &myOperation{value: c})
		So(o.Check(3), ShouldBeTrue)
		So(o.Check(6), ShouldBeTrue)
		So(o.Check(10), ShouldBeTrue)
		So(o.Check(6), ShouldBeFalse)
	})
}

type myOperation struct {
	value interface{}
}

func (o *myOperation) Equal(value interface{}) bool {
	return true
}
func (o *myOperation) Less(value interface{}) bool {
	return true
}
func (o *myOperation) In(value []interface{}) bool {
	return true
}

func (o *myOperation) SetValue(value interface{}) {
	o.value = value
}
