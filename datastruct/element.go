package datastruct

import (
	"github.com/Mintegral-official/juno/document"
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

func (element *Element) setNext(n int, x *Element) {
	if element == nil {
		return
	}
	if n < 0 || n > len(element.next) {
		return
	}
	element.next[n] = unsafe.Pointer(x)
}

func (element *Element) Next(n int) *Element {
	if element == nil {
		return nil
	}
	if n < 0 || n > len(element.next) {
		return nil
	}
	return (*Element)(element.next[n])
}

func (element *Element) Key() document.DocId {
	if element == nil {
		return 0
	}
	return element.key
}

func (element *Element) Value() interface{} {
	if element == nil {
		return nil
	}
	return element.value
}

func ElementCopy(element *Element) *Element {
	if element == nil {
		return nil
	}
	var e = newNode(0, nil, len(element.next))
	e.key = element.key
	e.value = element.value
	for i, v := range element.next {
		e.next[i] = v
	}
	return e
}
