package ast

import (
	"bytes"
	"ego/internal/frontend/token"
)

type ElseInlineExpression struct {
	Token          token.Token // the 'else' token
	TryExpression  Expression
	ElseExpression Expression
}

func (eie *ElseInlineExpression) expressionNode() {}

func (eie *ElseInlineExpression) TokenLiteral() string {
	return eie.Token.Literal
}

func (eie *ElseInlineExpression) String() string {
	var out bytes.Buffer

	out.WriteString(eie.TryExpression.String())
	out.WriteString(" else ")
	out.WriteString(eie.ElseExpression.String())

	return out.String()
}

type ElseBlockExpression struct {
	Token         token.Token // the 'else' token
	TryExpression Expression
	ElseBlock     *BlockStatement
}

func (ebe *ElseBlockExpression) expressionNode() {}

func (ebe *ElseBlockExpression) TokenLiteral() string {
	return ebe.Token.Literal
}

func (ebe *ElseBlockExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ebe.TryExpression.String())
	out.WriteString(" else ")
	out.WriteString(ebe.ElseBlock.String())

	return out.String()
}
