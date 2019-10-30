package index

import "github.com/Mintegral-official/juno/document"

type SkipList struct {
}

func NewSkipList(level int) *SkipList {
	return &SkipList{}
}

func (sl *SkipList) Add(id document.DocId) {

}

func (sl *SkipList) Del(id document.DocId) {

}

func (sl *SkipList) Contains(id document.DocId) bool {
	return false
}

func (sl *SkipList) Iterator() InvertedIterator {
	return nil
}
