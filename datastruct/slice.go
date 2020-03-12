package datastruct

import (
	"errors"
	"github.com/MintegralTech/juno/document"
)

type Slice []*Element

func (s Slice) Iterator() *SliceIterator {
	return &SliceIterator{index: 0, data: &s}
}

func (s Slice) Len() int {
	return len(s)
}

func NewSlice() *Slice {
	return &Slice{}
}

func (s *Slice) Add(id document.DocId, value interface{}) {
	if s.Len() == 0 {
		*s = append(*s, &Element{key: id, value: value})
	} else {
		for i := 0; i < s.Len(); i++ {
			if (*s)[i].key == id {
				(*s)[i].value = value
				break
			} else if (*s)[i].key > id {
				tmp := append([]*Element{}, (*s)[i:]...)
				*s = append((*s)[0:i], &Element{key: id, value: value})
				*s = append(*s, tmp...)
				break
			}
		}
		if (*s)[s.Len()-1].key < id {
			*s = append(*s, &Element{key: id, value: value})
		}
	}
}

func (s *Slice) Get(id document.DocId) (*Element, error) {
	for _, v := range *s {
		if v.key == id {
			return v, nil
		}
	}
	return nil, errors.New("not found")
}
