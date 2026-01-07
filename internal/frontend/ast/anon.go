package ast

import (
	"bytes"
	"ego/internal/frontend/token"
	"strings"
)

type AnonymousFunction struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (af *AnonymousFunction) expressionNode() {}

func (af *AnonymousFunction) TokenLiteral() string {
	return af.Token.Literal
}

func (af *AnonymousFunction) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range af.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(af.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(af.Body.String())

	return out.String()
}
