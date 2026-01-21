package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

func (p *Parser) parseDeclareStatement() *ast.DeclareStatement {
	p.stackTrace.Push(DECLARE)
	stmt := &ast.DeclareStatement{Token: token.Token{Type: token.DECLARE, Literal: token.DECLARE_INTERNAL}}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.curTokenIs(token.CONST) {
		stmt.Name.Mutable = false
	} else {
		stmt.Name.Mutable = true
	}

	p.nextToken() // consume indentifier
	p.nextToken() // consume ':='

	stmt.Value = p.parseExpression(LOWEST)

	if p.errorExist() {
		return nil
	}

	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // consume current expression
		elseExpr := p.parseElseExpression(stmt.Value)
		if p.errorExist() {
			return nil
		}

		stmt.Value = elseExpr
	}

	if p.peekTokenIs(token.EOL) {
		p.nextToken()
	}

	p.stackTrace.Pop()
	return stmt
}
