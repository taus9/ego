package parser

import (
	"ego/internal/frontend/ast"
)

func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	return &ast.BreakStatement{Token: p.curToken}
}

func (p *Parser) parseAgainStatement() *ast.AgainStatement {
	return &ast.AgainStatement{Token: p.curToken}
}
