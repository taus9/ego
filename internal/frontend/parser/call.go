package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	p.stackTrace.Push(CALL)
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	if p.errorExist() {
		return nil
	}
	p.stackTrace.Pop()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	p.stackTrace.Push(CALL_ARGS)
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		p.stackTrace.Pop()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	if p.errorExist() {
		return nil
	}

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
		if p.errorExist() {
			return nil
		}
	}

	if !p.expectPeek(token.RPAREN) {
		p.createErrorMessage("expected closing parenthesis for function call")
		return nil
	}

	p.stackTrace.Pop()
	return args
}
