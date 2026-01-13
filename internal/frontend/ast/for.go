package ast

import (
	"bytes"
	"ego/internal/frontend/token"
)

type ForToStatement struct {
	Token    token.Token
	Iterator *Identifier
	Start    Expression
	End      Expression
	Body     *BlockStatement
}

func (fts *ForToStatement) statementNode() {}

func (fts *ForToStatement) TokenLiteral() string {
	return fts.Token.Literal
}

func (fts *ForToStatement) String() string {
	var out bytes.Buffer

	out.WriteString(fts.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(fts.Iterator.String())
	out.WriteString(" = ")
	out.WriteString(fts.Start.String())
	out.WriteString(" to ")
	out.WriteString(fts.End.String())
	out.WriteString(" ")
	out.WriteString(fts.Body.String())

	return out.String()
}
