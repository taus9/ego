package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
	"fmt"
)

func (p *Parser) parseMapLiteral() ast.Expression {
	p.stackTrace.Push(MAP)
	mapLiteral := &ast.MapLiteral{Token: p.curToken}
	mapLiteral.Pairs = make(map[ast.Expression]ast.Expression)

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) {
		p.nextToken()
		for p.curTokenIs(token.EOL) {
			p.nextToken()
		}
		key := p.parseExpression(LOWEST)

		if p.errorExist() {
			return nil
		}

		p.nextToken()

		if !p.curTokenIs(token.COLON) {
			p.createErrorMessage("expected ':' after map key")
			return nil
		}

		p.nextToken() // consuming ':'

		for p.curTokenIs(token.EOL) {
			p.nextToken()
		}
		value := p.parseExpression(LOWEST)

		if p.errorExist() {
			return nil
		}

		mapLiteral.Pairs[key] = value
		p.nextToken()
		for p.curTokenIs(token.EOL) {
			p.nextToken()
		}

		if !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.COMMA) {
			p.createErrorMessage(fmt.Sprintf("expected ',' or '}' after map pair, got %s instead", p.curToken.Type))
			return nil
		}
	}

	p.stackTrace.Pop()
	return mapLiteral
}
