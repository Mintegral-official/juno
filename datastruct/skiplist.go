package datastruct

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
	if element == nil {
		return nil
	}
	return (*Element)(atomic.LoadPointer(&element.next[n]))
}

func (element *Element) setNext(n int, x *Element) {
	if element == nil {
		return
	}
	atomic.StorePointer(&element.next[n], unsafe.Pointer(x))
}

func (element *Element) Next(n int) *Element {
	if element == nil {
		return nil
	}
	return (*Element)(element.next[n])
}

func (element *Element) Key() interface{} {
	return element.key
}

type SkipList struct {
	cmp               helpers.Comparable
	randSource        rand.Source
	header            *Element
	level             int32
	length            int64
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
	prev := skipList.previousNodeCache
	if m, ok := skipList.findGE(key, true, prev); ok && skipList.cmp.Compare(m.key, key) == 0 {
		h := int32(len(m.next))
		x := newNode(key, value, h)
		for i, n := range skipList.previousNodeCache[:h] {
			x.setNext(i, m.getNext(i))
			n.setNext(i, x)
		}
		return
	}

	h := skipList.randLevel()

	x := newNode(key, value, h)
	for i, n := range skipList.previousNodeCache[:h] {
		x.setNext(i, n.getNext(i))
		n.setNext(i, x)
	}
	skipList.length++
}

func (skipList *SkipList) Del(key interface{}) {
	prev := skipList.previousNodeCache
	x, ok := skipList.findGE(key, true, prev)
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
	prev := skipList.previousNodeCache
	_, ok := skipList.findGE(key, true, prev)
	return ok
}

func (skipList *SkipList) Get(key interface{}) (*Element, error) {
	prev := skipList.previousNodeCache
	if x, ok := skipList.findGE(key, true, prev); ok {
		return x, nil
	}
	return nil, helpers.ElementNotfound
}

func (skipList *SkipList) Len() int64 {
	return atomic.LoadInt64(&skipList.length)
}

func (skipList *SkipList) findGE(key interface{}, flag bool, element [DEFAULT_MAX_LEVEL]*Element) (*Element, bool) {
	x := skipList.header
	h := int(atomic.LoadInt32(&skipList.level)) - 1
	for h >= 0 {
		if x == nil {
			return nil, false
		}
		next := x.getNext(h)
		cmp := 1
		if next != nil {
			cmp = skipList.cmp.Compare(next.key, key)
		}
		if cmp < 0 {
			x = next
		} else {
			if flag {
				element[h] = x
				skipList.previousNodeCache[h] = element[h]
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

func (skipList *SkipList) findLT(key interface{}) (*Element, bool) {
	x := skipList.header
	h := int(atomic.LoadInt32(&skipList.level)) - 1
	for h >= 0 {
		next := x.getNext(h)
		if next == nil || skipList.cmp.Compare(next.key, key) >= 0 {
			if h == 0 {
				if x == skipList.header {
					return nil, false
				}
				return x, true
			}
			h--
		} else {
			x = next
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
