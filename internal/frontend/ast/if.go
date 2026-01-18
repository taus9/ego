package ast

import (
	"bytes"
	"ego/internal/frontend/token"
)

type IfBlockExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression  // The condition expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ibe *IfBlockExpression) expressionNode() {}

func (ibe *IfBlockExpression) TokenLiteral() string {
	return ibe.Token.Literal
}

func (ibe *IfBlockExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ibe.Condition.String())
	out.WriteString(" ")
	out.WriteString(ibe.Consequence.String())

	if ibe.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ibe.Alternative.String())
	}

	return out.String()
}

type IfTernaryExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression  // The condition expression
	Consequence Expression
	Alternative Expression
}

func (ite *IfTernaryExpression) expressionNode() {}

func (ite *IfTernaryExpression) TokenLiteral() string {
	return ite.Token.Literal
}

func (ite *IfTernaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ite.Condition.String())
	out.WriteString(" ")
	out.WriteString(ite.Consequence.String())
	out.WriteString(" else ")
	out.WriteString(ite.Alternative.String())

	return out.String()
}
