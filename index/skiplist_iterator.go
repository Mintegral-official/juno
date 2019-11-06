package index

import (
	"github.com/Mintegral-official/juno/helpers"
)

type SkipListIterator struct {
	*SkipList
	index   int64
	element *Element
}

func NewSkipListIterator(level int32, cmp helpers.Comparable) *SkipListIterator {
	if level <= 0 || level > DEFAULT_MAX_LEVEL {
		level = DEFAULT_MAX_LEVEL
	}
	return &SkipListIterator{
		SkipList: NewSkipList(DEFAULT_MAX_LEVEL, cmp),
		index:    0,
		element: nil,
	}
}

func (slIterator *SkipListIterator) Valid() bool {
	return slIterator.element != nil
}

func (slIterator *SkipListIterator) First() bool {
	slIterator.element = slIterator.header.getNext(0)
	return slIterator.Valid()
}

func (slIterator *SkipListIterator) Iterator() InvertedIterator {
	if slIterator != nil {
		return slIterator
	}
	return nil
}

func (slIterator *SkipListIterator) HasNext() bool {

	if slIterator.element == nil {
		return slIterator.First()
	}
	slIterator.element = slIterator.element.getNext(0)
	return slIterator.Valid()
}

func (slIterator *SkipListIterator) Next() *Element {
	if slIterator.element == nil {
		return nil
	}
	v := slIterator.element.getNext(0)
	slIterator.index++
	return v
}

func (slIterator *SkipListIterator) GetGE(key interface{}) interface{} {

	if slIterator.index == slIterator.Len() {
		return nil
	}
	prev := slIterator.Next()
	if prev == nil || prev.getNext(0) == nil {
		return nil
	}
	k := prev.key
	//fmt.Println(k, prev)
	if slIterator.cmp.Compare(k, key) > 0 {
		return nil
	} else {
		for {
			if prev, ok := slIterator.findGE(key, true, slIterator.previousNodeCache); ok {
				//fmt.Println(prev, k)
				if slIterator.cmp.Compare(prev.key, key) < 0 {
					k = prev.key
				} else {
					return prev
				}
			} else {
				return nil
			}
		}
	}
}
