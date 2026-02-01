package ast

import (
	"bytes"
	"ego/internal/frontend/token"
	"strings"
)

type ExecLiteral struct {
	Token   token.Token // the 'exec' token
	Command string
}

func (e *ExecLiteral) expressionNode() {}

func (e *ExecLiteral) TokenLiteral() string {
	return e.Token.Literal
}

func (e *ExecLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("exec ")
	out.WriteString(strings.TrimSpace(e.Command))

	return out.String()
}
