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