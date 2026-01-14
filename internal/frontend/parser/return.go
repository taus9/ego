package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	p.stackTrace.Push(RETURN)
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	if p.curTokenIs(token.EOL) {
		stmt.ReturnValue = nil
		return stmt
	}

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.errorExist() {
		return nil
	}

	if p.peekTokenIs(token.EOL) {
		p.nextToken()
	}

	p.stackTrace.Pop()
	return stmt
}
