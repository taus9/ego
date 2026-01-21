package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

func (p *Parser) parseMapLiteral() ast.Expression {
	p.stackTrace.Push(MAP)
	mapLiteral := &ast.MapLiteral{Token: p.curToken}
	mapLiteral.Pairs = make(map[ast.Expression]ast.Expression)

	p.nextToken() // consuming '{'

	hasEOLs := false

	for !p.curTokenIs(token.RBRACE) {

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

		value := p.parseExpression(LOWEST)

		if p.errorExist() {
			return nil
		}

		mapLiteral.Pairs[key] = value
		p.nextToken()

		// at this point multiple scenarios are possible
		// 1. we have a comma with more pairs to come
		// 2. we have a trailing comma followed by a closing brace
		// 3. we have a comma, then EOLs, then more pairs
		// 4. we have a trailing comma, then EOLs, then closing brace
		// 5. we have EOLs, then closing brace
		// 6. we have a closing brace right away

		// I don't want to allow EOLs then a comma
		// there might even be more but this is all I could think of

		for p.curTokenIs(token.EOL) {
			p.nextToken()
			hasEOLs = true
		}

		if p.curTokenIs(token.COMMA) {
			if hasEOLs {
				p.createErrorMessage("seperator must be the same line as map pair")
				return nil
			}

			p.nextToken() // consume comma

			for p.curTokenIs(token.EOL) {
				p.nextToken()
			}

			if p.curTokenIs(token.RBRACE) {
				break
			}

			//p.nextToken() // consume comma
			hasEOLs = false
			continue
		}

		if p.curTokenIs(token.RBRACE) {
			break
		}

		p.createErrorMessage("expected ',' or '}' after map pair")
		return nil
	}

	p.stackTrace.Pop()
	return mapLiteral
}
