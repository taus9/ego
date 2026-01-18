package ast

import "ego/internal/frontend/token"

type Constant struct {
	Token token.Token // the token.CONST token
	Value string
}

func (c *Constant) expressionNode()      {}
func (c *Constant) TokenLiteral() string { return c.Token.Literal }
func (c *Constant) String() string       { return c.Token.Literal }
