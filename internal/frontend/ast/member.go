package ast

import (
	"bytes"
	"ego/internal/frontend/token"
)

type MemberExpression struct {
	Token          token.Token
	ModuleName     *Identifier
	ModuleProperty Expression
}

func (me *MemberExpression) expressionNode() {}

func (me *MemberExpression) TokenLiteral() string {
	return me.Token.Literal
}

func (me *MemberExpression) String() string {
	var out bytes.Buffer

	out.WriteString(me.ModuleName.String())
	out.WriteString(".")
	out.WriteString(me.ModuleProperty.String())

	return out.String()
}
