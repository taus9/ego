package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

func (p *Parser) parseUseStatement() *ast.UseStatement {
	p.stackTrace.Push(USE)
	stmt := &ast.UseStatement{Token: p.curToken}

	p.nextToken() // consume use

	if !p.curTokenIs(token.STRING) {
		p.createErrorMessage("STRING expected after USE")
		return nil
	}

	stmt.Module = &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.AS) {
		p.nextToken() // consume module name
		p.nextToken() // consume as
		if !p.curTokenIs(token.IDENT) {
			p.createErrorMessage("expected IDENTIFIER after AS")
			return nil
		}

		stmt.Alias = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal, Mutable: false}
	}

	if !p.peekTokenIs(token.EOL) {
		p.createErrorMessage("expected EOL after USE statement")
		return nil
	}

	p.nextToken() // advance to EOL
	p.stackTrace.Pop()
	return stmt
}
