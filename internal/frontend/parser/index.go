package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
	"fmt"
)

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expression := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	expression.Index = p.parseExpression(LOWEST)
	if p.errorsExist() {
		return nil
	}

	if !p.expectPeek(token.RBRACKET) {
		p.createErrorMessage(fmt.Sprintf("expected ']' after index expression, got %s instead", p.peekToken.Literal))
		return nil
	}

	return expression
}
