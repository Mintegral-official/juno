package check

import (
	"encoding/json"
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
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

	Convey("EQ checker", t, func() {
		c := NewChecker(sl.Iterator(), 10, operation.EQ, nil, false)
		So(c.Check(3), ShouldBeFalse)
		So(c.Check(10), ShouldBeTrue)
	})

	Convey("And checker GE & LT", t, func() {
		c := NewChecker(sl.Iterator(), 3, operation.GE, nil, false)
		d := NewChecker(sl1.Iterator(), 10, operation.LT, nil, false)
		a := NewAndChecker([]Checker{
			c, d,
		})
		So(a.Check(3), ShouldBeFalse)
		So(a.Check(6), ShouldBeTrue)
		So(a.Check(10), ShouldBeFalse)
	})

	Convey("Or checker EQ EQ", t, func() {
		c := NewChecker(sl.Iterator(), 6, operation.EQ, nil, false)
		d := NewChecker(sl1.Iterator(), 10, operation.EQ, nil, false)
		o := NewOrChecker([]Checker{
			c, d,
		})
		So(o.Check(3), ShouldBeTrue)
		So(o.Check(6), ShouldBeFalse)
		So(o.Check(10), ShouldBeFalse)
		So(o.Check(6), ShouldBeFalse)
	})

	Convey("In checker", t, func() {
		var a = []int{1, 6, 3, 10}
		c := make([]int, len(a))
		for _, v := range a {
			c = append(c, v)
		}
		o := NewInChecker(sl.Iterator(), c, nil, false)
		So(o.Check(3), ShouldBeTrue)
		So(o.Check(6), ShouldBeFalse)
		So(o.Check(10), ShouldBeTrue)
		So(o.Check(6), ShouldBeFalse)
	})

	Convey("Not checker", t, func() {
		var a = []int{1, 6, 3, 10}
		c := make([]int, len(a))
		for _, v := range a {
			c = append(c, v)
		}
		o := NewNotChecker(sl.Iterator(), c, nil, false)
		So(o.Check(3), ShouldBeFalse)
		So(o.Check(6), ShouldBeTrue)
		So(o.Check(10), ShouldBeFalse)
		So(o.Check(6), ShouldBeFalse)
	})
}

func TestNewInChecker(t *testing.T) {
	sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
	sl.Add(1, []int{1, 2, 3})
	sl.Add(3, []int{4, 5, 3})
	sl.Add(6, []int{6, 8, 1})
	sl.Add(10, []int{10})
	Convey("in checker 2", t, func() {
		var a = []int{1, 6, 3, 10}
		c := make([]interface{}, 1)
		c = append(c, a)
		o := NewInChecker(sl.Iterator(), c, &myOperation{value: c}, false)
		So(o.Check(3), ShouldBeTrue)
		So(o.Check(6), ShouldBeTrue)
		So(o.Check(10), ShouldBeTrue)
		So(o.Check(6), ShouldBeFalse)
	})
}

func TestUtilCheck(t *testing.T) {
	ss := index.NewIndex("")
	s := ss.GetStorageIndex()
	Convey("Add", t, func() {
		So(s.Add("fieldName_1", 1, 1), ShouldBeNil)
		So(s.Add("fieldName_1", 5, 1), ShouldBeNil)
		So(s.Add("fieldName_1", 6, 2), ShouldBeNil)
		So(s.Add("fieldName_1", 7, 2), ShouldBeNil)
		So(s.Add("fieldName_1", 8, 3), ShouldBeNil)
		So(s.Add("fieldName_1", 9, 3), ShouldBeNil)
		c := NewChecker(s.Iterator("fieldName_1"), 3, operation.EQ, nil, false)
		So(c.Check(1), ShouldBeFalse)
		So(c.Check(8), ShouldBeTrue)
		So(c.Check(9), ShouldBeTrue)
		tmp := c.Marshal()
		res, _ := json.Marshal(tmp)
		fmt.Println(string(res))
		cc := c.Unmarshal(ss, tmp, nil)
		So(cc.Check(1), ShouldBeFalse)
		So(cc.Check(8), ShouldBeTrue)
		So(cc.Check(9), ShouldBeTrue)
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
func (o *myOperation) In(value interface{}) bool {
	return true
}

func (o *myOperation) SetValue(value interface{}) {
	o.value = value
}
