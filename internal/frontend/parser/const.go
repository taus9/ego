package parser

import (
	"ego/internal/frontend/ast"
)

func (p *Parser) parseConstant() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal, Mutable: false}
}
