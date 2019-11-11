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
    a := NewAndQuery(nil)
    fmt.Println(a)
}

func TestAndQuery_Next(t *testing.T) {
	sl := index.NewSkipList(index.DEFAULT_MAX_LEVEL, helpers.DocIdFunc)

	sl.Add(document.DocId(1), [1]byte{})
	sl.Add(document.DocId(2), [1]byte{})
	sl.Add(document.DocId(3), [1]byte{})
	sl.Add(document.DocId(4), [1]byte{})

	sll := &index.SkipListIterator{
		SkipList: sl,
	}
	a := NewAndQuery(&TermQuery{sll})

	Convey("Next", t, func() {
		v, e := a.Next()
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)
		v, e = a.GetGE(document.DocId(1))
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)
	})

	fmt.Println(a.Next())
	fmt.Println(a.GetGE(document.DocId(1)))

	sl.Del(document.DocId(2))

	fmt.Println(a.Next())
	fmt.Println(a.GetGE(document.DocId(2)))

	Convey("GetGE", t, func() {
		v, e := a.Next()
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)
		v, e = a.GetGE(document.DocId(2))
		So(v, ShouldEqual, 3)
		So(e, ShouldBeNil)
	})

}
