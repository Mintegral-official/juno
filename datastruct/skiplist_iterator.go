package datastruct

import (
	"github.com/Mintegral-official/juno/helpers"
)

type SkipListIterator struct {
	*SkipList
	Index   int64
	Element *Element
}

func NewSkipListIterator(level int32, cmp helpers.Comparable) *SkipListIterator {
	if level <= 0 || level > DEFAULT_MAX_LEVEL {
		level = DEFAULT_MAX_LEVEL
	}
	return &SkipListIterator{
		SkipList: NewSkipList(DEFAULT_MAX_LEVEL, cmp),
		Index:    0,
		Element:  nil,
	}
}

func (slIterator *SkipListIterator) Valid() bool {
	return slIterator.Element != nil
}

func (slIterator *SkipListIterator) First() bool {
	slIterator.Element = slIterator.header.getNext(0)
	return slIterator.Valid()
}

func (slIterator *SkipListIterator) Iterator() Iterator {
	if slIterator != nil {
		slIterator.Index = 0
		slIterator.Element = nil
		return slIterator
	}
	return nil
}

func (slIterator *SkipListIterator) HasNext() bool {

	if slIterator.Element == nil && slIterator.Index == 0 {
		return slIterator.First()
	}
	// slIterator.element = slIterator.element.Next(0)
	return slIterator.Valid()
}

func (slIterator *SkipListIterator) Next() interface{} {
	//_ = slIterator.Element
	if slIterator.Element == nil {
		if slIterator.Index == 0 {
			slIterator.Element = slIterator.header.getNext(0)
			slIterator.Index++
			if slIterator.Element != nil {
				slIterator.Element = slIterator.Element.getNext(0)
			}
			return slIterator.header.getNext(0)
		}
		return nil
	}
	v := slIterator.Element
	slIterator.Element = slIterator.Element.getNext(0)
	slIterator.Index++
	return v
}

func (slIterator *SkipListIterator) GetGE(key interface{}) interface{} {
	var prev interface{}
	if slIterator.Index == slIterator.length {
		return nil
	} else if slIterator.Element == nil && slIterator.Index == 0 {
		slIterator.Element = slIterator.header.getNext(0)
		prev = slIterator.Element
	} else {
		prev = slIterator.Element
	}
	//fmt.Println(prev)
	if prev == nil {
		return nil
	}
	k := prev.(*Element).key
	//fmt.Println(k, prev)
	if slIterator.cmp.Compare(k, key) > 0 {
		return prev
	} else {
		for {
			if prev, ok := slIterator.findGE(key, true, slIterator.previousNodeCache); ok {
				//fmt.Println(prev, k)
				if slIterator.cmp.Compare(prev.key, key) < 0 {
					k = prev.key
				} else {
					return prev
				}
			} else if prev != nil {
				return prev
			} else {
				return nil
			}
		}
	}
}

func (slIterator *SkipListIterator) Current() (interface{}, error) {
	if slIterator.Element == nil && slIterator.Index == 0 {
		return slIterator.header.getNext(0).Key(), nil
	} else if slIterator.Element == nil {
		return 0, helpers.ElementNotfound
	}
	return slIterator.Element.Key(), nil
}
