package parser

import (
	"ego/ast"
	"ego/token"
)

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: token.Token{Type: token.LET_INTERNAL, Literal: token.LET_INTERNAL}}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	for !p.curTokenIs(token.EOL) {
		p.nextToken()
	}

	return stmt
}
