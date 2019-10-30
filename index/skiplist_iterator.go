package index

import (
	"github.com/Mintegral-official/juno/helpers"
)

type SkipListIterator struct {
	*SkipList
	index int
}

func NewSKipListIterator(level int, keyFunc helpers.Comparable) *SkipListIterator {
	if level <= 0 || level > DEFAULT_MAX_LEVEL {
		level = DEFAULT_MAX_LEVEL
	}
	return &SkipListIterator{
		SkipList: NewSkipList(level, keyFunc),
		index:    0,
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
	return slIterator.index < slIterator.length
}

func (slIterator *SkipListIterator) Next() interface{} {
	if slIterator == nil {
		return 0
	}
	//fmt.Println(slIterator.length)
	slIterator.index++
	v := slIterator.elementNode.next[0].key
	slIterator.elementNode.next[0] = slIterator.elementNode.next[0].next[0]
	return v

}

func (slIterator *SkipListIterator) GetGE(key interface{}) interface{} {

	prev := &slIterator.elementNode
	for i := slIterator.level - 1; i >= 0; i-- {
		for {
			if prev.next == nil || prev.next[i] == nil || prev.next[i].key == nil {
				break
			}
			if slIterator.keyFunc.Compare(prev.next[i].key, key) == 0 {
				return prev.next[i].value
			}

			if slIterator.keyFunc.Compare(prev.next[i].key, key) < 0 {
				prev = &prev.next[i].elementNode
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
		} else if slIterator.keyFunc.Compare(prev.next[0].key, key) < 0 {
			prev = &prev.next[0].elementNode
		} else {
			return prev.next[0].value
		}
	}
}
