package query

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewNotAndQuery_Next1(t *testing.T) {
	sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl.Add(document.DocId(1), [1]byte{})
	sl.Add(document.DocId(3), [1]byte{})
	sl.Add(document.DocId(6), [1]byte{})
	sl.Add(document.DocId(10), [1]byte{})

	Convey("not and query next1", t, func() {
		a := NewNotAndQuery([]Query{NewTermQuery(sl.Iterator())}, nil)
		v, e := a.Current()
		a.Next()
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.Current()
		a.Next()
		So(v, ShouldEqual, 3)
		So(e, ShouldBeNil)

		v, e = a.Current()
		a.Next()
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.Current()
		a.Next()
		So(v, ShouldEqual, 10)
		So(e, ShouldBeNil)
	})
}

func TestNewNotAndQuery_Next2(t *testing.T) {
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

	sl2 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
	sl2.Add(document.DocId(3), [1]byte{})

	Convey("not and query next2", t, func() {
		a := NewNotAndQuery([]Query{
			NewTermQuery(sl.Iterator()), NewTermQuery(sl1.Iterator()), NewTermQuery(sl2.Iterator()),
		}, nil)

		v, e := a.Current()
		a.Next()
		So(v, ShouldEqual, 10)
		So(e, ShouldBeNil)

		v, e = a.Current()
		a.Next()
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})
}

func TestNewNotAndQuery_GetGE(t *testing.T) {
	sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl.Add(document.DocId(1), [1]byte{})
	sl.Add(document.DocId(3), [1]byte{})
	sl.Add(document.DocId(6), [1]byte{})
	sl.Add(document.DocId(10), [1]byte{})

	Convey("not and query get1", t, func() {
		s1 := sl.Iterator()
		a := NewNotAndQuery([]Query{NewTermQuery(s1)}, nil)
		v, e := a.GetGE(1)
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.GetGE(2)
		So(v, ShouldEqual, 3)
		So(e, ShouldBeNil)

		v, e = a.GetGE(3)
		So(v, ShouldEqual, 3)
		So(e, ShouldBeNil)

		v, e = a.GetGE(4)
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.GetGE(5)
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.GetGE(6)
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.GetGE(7)
		So(v, ShouldEqual, 10)
		So(e, ShouldBeNil)

		v, e = a.GetGE(9)
		So(v, ShouldEqual, 10)
		So(e, ShouldBeNil)

		v, e = a.GetGE(10)
		So(v, ShouldEqual, 10)
		So(e, ShouldBeNil)

		v, e = a.GetGE(11)
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})
}

func TestNewNotAndQuery_GetGE2(t *testing.T) {
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

	Convey("not and query next2", t, func() {
		s1 := sl.Iterator()
		s2 := sl1.Iterator()
		a := NewNotAndQuery([]Query{NewTermQuery(s1), NewTermQuery(s2)}, nil)
		v, e := a.GetGE(1)
		So(v, ShouldEqual, 3)
		So(e, ShouldBeNil)

		v, e = a.GetGE(2)
		So(v, ShouldEqual, 3)
		So(e, ShouldBeNil)

		v, e = a.GetGE(3)
		So(v, ShouldEqual, 3)
		So(e, ShouldBeNil)

		v, e = a.GetGE(4)
		So(v, ShouldEqual, 10)
		So(e, ShouldBeNil)

		v, e = a.GetGE(5)
		So(v, ShouldEqual, 10)
		So(e, ShouldBeNil)

		v, e = a.GetGE(6)
		So(v, ShouldEqual, 10)
		So(e, ShouldBeNil)

		v, e = a.GetGE(7)
		So(v, ShouldEqual, 10)
		So(e, ShouldBeNil)

		v, e = a.GetGE(9)
		So(v, ShouldEqual, 10)
		So(e, ShouldBeNil)

		v, e = a.GetGE(10)
		So(v, ShouldEqual, 10)
		So(e, ShouldBeNil)

		v, e = a.GetGE(11)
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})
}

func TestNewAndQuery2(t *testing.T) {
	sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl.Add(document.DocId(1), [1]byte{})
	sl.Add(document.DocId(3), [1]byte{})
	sl.Add(document.DocId(6), [1]byte{})
	sl.Add(document.DocId(10), [1]byte{})
	sl.Add(document.DocId(11), [1]byte{})

	sl1 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl1.Add(document.DocId(1), [1]byte{})
	sl1.Add(document.DocId(4), [1]byte{})
	sl1.Add(document.DocId(6), [1]byte{})
	sl1.Add(document.DocId(9), [1]byte{})

	sl2 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl2.Add(document.DocId(1), [1]byte{})
	sl2.Add(document.DocId(4), [1]byte{})
	sl2.Add(document.DocId(6), [1]byte{})
	sl2.Add(document.DocId(9), [1]byte{})

	//TODO
	//q := NewNotAndQuery([]Query{
	//	NewTermQuery(sl.Iterator()),
	//	NewTermQuery(sl1.Iterator()),
	//	NewOrQuery([]Query{
	//		NewTermQuery(sl2.Iterator()),
	//		NewTermQuery(sl1.Iterator()),
	//	}, nil),
	//}, nil)
	//
	//fmt.Println(q.Next())
	//fmt.Println(q.Next())
	//fmt.Println(q.Next())
	//fmt.Println(q.Next())
}
