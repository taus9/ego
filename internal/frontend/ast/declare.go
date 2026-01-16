package ast

import (
	"bytes"
	"ego/internal/frontend/token"
)

type DeclareStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *DeclareStatement) statementNode() {}

func (ls *DeclareStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *DeclareStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	return out.String()
}
