package index

import (
	"github.com/Mintegral-official/juno/helpers"
	"math/rand"
	"time"
)

const (
	DEFAULT_MAX_LEVEL   = 12
	DEFAULT_PROBABILITY = 0x3FFF
)

type elementNode struct {
	next []*Element
}

type Element struct {
	elementNode
	key, value interface{}
	score      float64
}

func (element *Element) Key() interface{} {
	return element.key
}

func (element *Element) Next() *Element {
	return element.next[0]
}

func (element *Element) Value() interface{} {
	return element.value
}

type SkipList struct {
	elementNode
	level          int
	length         int
	keyFunc        helpers.Comparable
	randSource     rand.Source
	prevNodesCache []*elementNode
}

func (skipList *SkipList) Front() *Element {
	return skipList.next[0]
}

func (skipList *SkipList) Len() int {
	return skipList.length
}

func NewSkipList(level int, keyFunc helpers.Comparable) *SkipList {
	if level <= 0 || level > DEFAULT_MAX_LEVEL {
		level = DEFAULT_MAX_LEVEL
	}
	return &SkipList{
		elementNode:    elementNode{next: make([]*Element, DEFAULT_MAX_LEVEL)},
		prevNodesCache: make([]*elementNode, level),
		level:          DEFAULT_MAX_LEVEL,
		keyFunc:        keyFunc,
		randSource:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (skipList *SkipList) Compare(a, b interface{}) int {
	// TODO
	return 0
}

func (skipList *SkipList) randLevel() int {
	l := 1

	for ((skipList.randSource.Int63() >> 32) & 0xFFFF) < DEFAULT_PROBABILITY {
		l++
	}

	if l > skipList.level {
		l = skipList.level
	}

	return l
}

func (skipList *SkipList) findPreElements(key interface{}) []*elementNode {
	prev := &skipList.elementNode
	var next *Element

	prevs := skipList.prevNodesCache

	for i := skipList.level - 1; i >= 0; i-- {
		next = prev.next[i]

		for next != nil && skipList.keyFunc.Compare(next.key, key) < 0 {
			prev = &next.elementNode
			next = next.next[i]
		}

		prevs[i] = prev
	}

	return prevs
}

func (skipList *SkipList) Add(key, value interface{}) *Element {
	var element *Element

	prevs := skipList.findPreElements(key)

	if element = prevs[0].next[0]; element != nil && skipList.keyFunc.Compare(element.key, key) <= 0 {
		element.value = value
		return element
	}

	element = &Element{
		elementNode: elementNode{
			next: make([]*Element, skipList.randLevel()),
		},
		key:   key,
		value: value,
	}
	// fmt.Println(skipList.level)

	for i := range element.next {
		element.next[i] = prevs[i].next[i]
		prevs[i].next[i] = element
	}

	skipList.length++
	return element
}

func (skipList *SkipList) GetK(key interface{}) *Element {
	prev := &skipList.elementNode
	var next *Element

	for i := skipList.level - 1; i >= 0; i-- {
		next = prev.next[i]

		for next != nil && skipList.keyFunc.Compare(next.key, key) < 0 {
			prev = &next.elementNode
			next = next.next[i]
		}
	}

	if next != nil && skipList.keyFunc.Compare(next.key, key) == 0 {
		return next
	}

	return nil
}

func (skipList *SkipList) GetV(key interface{}) (interface{}, bool) {
	element := skipList.GetK(key)
	if element != nil {
		return element.value, true
	}
	return nil, false
}

func (skipList *SkipList) Contains(key interface{}) bool {
	element := skipList.GetK(key)
	if element == nil {
		return false
	}
	return true
}

func (skipList *SkipList) Del(key interface{}) *Element {

	prevs := skipList.findPreElements(key)
	if prevs == nil {
		return nil
	}

	if element := prevs[0].next[0]; element != nil && skipList.keyFunc.Compare(element.key, key) == 0 {
		for k, v := range element.next {
			prevs[k].next[k] = v
		}
		skipList.length--
		return element
	}

	return nil
}

//import "github.com/Mintegral-official/juno/document"
//
//type SkipList struct {
//}
//
//func NewSkipList(level int) *SkipList {
//	return &SkipList{}
//}
//
//func (sl *SkipList) Add(id document.DocId) {
//
//}
//
//func (sl *SkipList) Del(id document.DocId) {
//
//}
//
//func (sl *SkipList) Contains(id document.DocId) bool {
//	return false
//}
//
//func (sl *SkipList) Iterator() InvertedIterator {
//	return nil
//}
