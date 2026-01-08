package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
	"fmt"
)

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expression := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken() // consume '['

	expression.Index = p.parseExpression(LOWEST)

	if p.errorExist() {
		return nil
	}

	p.nextToken() // consume token

	if !p.curTokenIs(token.RBRACKET) {
		p.createErrorMessage(fmt.Sprintf("expected ']' after index expression, got %s instead", p.peekToken.Literal))
		return nil
	}

	p.nextToken() // consume ']'
	return expression
}
