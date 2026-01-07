package parser

type Stack struct {
	elements []any
}

func NewStack() *Stack {
	return &Stack{
		elements: make([]any, 0),
	}
}

func (s *Stack) Push(t any) {
	s.elements = append(s.elements, t)
}

func (s *Stack) Pop() any {
	if len(s.elements) == 0 {
		return nil
	}

	n := len(s.elements) - 1
	e := s.elements[n]
	s.elements = s.elements[:n]
	return e
}

func (s *Stack) Peek() any {
	if len(s.elements) == 0 {
		return nil
	}
	return s.elements[len(s.elements)-1]
}

func (s *Stack) IsEmpty() bool {
	return len(s.elements) == 0
}

func (s *Stack) Size() int {
	return len(s.elements)
}
