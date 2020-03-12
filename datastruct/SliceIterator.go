package datastruct

import (
	"github.com/MintegralTech/juno/document"
)

type SliceIterator struct {
	data      *Slice
	index     int
	FieldName string
}

func (s *SliceIterator) HasNext() bool {
	return s.index < s.data.Len()
}

func (s *SliceIterator) Next() {
	if s.data.Len() == 0 {
		return
	}
	s.index++
}

func (s *SliceIterator) Current() *Element {
	if s.index >= s.data.Len() {
		return nil
	}
	return (*s.data)[s.index]
}

func (s *SliceIterator) GetGE(id document.DocId) *Element {
	idx := binarySearch1(s.data, id)
	if idx == -1 {
		*s.data = (*s.data)[0:0]
		s.index = s.data.Len()
		return nil

	} else if idx == 0 {
		return (*s.data)[0]
	} else {
		res := (*s.data)[idx]
		*s.data = (*s.data)[idx:]
		s.index = 0
		return res
	}
}

func binarySearch1(data *Slice, id document.DocId) int {
	left, right, mid := 0, data.Len()-1, 0
	if data.Len() == 0 {
		return -1
	}
	if (*data)[0].key >= id {
		return 0
	}
	for left < right {
		mid = (left + right) / 2
		if (*data)[mid].key > id {
			right = mid
		} else if (*data)[mid].key < id {
			left = mid + 1
		} else if (*data)[mid].key == id {
			return mid
		}
	}
	if (*data)[right].key >= id {
		return right
	} else {
		return -1
	}
}
