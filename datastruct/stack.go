package datastruct

type node struct {
	value interface{}
	next  *node
}

type Stack struct {
	top *node
	len int
}

func NewStack() *Stack {
	return &Stack{
		top: nil,
		len: 0,
	}
}

func (s *Stack) Empty() bool {
	return s.len == 0
}

func (s *Stack) Len() int {
	return s.len
}

func (s *Stack) Push(value interface{}) {
	n := &node{
		value: value,
		next:  s.top,
	}
	s.top = n
	s.len++
}

func (s *Stack) Pop() interface{} {
	if s.len == 0 {
		return nil
	}
	top := s.top
	s.top = top.next
	s.len--
	return top.value
}

func (s *Stack) Peek() interface{} {
	if s.len == 0 {
		return nil
	}
	return s.top.value
}
