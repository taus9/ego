package parser

import (
	"ego/ast"
	"ego/token"
)

func (p *Parser) parseAssignStatement() *ast.AssignStatement {
	stmt := &ast.AssignStatement{Token: token.Token{Type: token.ASSIGN, Literal: "="}}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return stmt
}
