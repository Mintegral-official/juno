package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAndQuery(t *testing.T) {
	a := NewAndQuery(nil, nil)
	Convey("and query", t, func() {
		So(a, ShouldNotBeNil)
	})
}

func TestAndQuery_Next(t *testing.T) {
	sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl.Add(document.DocId(1), [1]byte{})
	sl.Add(document.DocId(3), [1]byte{})
	sl.Add(document.DocId(6), [1]byte{})
	sl.Add(document.DocId(10), [1]byte{})

	sl1 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl1.Add(document.DocId(1), [1]byte{})
	sl1.Add(document.DocId(4), [1]byte{})
	sl1.Add(document.DocId(6), [1]byte{})
	sl1.Add(document.DocId(9), [1]byte{})

	Convey("GetGE1", t, func() {
		a := NewAndQuery([]Query{NewTermQuery(sl.Iterator())}, nil)

		v, e := a.GetGE(document.DocId(1))
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(1))
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(1))
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(3))
		So(v, ShouldEqual, 3)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(4))
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(6))
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(7))
		So(v, ShouldEqual, 10)
		So(e, ShouldBeNil)
	})

	Convey("GetGE", t, func() {
		a := NewAndQuery([]Query{NewTermQuery(sl.Iterator()), NewTermQuery(sl1.Iterator())}, nil)
		v, e := a.GetGE(document.DocId(1))
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(1))
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(1))
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(1))
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(3))
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(3))
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(4))
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(6))
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)
		fmt.Println()

		v, e = a.GetGE(document.DocId(10))
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
		So(a.DebugInfo(), ShouldNotBeNil)
	})

	Convey("Next1", t, func() {
		a := NewAndQuery([]Query{NewTermQuery(sl.Iterator())}, nil)
		v, e := a.Current()
		//fmt.Println(v, e)
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.Next()
		// fmt.Println(v, e)
		So(v, ShouldEqual, 3)
		So(e, ShouldBeNil)

		v, e = a.Next()
		// fmt.Println(v, e)
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.Next()
		// fmt.Println(v, e)
		So(v, ShouldEqual, 10)
		So(e, ShouldBeNil)
	})

	Convey("Next2", t, func() {
		a := NewAndQuery([]Query{NewTermQuery(sl.Iterator()), NewTermQuery(sl1.Iterator())}, nil)
		v, e := a.Current()
		// fmt.Println(v, e)
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.Next()
		// fmt.Println(v, e)
		//	v, e = a.Current()
		//	fmt.Println(v, e)
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

	})
}

//func f(a, b interface{}) bool {
//	aa, bb := unsafe.Pointer(&a), unsafe.Pointer(&b)
//	//fmt.Println(*(*float64)(aa))
//	//fmt.Println(*(*float64)(bb))
//	return *(*float64)(aa) == *(*float64)(bb)
//}
//
//func TestAndQuery_Current(t *testing.T) {
//    fmt.Println(f(1.1, 1.1))
//    fmt.Println(f(1.1, 1.2))
//    fmt.Println(f(1, 1))
//    fmt.Println(f(2, 1))
//}
