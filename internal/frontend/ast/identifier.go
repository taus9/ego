package ast

import "ego/internal/frontend/token"

type Identifier struct {
	Token   token.Token
	Value   string
	Mutable bool
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}
