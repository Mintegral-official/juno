package query

import (
	"container/heap"
	"github.com/MintegralTech/juno/datastruct"
	"github.com/MintegralTech/juno/document"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestHeap_Compare(t *testing.T) {
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

	h := &Heap{}
	heap.Push(h, Query(NewTermQuery(sl.Iterator())))
	heap.Push(h, Query(NewTermQuery(sl1.Iterator())))
	Convey("heap", t, func() {
		So(h.Len(), ShouldEqual, 2)
	})

}