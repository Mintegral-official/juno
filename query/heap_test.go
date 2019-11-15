package query

import (
	"container/heap"
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"testing"
)

func TestHeap_Compare(t *testing.T) {
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

	h := &Heap{}
	heap.Push(h, Query(&TermQuery{sll.Iterator()}))
	heap.Push(h, Query(&TermQuery{sll1.Iterator()}))
	fmt.Println(h.Len())
	fmt.Println(h.Pop().(Query))
}