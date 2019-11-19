package datastruct

import (
	"sync/atomic"
	"unsafe"
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

func ElementCopy(element *Element) *Element {
	if element == nil {
		return nil
	}
	var e = newNode(nil, nil, int32(len(element.next)))
	e.key = element.key
	e.value = element.value
	for i, v := range element.next {
		e.next[i] = v
	}
	return e
}