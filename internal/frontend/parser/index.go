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

	if !p.expectPeek(token.RBRACKET) {
		msg := fmt.Sprintf("expected ']' after index expression, got %s instead", p.peekToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	return expression
}
