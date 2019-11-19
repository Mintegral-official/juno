package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAndQuery(t *testing.T) {
	a := NewAndQuery(nil, nil)
	fmt.Println(a)
}

func TestAndQuery_Next(t *testing.T) {
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

	Convey("GetGE1", t, func() {
		a := NewAndQuery([]Query{&TermQuery{sl.Iterator()}}, nil)

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
		a := NewAndQuery([]Query{&TermQuery{sl.Iterator()}, &TermQuery{sl1.Iterator()}}, nil)
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
        fmt.Println("******************")

		v, e = a.GetGE(document.DocId(3))
		fmt.Println(v, e)
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(4))
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(6))
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(10))
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)

	})

	Convey("Next1", t, func() {
		a := NewAndQuery([]Query{&TermQuery{sl.Iterator()}}, nil)
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
		a := NewAndQuery([]Query{&TermQuery{sl.Iterator()}, &TermQuery{sl1.Iterator()}}, nil)
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
