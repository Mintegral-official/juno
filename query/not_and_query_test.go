package query

import (
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewNotAndQuery(t *testing.T) {
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

	Convey("Next1", t, func() {
		a := NewNotAndQuery([]Query{&TermQuery{sl.Iterator()}}, nil)
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


