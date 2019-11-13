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

	Convey("Next", t, func() {
		a := NewAndQuery([]Query{&TermQuery{sll}}, nil)
		v, e := a.GetGE(document.DocId(1))
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)
		v, e = a.Next()
		So(e, ShouldBeNil)
		So(v, ShouldEqual, 4)
		So(e, ShouldBeNil)
		v, e = a.GetGE(document.DocId(1))
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)

	})

	Convey("GetGE", t, func() {
		a := NewAndQuery([]Query{&TermQuery{sll.Iterator()}, &TermQuery{sll1.Iterator()}}, nil)
		v, e := a.Next()
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)
		v, e = a.GetGE(document.DocId(2))
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})

}
