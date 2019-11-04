package index

import (
	"github.com/Mintegral-official/juno/helpers"
	"math/rand"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	DEFAULT_MAX_LEVEL   = 12
	DEFAULT_PROBABILITY = 0x3FFF
)

type Element struct {
	key, value interface{}
	next       []unsafe.Pointer
}

func newNode(key, value interface{}, level int32) *Element {
	return &Element{key, value, make([]unsafe.Pointer, level)}
}

func (element *Element) getNext(n int) *Element {
	return (*Element)(atomic.LoadPointer(&element.next[n]))
}

func (element *Element) setNext(n int, x *Element) {
	atomic.StorePointer(&element.next[n], unsafe.Pointer(x))
}

func (element *Element) Next(n int) *Element {
	return (*Element)(element.next[n])
}

type SkipList struct {
	cmp              helpers.Comparable
	randSource       rand.Source
	header           *Element
	level            int32
	length           int64
	previousNodeCache [DEFAULT_MAX_LEVEL]*Element
}

func NewSkipList(level int32, cmp helpers.Comparable) *SkipList {
	return &SkipList{
		cmp:        cmp,
		randSource: rand.New(rand.NewSource(time.Now().UnixNano())),
		level:      level,
		length:     0,
		header:     newNode(nil, nil, DEFAULT_MAX_LEVEL),
	}
}

func (skipList *SkipList) Add(key, value interface{}) {

	if m, ok := skipList.findGE(key, true); ok {
		h := int32(len(m.next))
		x := newNode(key, value, h)
		for i, n := range skipList.previousNodeCache[:h] {
			x.setNext(i, m.getNext(i))
			n.setNext(i, x)
		}
		return
	}

	h := skipList.randLevel()

	if h > skipList.level {
		for i := skipList.level; i < h; i++ {
			skipList.previousNodeCache[i] = skipList.header
		}
		atomic.StoreInt32(&skipList.level, h)
	}
	x := newNode(key, value, h)
	for i, n := range skipList.previousNodeCache[:h] {
		x.setNext(i, n.getNext(i))
		n.setNext(i, x)
	}
	atomic.AddInt64(&skipList.length, 1)
}

func (skipList *SkipList) Del(key interface{}) {
	x, ok := skipList.findGE(key, true)
	if !ok {
		return
	}

	h := len(x.next)
	for i, n := range skipList.previousNodeCache[:h] {
		if n.Next(i) != nil {
			n.setNext(i, n.Next(i).Next(i))
		}
	}
	atomic.AddInt64(&skipList.length, -1)
}

func (skipList *SkipList) Contains(key interface{}) bool {
	_, ok := skipList.findGE(key, false)
	return ok
}

func (skipList *SkipList) Get(key interface{}) (value interface{}, err error) {
	if x, ok := skipList.findGE(key, false); ok {
		return x.value, nil
	}
	return nil, helpers.ERROR_ELEMENT_ERROR
}

func (skipList *SkipList) Len() int {
	return int(atomic.LoadInt64(&skipList.length))
}

func (skipList *SkipList) findGE(key interface{}, flag bool) (*Element, bool) {
	x := skipList.header
	h := int(atomic.LoadInt32(&skipList.level)) - 1
	for {
		next := x.getNext(h)
		cmp := 1
		if next != nil {
			cmp = skipList.cmp.Compare(next.key, key)
		}
		if cmp < 0 {
			x = next
		} else {
			if flag {
				skipList.previousNodeCache[h] = x
			} else if cmp == 0 {
				return next, true
			}
			if h == 0 {
				return next, cmp == 0
			}
			h--
		}
	}
	return nil, false
}

func (skipList *SkipList) randLevel() int32 {
	var l int32 = 1
	for ((skipList.randSource.Int63() >> 32) & 0xFFFF) < DEFAULT_PROBABILITY {
		l++
	}
	if l > DEFAULT_MAX_LEVEL {
		l = DEFAULT_MAX_LEVEL
	}
	return l
}
