package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	if p.curTokenIs(token.RPAREN) {
		p.createErrorMessage("expected expression")
		return nil
	}

	exp := p.parseExpression(LOWEST)

	if p.errorsExist() {
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		p.createErrorMessage("expected closing parenthesis for grouped expression")
		return nil
	}

	return exp
}
