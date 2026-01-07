package parser

import (
	"ego/internal/frontend/ast"
)

func (p *Parser) parseAnonymousFunction() ast.Expression {
	p.stackTrace.Push(ANON_FUNCTION)
	anon := &ast.AnonymousFunction{Token: p.curToken}

	anon.Parameters = p.parseFunctionParameters()

	if p.errorExist() {
		return nil
	}

	anon.Body = p.parseBlockStatement()

	if p.errorExist() {
		return nil
	}

	p.stackTrace.Pop()
	return anon
}
