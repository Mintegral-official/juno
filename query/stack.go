package query

type node struct {
	value string
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
	return s == nil || s.len == 0
}

func (s *Stack) Len() int {
	return s.len
}

func (s *Stack) Push(value string) {
	n := &node{
		value: value,
		next:  s.top,
	}
	s.top = n
	s.len++
}

func (s *Stack) Pop() string {
	if s.len == 0 {
		return ""
	}
	top := s.top
	s.top = top.next
	s.len--
	return top.value
}

func (s *Stack) Peek() string {
	if s.len == 0 {
		return ""
	}
	return s.top.value
}
