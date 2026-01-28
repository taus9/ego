package ast

import (
	"bytes"
	"ego/internal/frontend/token"
)

type UseStatement struct {
	Token  token.Token // the 'use' token
	Module *StringLiteral
	Alias  *Identifier // optional alias
}

func (us *UseStatement) statementNode()       {}
func (us *UseStatement) TokenLiteral() string { return us.Token.Literal }
func (us *UseStatement) String() string {
	var out bytes.Buffer
	out.WriteString(us.TokenLiteral() + " ")
	out.WriteString(us.Module.String())
	if us.Alias != nil {
		out.WriteString(" as ")
		out.WriteString(us.Alias.String())
	}
	out.WriteString(";")
	return out.String()
}
