package index

type SkipList struct {
}

func NewSkipList(level int) *SkipList {
	return &SkipList{}
}

func Add(id DocId) {

}

func Del(id DocId) {

}

func Contains(id DocId) bool {
	return false
}

func (sl *SkipList) Iterator() InvertedIterator {
	return nil
}
