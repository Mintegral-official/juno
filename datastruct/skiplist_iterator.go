package datastruct

import (
	"github.com/Mintegral-official/juno/document"
)

type SkipListIterator struct {
	Element *Element
}

func NewSkipListIterator(element *Element) *SkipListIterator {
	sli := &SkipListIterator{element}
	sli.Next()
	return sli
}

func (si *SkipListIterator) HasNext() bool {
	return si.Element != nil
}

func (si *SkipListIterator) Next() {
	if si.Element == nil {
		return
	}
	next := si.Element.Next(0)
	if next == nil {
		si.Element = nil
		return
	}
	for i, v := range next.next {
		si.Element.next[i] = v
	}
	si.Element.key, si.Element.value = next.key, next.value
}

func (si *SkipListIterator) GetLE(key document.DocId) *Element {
	if si.Element == nil {
		return nil
	}
	for i := len(si.Element.next) - 1; i >= 0; {
		next := si.Element.Next(i)
		if next == nil {
			i--
			continue
		}
		cmp := int(key - next.key)
		if cmp == 0 {
			for ; i >= 0; i-- {
				si.Element.next[i] = next.next[i]
			}
			si.Element.key, si.Element.value = next.key, next.value
			return si.Element
		} else if cmp > 0 {
			si.Element.next[i] = next.next[i]
		} else {
			i--
		}
	}
	return si.Element
}

func (si *SkipListIterator) GetGE(key document.DocId) *Element {
	e := si.GetLE(key)
	if e != nil {
		if int(key-e.key) > 0 {
			si.Next()
		}
		return si.Element
	}
	return nil
}

func (si *SkipListIterator) Current() *Element {
	return si.Element
}
