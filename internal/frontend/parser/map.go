package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
	"fmt"
)

func (p *Parser) parseMapLiteral() ast.Expression {
	mapLiteral := &ast.MapLiteral{Token: p.curToken}
	mapLiteral.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if p.errorsExist() {
			return nil
		}

		if !p.expectPeek(token.COLON) {
			p.createErrorMessage(fmt.Sprintf("expected ':' after map key, got %s instead", p.peekToken.Type))
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		if p.errorsExist() {
			return nil
		}

		mapLiteral.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			p.createErrorMessage(fmt.Sprintf("expected ',' or '}' after map pair, got %s instead", p.peekToken.Type))
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		p.createErrorMessage(fmt.Sprintf("expected '}' at end of map literal, got %s instead", p.peekToken.Type))
		return nil
	}

	return mapLiteral
}
