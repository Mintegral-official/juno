package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/index"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAndQuery(t *testing.T) {
	a := NewAndQuery(nil, nil)
	fmt.Println(a)
}

func TestAndQuery_Next(t *testing.T) {
	sl := index.NewSkipList(index.DEFAULT_MAX_LEVEL, helpers.DocIdFunc)

	sl.Add(document.DocId(1), [1]byte{})
	sl.Add(document.DocId(3), [1]byte{})
	sl.Add(document.DocId(6), [1]byte{})
	sl.Add(document.DocId(10), [1]byte{})

	sl1 := index.NewSkipList(index.DEFAULT_MAX_LEVEL, helpers.DocIdFunc)

	sl1.Add(document.DocId(1), [1]byte{})
	sl1.Add(document.DocId(4), [1]byte{})
	sl1.Add(document.DocId(6), [1]byte{})
	sl1.Add(document.DocId(9), [1]byte{})

	sll := &index.SkipListIterator{
		SkipList: sl,
		Element:  nil,
	}

	sll1 := &index.SkipListIterator{
		SkipList: sl1,
		Element:  nil,
	}

	Convey("Next", t, func() {
		a := NewAndQuery(nil, &TermQuery{sll.Iterator()}, &TermQuery{sll1.Iterator()})
		v, e := a.GetGE(document.DocId(1))
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)
		v, e = a.Next()
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)
		v, e = a.GetGE(document.DocId(1))
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)

	})

	Convey("GetGE", t, func() {
		a := NewAndQuery(nil, &TermQuery{sll.Iterator()}, &TermQuery{sll1.Iterator()})
		v, e := a.Next()
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)
		v, e = a.GetGE(document.DocId(2))
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})

}
