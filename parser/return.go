package parser

import (
	"ego/ast"
	"ego/token"
)

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	for !p.curTokenIs(token.EOL) {
		p.nextToken()
	}

	return stmt
}
