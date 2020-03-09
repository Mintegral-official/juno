package datastruct

import (
	"github.com/MintegralTech/juno/document"
	"unsafe"
)

type Element struct {
	key   document.DocId
	value interface{}
	next  []unsafe.Pointer
}

func newNode(key document.DocId, value interface{}, level int) *Element {
	if level <= 0 || level > DefaultMaxLevel {
		level = DefaultMaxLevel
	}
	return &Element{
		key:   key,
		value: value,
		next:  make([]unsafe.Pointer, level),
	}
}

func (e *Element) setNext(n int, x *Element) {
	if e == nil || n < 0 || n > len(e.next) {
		return
	}
	e.next[n] = unsafe.Pointer(x)
}

func (e *Element) Next(n int) *Element {
	if e == nil || n < 0 || n > len(e.next) {
		return nil
	}
	return (*Element)(e.next[n])
}

func (e *Element) Key() document.DocId {
	if e == nil {
		return 0
	}
	return e.key
}

func (e *Element) Value() interface{} {
	if e == nil {
		return nil
	}
	return e.value
}

func ElementCopy(element *Element) (e *Element) {
	if element == nil {
		return nil
	}
	e = newNode(0, nil, len(element.next))
	e.key, e.value = element.key, element.value
	for i, v := range element.next {
		e.next[i] = v
	}
	return
}
