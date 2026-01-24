package ast

import (
	"bytes"
	"ego/internal/frontend/token"
)

type WhenExpression struct {
	Token     token.Token // the 'when' token
	Selector  Expression
	Cases     []*WhenCase
	ElseBlock *BlockStatement
}

type WhenCase struct {
	Token     token.Token // the 'is' token
	Condition Expression
	Block     *BlockStatement
}

func (we *WhenExpression) expressionNode()      {}
func (we *WhenExpression) TokenLiteral() string { return we.Token.Literal }
func (we *WhenExpression) String() string {
	var out bytes.Buffer

	out.WriteString("when ")

	if we.Selector != nil {
		out.WriteString(we.Selector.String())
		out.WriteString("\n")
	} else {
		out.WriteString("\n")
	}

	for _, whenCase := range we.Cases {
		out.WriteString(whenCase.String())
	}

	if we.ElseBlock != nil {
		out.WriteString("else\n")
		out.WriteString(we.ElseBlock.String())
	}

	return out.String()
}

func (wc *WhenCase) statementNode()       {}
func (wc *WhenCase) TokenLiteral() string { return wc.Token.Literal }
func (wc *WhenCase) String() string {
	var out bytes.Buffer

	out.WriteString("is ")
	out.WriteString(wc.Condition.String())
	out.WriteString("\n")
	if wc.Block != nil {
		out.WriteString(wc.Block.String())
	}

	return out.String()
}
