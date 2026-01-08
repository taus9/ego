package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	p.stackTrace.Push(CALL)
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	p.nextToken() // consume '('
	exp.Arguments = p.parseCallArguments()
	if p.errorExist() {
		return nil
	}
	// ')' has been consumed in parseCallArguments
	p.stackTrace.Pop()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	p.stackTrace.Push(CALL_ARGS)
	args := []ast.Expression{}

	if p.curTokenIs(token.RPAREN) {
		p.nextToken() // consume ')'
		p.stackTrace.Pop()
		return args
	}

	args = append(args, p.parseExpression(LOWEST))

	if p.errorExist() {
		return nil
	}

	p.nextToken() // consume token

	for p.curTokenIs(token.COMMA) {
		p.nextToken() // consume ','
		args = append(args, p.parseExpression(LOWEST))
		if p.errorExist() {
			return nil
		}
		p.nextToken() // consume token
	}

	if !p.curTokenIs(token.RPAREN) {
		p.createErrorMessage("expected closing parenthesis for function call")
		return nil
	}

	p.nextToken() // consume ')'
	p.stackTrace.Pop()
	return args
}
