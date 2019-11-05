package index

import (
	"github.com/Mintegral-official/juno/helpers"
)

type SkipListIterator struct {
	*SkipList
	index   int64
	element *Element
}

func NewSkipListIterator(level int32, keyFunc helpers.Comparable) *SkipListIterator {
	if level <= 0 || level > DEFAULT_MAX_LEVEL {
		level = DEFAULT_MAX_LEVEL
	}
	return &SkipListIterator{
		SkipList: NewSkipList(DEFAULT_MAX_LEVEL, keyFunc),
		index:    0,
		element:  newNode(nil, nil, DEFAULT_MAX_LEVEL),
	}
}

func (slIterator *SkipListIterator) Iterator() InvertedIterator {
	if slIterator != nil {
		return slIterator
	}
	return nil
}

func (slIterator *SkipListIterator) HasNext() bool {

	if slIterator == nil {
		return false
	}
	return slIterator.header.next[0] != nil
}

func (slIterator *SkipListIterator) Next() *Element {
	if slIterator == nil {
		return nil
	}

	v := slIterator.header.next[0]
	if slIterator.header.getNext(0) == nil {
		return nil
	}
	slIterator.header.next[0] = slIterator.header.getNext(0).next[0]
	slIterator.index++
	return (*Element)(v)

}

func (slIterator *SkipListIterator) GetGE(key interface{}) interface{} {

	if slIterator == nil {
		return nil
	}

	prev := slIterator.header
	for i := int(slIterator.level) - 1; i >= 0; i-- {
		for {
			if prev.next == nil || prev.next[i] == nil || prev.getNext(i).key == nil {
				break
			}
			if slIterator.cmp.Compare(prev.getNext(i).key, key) == 0 {
				return prev.getNext(i).value
			}

			if slIterator.cmp.Compare(prev.getNext(i).key, key) < 0 {
				prev = prev.getNext(i)
				continue
			} else {
				//	i--
				break
			}
		}
	}
	for {
		if prev.next == nil || prev.next[0] == nil {
			return nil
		} else if slIterator.cmp.Compare(prev.getNext(0).key, key) < 0 {
			prev = prev.getNext(0)
		} else {
			return prev.getNext(0).value
		}
	}
}
