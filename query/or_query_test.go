package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestOrQuery(t *testing.T) {
	sl := datastruct.NewSkipList(datastruct.DEFAULT_MAX_LEVEL, helpers.DocIdFunc)

	sl.Add(document.DocId(1), [1]byte{})
	sl.Add(document.DocId(3), [1]byte{})
	sl.Add(document.DocId(6), [1]byte{})
	sl.Add(document.DocId(10), [1]byte{})

	sl1 := datastruct.NewSkipList(datastruct.DEFAULT_MAX_LEVEL, helpers.DocIdFunc)

	sl1.Add(document.DocId(1), [1]byte{})
	sl1.Add(document.DocId(4), [1]byte{})
	sl1.Add(document.DocId(6), [1]byte{})
	sl1.Add(document.DocId(9), [1]byte{})

	sll := &datastruct.SkipListIterator{
		SkipList: sl,
		Element:  nil,
	}

	sll1 := &datastruct.SkipListIterator{
		SkipList: sl1,
		Element:  nil,
	}

	//Convey("or_query", t, func() {
	//	a := NewOrQuery([]Query{&TermQuery{sll.Iterator()}}, nil)
	//	v, e := a.Next()
	//	fmt.Println(v, e)
	//	So(v, ShouldEqual, 1)
	//	So(e, ShouldBeNil)
	//	v, e = a.Next()
	//	fmt.Println(v, e)
	//	So(v, ShouldEqual, 3)
	//	So(e, ShouldBeNil)
	//	v, e = a.Next()
	//	fmt.Println(v, e)
	//	So(v, ShouldEqual, 6)
	//	So(e, ShouldBeNil)
	//	v, e = a.Next()
	//	fmt.Println(v, e)
	//	So(v, ShouldEqual, 10)
	//	So(e, ShouldBeNil)
	//})

	Convey("Next", t, func() {
		a := NewOrQuery([]Query{&TermQuery{sll.Iterator()}}, nil)
		v, e := a.Next()
		fmt.Println(v, e)

		v, e = a.Next()
		fmt.Println(v, e)

		v, e = a.Next()
		fmt.Println(v, e)

		v, e = a.Next()
		fmt.Println(v, e)

	})

	Convey("getGE", t, func() {
		a := NewOrQuery([]Query{&TermQuery{sll.Iterator()}, &TermQuery{sll1.Iterator()}}, nil)
		v, e := a.GetGE(1)
		fmt.Println(v, e)

		v, e = a.GetGE(2)
		fmt.Println(v, e)

		v, e = a.GetGE(3)
		fmt.Println(v, e)

		v, e = a.GetGE(4)
		fmt.Println(v, e)

		v, e = a.GetGE(5)
		fmt.Println(v, e)

		v, e = a.GetGE(6)
		fmt.Println(v, e)

		v, e = a.GetGE(7)
		fmt.Println(v, e)

		v, e = a.GetGE(9)
		fmt.Println(v, e)

	})

	Convey("Next1", t, func() {
		a := NewOrQuery([]Query{&TermQuery{sll.Iterator()}, &TermQuery{sll1.Iterator()}}, nil)
		v, e := a.Next()
		fmt.Println(v, e)

		v, e = a.Next()
		fmt.Println(v, e)

		v, e = a.Next()
		fmt.Println(v, e)

		v, e = a.Next()
		fmt.Println(v, e)

		v, e = a.Next()
		fmt.Println(v, e)

		v, e = a.Next()
		fmt.Println(v, e)

		v, e = a.Next()
		fmt.Println(v, e)

		v, e = a.Next()
		fmt.Println(v, e)

	})

}
