package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

func (p *Parser) parseAnonymousFunction() ast.Expression {
	p.stackTrace.Push(ANON_FUNCTION)
	anon := &ast.AnonymousFunction{Token: p.curToken}

	p.nextToken() // consume ':>' token

	anon.Parameters = p.parseFunctionParameters()

	if p.errorExist() {
		return nil
	}

	p.blockStack.Push(token.ANON_FUNCTION)
	anon.Body = p.parseBlockStatement()
	p.blockStack.Pop()

	if p.errorExist() {
		return nil
	}

	//END BLOCK token already consumed in parseBlockStatement
	p.stackTrace.Pop()
	return anon
}
