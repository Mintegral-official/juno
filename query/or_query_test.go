package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewOrQuery_Next1(t *testing.T) {
	sl, _ := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl.Add(document.DocId(1), [1]byte{})
	sl.Add(document.DocId(3), [1]byte{})
	sl.Add(document.DocId(6), [1]byte{})
	sl.Add(document.DocId(10), [1]byte{})

	//sl1 := datastruct.NewSkipList(datastruct.DEFAULT_MAX_LEVEL, helpers.DocIdFunc)
	//
	//sl1.Add(document.DocId(1), [1]byte{})
	//sl1.Add(document.DocId(4), [1]byte{})
	//sl1.Add(document.DocId(6), [1]byte{})
	//sl1.Add(document.DocId(9), [1]byte{})

	Convey("Next1", t, func() {
		a := NewOrQuery([]Query{NewTermQuery(sl.Iterator())}, nil)
		v, e := a.Next()
		// fmt.Println(v, e)
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)
		v, e = a.Next()
		//fmt.Println(v, e)
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
}

func TestOrQuery_GetGE(t *testing.T) {
	sl, _ := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl.Add(document.DocId(1), [1]byte{})
	sl.Add(document.DocId(3), [1]byte{})
	sl.Add(document.DocId(6), [1]byte{})
	sl.Add(document.DocId(10), [1]byte{})

	sl1, _ := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl1.Add(document.DocId(1), [1]byte{})
	sl1.Add(document.DocId(4), [1]byte{})
	sl1.Add(document.DocId(6), [1]byte{})
	sl1.Add(document.DocId(9), [1]byte{})

	Convey("getGE", t, func() {
		s1 := sl.Iterator()
		s2 := sl1.Iterator()
		a := NewOrQuery([]Query{NewTermQuery(s1),NewTermQuery(s2)}, nil)
		v, e := a.GetGE(1)
		//fmt.Println(v, e)
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.GetGE(2)
		//fmt.Println(v, e)
		So(v, ShouldEqual, 3)
		So(e, ShouldBeNil)

		v, e = a.GetGE(3)
		// fmt.Println(v, e)
		So(v, ShouldEqual, 3)
		So(e, ShouldBeNil)

		v, e = a.GetGE(4)
		// fmt.Println(v, e)
		So(v, ShouldEqual, 4)
		So(e, ShouldBeNil)

		v, e = a.GetGE(5)
		// fmt.Println(v, e)
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.GetGE(6)
		// fmt.Println(v, e)
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.GetGE(7)
		// fmt.Println(v, e)
		So(v, ShouldEqual, 9)
		So(e, ShouldBeNil)

		v, e = a.GetGE(9)
		// fmt.Println(v, e)
		So(v, ShouldEqual, 9)
		So(e, ShouldBeNil)

		v, e = a.GetGE(10)
		// fmt.Println(v, e)
		So(v, ShouldEqual, 10)
		So(e, ShouldBeNil)

		v, e = a.GetGE(11)
		// fmt.Println(v, e)
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})

	Convey("GetGE2", t, func() {
		a := NewOrQuery([]Query{NewTermQuery(sl.Iterator()),NewTermQuery(sl1.Iterator())}, nil)
		v, e := a.GetGE(8)
		So(v, ShouldEqual, 9)
		So(e, ShouldBeNil)

		v, e = a.GetGE(2)
		So(v, ShouldEqual, 9)
		So(e, ShouldBeNil)

	})
}

func TestNewOrQuery_Next2(t *testing.T) {

	sl, _ := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl.Add(document.DocId(1), [1]byte{})
	sl.Add(document.DocId(3), [1]byte{})
	sl.Add(document.DocId(6), [1]byte{})
	sl.Add(document.DocId(10), [1]byte{})

	sl1, _ := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl1.Add(document.DocId(1), [1]byte{})
	sl1.Add(document.DocId(4), [1]byte{})
	sl1.Add(document.DocId(6), [1]byte{})
	sl1.Add(document.DocId(9), [1]byte{})

	sl2, _ := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl2.Add(document.DocId(2), [1]byte{})
	sl2.Add(document.DocId(5), [1]byte{})
	sl2.Add(document.DocId(7), [1]byte{})
	sl2.Add(document.DocId(8), [1]byte{})

	Convey("Next1", t, func() {
		a := NewOrQuery([]Query{NewTermQuery(sl2.Iterator()), NewTermQuery(sl.Iterator()), NewTermQuery(sl1.Iterator())}, nil)
		v, e := a.Current()
		// fmt.Println(v, e)
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.Next()
		// fmt.Println(v, e)
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.Next()
		// fmt.Println(v, e)
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.Next()
		// fmt.Println(v, e)
		So(v, ShouldEqual, 2)
		So(e, ShouldBeNil)

		v, e = a.Next()
		// fmt.Println(v, e)
		So(v, ShouldEqual, 3)
		So(e, ShouldBeNil)

		v, e = a.Next()
		// fmt.Println(v, e)
		So(v, ShouldEqual, 4)
		So(e, ShouldBeNil)

		v, e = a.Next()
		// fmt.Println(v, e)
		So(v, ShouldEqual, 5)
		So(e, ShouldBeNil)

		v, e = a.Next()
		// fmt.Println(v, e)
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.Next()
		// fmt.Println(v, e)
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.Next()
		// fmt.Println(v, e)
		So(v, ShouldEqual, 7)
		So(e, ShouldBeNil)

		v, e = a.Next()
		// fmt.Println(v, e)
		So(v, ShouldEqual, 8)
		So(e, ShouldBeNil)

		v, e = a.Next()
		// fmt.Println(v, e)
		So(v, ShouldEqual, 9)
		So(e, ShouldBeNil)

		v, e = a.Next()
		// fmt.Println(v, e)
		So(v, ShouldEqual, 10)
		So(e, ShouldBeNil)

		v, e = a.Next()
		// fmt.Println(v, e)
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)

		fmt.Println(a.String())

	})
}
