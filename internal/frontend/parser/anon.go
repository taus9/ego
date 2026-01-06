package parser

import (
	"ego/internal/frontend/ast"
)

func (p *Parser) parseFunctionLiteral() ast.Expression {
	anon := &ast.FunctionLiteral{Token: p.curToken}

	anon.Parameters = p.parseFunctionParameters()

	anon.Body = p.parseBlockStatement()

	return anon
}
