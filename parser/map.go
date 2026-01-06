package parser

import (
	"ego/ast"
	"ego/token"
	"fmt"
)

func (p *Parser) parseMapLiteral() ast.Expression {
	mapLiteral := &ast.MapLiteral{Token: p.curToken}
	mapLiteral.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			msg := fmt.Sprintf("expected ':' after map key, got %s instead", p.peekToken.Type)
			p.errors = append(p.errors, msg)
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		mapLiteral.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			msg := fmt.Sprintf("expected ',' or '}' after map pair, got %s instead", p.peekToken.Type)
			p.errors = append(p.errors, msg)
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		msg := fmt.Sprintf("expected '}' at end of map literal, got %s instead", p.peekToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	return mapLiteral
}
