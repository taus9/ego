package ast

import "ego/internal/frontend/token"

type BreakStatement struct {
	Token token.Token // the 'break' token
}

func (bs *BreakStatement) statementNode()       {}
func (bs *BreakStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BreakStatement) String() string       { return bs.Token.Literal }

type AgainStatement struct {
	Token token.Token
}

func (as *AgainStatement) statementNode()       {}
func (as *AgainStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AgainStatement) String() string       { return as.Token.Literal }
