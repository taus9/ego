package parser

import (
	"ego/internal/frontend/ast"
)

func (p *Parser) parseFunctionLiteral() ast.Expression {
	anon := &ast.FunctionLiteral{Token: p.curToken}

	anon.Parameters = p.parseFunctionParameters()

	if p.errorsExist() {
		return nil
	}

	anon.Body = p.parseBlockStatement()

	if p.errorsExist() {
		return nil
	}

	return anon
}
