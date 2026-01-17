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

type ForWhileStatement struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (fws *ForWhileStatement) statementNode() {}

func (fws *ForWhileStatement) TokenLiteral() string {
	return fws.Token.Literal
}

func (fws *ForWhileStatement) String() string {
	var out bytes.Buffer

	out.WriteString(fws.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(fws.Condition.String())
	out.WriteString(" ")
	out.WriteString(fws.Body.String())

	return out.String()
}

type ForInStatement struct {
	Token    token.Token
	Index    *Identifier
	Value    *Identifier
	Iterable Expression
	Body     *BlockStatement
}

func (fis *ForInStatement) statementNode() {}

func (fis *ForInStatement) TokenLiteral() string {
	return fis.Token.Literal
}

func (fis *ForInStatement) String() string {
	var out bytes.Buffer

	out.WriteString(fis.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(fis.Index.String())
	out.WriteString(", ")
	out.WriteString(fis.Value.String())
	out.WriteString(" in ")
	out.WriteString(fis.Iterable.String())
	out.WriteString(" ")
	out.WriteString(fis.Body.String())

	return out.String()
}
