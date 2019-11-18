package datastruct

import (
	"github.com/Mintegral-official/juno/helpers"
	"math/rand"
	"sync/atomic"
	"time"
)

const (
	DEFAULT_MAX_LEVEL   = 12
	DEFAULT_PROBABILITY = 0x3FFF
)

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

func (sl *SkipList) Add(key, value interface{}) {
	prev := sl.previousNodeCache
	if m, ok := sl.findGE(key, true, prev); ok && sl.cmp.Compare(m.key, key) == 0 {
		h := int32(len(m.next))
		x := newNode(key, value, h)
		for i, n := range sl.previousNodeCache[:h] {
			x.setNext(i, m.getNext(i))
			n.setNext(i, x)
		}
		return
	}

	h := sl.randLevel()

	x := newNode(key, value, h)
	for i, n := range sl.previousNodeCache[:h] {
		x.setNext(i, n.getNext(i))
		n.setNext(i, x)
	}
	sl.length++
}

func (sl *SkipList) Del(key interface{}) {
	prev := sl.previousNodeCache
	x, ok := sl.findGE(key, true, prev)
	if !ok {
		return
	}

	h := len(x.next)
	for i, n := range sl.previousNodeCache[:h] {
		if n.Next(i) != nil {
			n.setNext(i, n.Next(i).Next(i))
		}
	}
	atomic.AddInt64(&sl.length, -1)
}

func (sl *SkipList) Contains(key interface{}) bool {
	prev := sl.previousNodeCache
	_, ok := sl.findGE(key, true, prev)
	return ok
}

func (sl *SkipList) Get(key interface{}) (*Element, error) {
	prev := sl.previousNodeCache
	if x, ok := sl.findGE(key, true, prev); ok {
		return x, nil
	}
	return nil, helpers.ElementNotfound
}

func (sl *SkipList) Len() int64 {
	return atomic.LoadInt64(&sl.length)
}

func (sl *SkipList) findGE(key interface{}, flag bool, element [DEFAULT_MAX_LEVEL]*Element) (*Element, bool) {
	x := sl.header
	h := int(atomic.LoadInt32(&sl.level)) - 1
	for h >= 0 {
		if x == nil {
			return nil, false
		}
		next := x.getNext(h)
		cmp := 1
		if next != nil {
			cmp = sl.cmp.Compare(next.key, key)
		}
		if cmp < 0 {
			x = next
		} else {
			if flag {
				element[h] = x
				sl.previousNodeCache[h] = element[h]
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

func (sl *SkipList) findLT(key interface{}) (*Element, bool) {
	x := sl.header
	h := int(atomic.LoadInt32(&sl.level)) - 1
	for h >= 0 {
		next := x.getNext(h)
		if next == nil || sl.cmp.Compare(next.key, key) >= 0 {
			if h == 0 {
				if x == sl.header {
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

func (sl *SkipList) randLevel() int32 {
	var l int32 = 1
	for ((sl.randSource.Int63() >> 32) & 0xFFFF) < DEFAULT_PROBABILITY {
		l++
	}
	if l > DEFAULT_MAX_LEVEL {
		l = DEFAULT_MAX_LEVEL
	}
	return l
}

func (sl *SkipList) Iterator() *SkipListIterator {
	return NewSkipListIterator(sl.header, sl.cmp)
}
