package ast

import (
	"bytes"
	"ego/internal/frontend/token"
)

type AssignStatement struct {
	Token token.Token // the token.IDENT token
	Name  *Identifier
	Index Expression // Optional index for array assignment
	Value Expression
}

func (as *AssignStatement) statementNode() {}

func (as *AssignStatement) TokenLiteral() string {
	return as.Token.Literal
}

func (as *AssignStatement) String() string {
	var out bytes.Buffer

	if as.Name != nil {
		out.WriteString(as.Name.String())
	} else if as.Index != nil {
		out.WriteString(as.Index.String())
	}
	out.WriteString(" = ")

	if as.Value != nil {
		out.WriteString(as.Value.String())
	}

	return out.String()
}
