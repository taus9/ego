package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: token.Token{Type: token.LET_INTERNAL, Literal: token.LET_INTERNAL}}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// this is a redundant check, but added for clarity
	if !p.expectPeek(token.ASSIGN) {
		p.createErrorMessage("expected '=' after LET identifier")
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.EOL) {
		p.nextToken()
	}

	return stmt
}
