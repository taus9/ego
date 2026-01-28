package ast

import "ego/internal/frontend/token"

type MemberExpression struct {
	Token    token.Token
	Object   *Identifier
	Property *Identifier
}
