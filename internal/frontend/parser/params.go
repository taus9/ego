package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	// no parameters to parse
	if p.curTokenIs(token.EOL) {
		p.nextToken() // consume EOL token
		return identifiers
	}

	if !p.curTokenIs(token.IDENT) {
		p.createErrorMessage("expected function parameter identifier")
		return nil
	}

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)
	p.nextToken() // consume identifier

	for p.curTokenIs(token.COMMA) {
		p.nextToken() // consume ',' token

		if !p.curTokenIs(token.IDENT) {
			p.createErrorMessage("expected function parameter identifier after ','")
			return nil
		}

		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)

		p.nextToken() // consume identifier
	}

	if !p.curTokenIs(token.EOL) {
		p.createErrorMessage("expected EOL after function parameters")
		return nil
	}

	p.nextToken() // consume EOL token
	return identifiers
}
