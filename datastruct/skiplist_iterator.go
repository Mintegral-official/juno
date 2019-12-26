package datastruct

import (
	"github.com/Mintegral-official/juno/document"
)

type SkipListIterator struct {
	Element *Element
	//cmp     helpers.Comparable
}

func NewSkipListIterator(element *Element) *SkipListIterator {
	sli := &SkipListIterator{element}
	sli.Next()
	return sli
}

func (slIterator *SkipListIterator) HasNext() bool {
	return slIterator.Element != nil
}

func (slIterator *SkipListIterator) Next() {
	if slIterator.Element == nil {
		return
	}
	next := slIterator.Element.Next(0)
	if next == nil {
		//res := slIterator.Element
		slIterator.Element = nil
		return
	}
	for i, v := range next.next {
		slIterator.Element.next[i] = v
	}
	slIterator.Element.key, slIterator.Element.value = next.key, next.value
	//return slIterator.Element
}

func (slIterator *SkipListIterator) GetLE(key document.DocId) interface{} {
	for i := len(slIterator.Element.next) - 1; i >= 0; {
		next := slIterator.Element.Next(i)
		if next == nil {
			i--
			continue
		}
		//	cmp := slIterator.cmp.Compare(key, next.key)
		cmp := int(key - next.key)
		if cmp == 0 {
			for ; i >= 0; i-- {
				slIterator.Element.next[i] = next.next[i]
			}
			slIterator.Element.key, slIterator.Element.value = next.key, next.value
			return slIterator.Element
		} else if cmp > 0 {
			slIterator.Element.next[i] = next.next[i]
		} else {
			i--
		}
	}
	return slIterator.Element
}

func (slIterator *SkipListIterator) GetGE(key document.DocId) interface{} {
	e := slIterator.GetLE(key).(*Element)
	if e == nil {
		return nil
	}
	c := int(key - e.key)
	if c > 0 {
		slIterator.Next()
	}
	return slIterator.Element
}

func (slIterator *SkipListIterator) Current() interface{} {
	return slIterator.Element
}
