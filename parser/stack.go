package parser

import "ego/token"

type BlockStack struct {
	elements []token.TokenType
}

func NewStack() *BlockStack {
	return &BlockStack{
		elements: make([]token.TokenType, 0),
	}
}

func (s *BlockStack) Push(t token.TokenType) {
	s.elements = append(s.elements, t)
}

func (s *BlockStack) Pop() token.TokenType {
	if len(s.elements) == 0 {
		return token.ILLEGAL
	}

	n := len(s.elements) - 1
	e := s.elements[n]
	s.elements = s.elements[:n]
	return e
}

func (s *BlockStack) Peek() token.TokenType {
	if len(s.elements) == 0 {
		return token.ILLEGAL
	}
	return s.elements[len(s.elements)-1]
}

func (s *BlockStack) IsEmpty() bool {
	return len(s.elements) == 0
}

func (s *BlockStack) Size() int {
	return len(s.elements)
}
