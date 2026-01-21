package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

func (p *Parser) parseElseExpression(tryExpr ast.Expression) ast.Expression {
	p.stackTrace.Push(ELSE)

	p.nextToken() // consume ELSE

	if p.curTokenIs(token.EOL) {
		return p.parseElseBlockExpression(tryExpr)
	}

	expression := p.parseExpression(LOWEST)

	if p.errorExist() {
		return nil
	}

	p.stackTrace.Pop()

	return &ast.ElseInlineExpression{
		Token:          token.Token{Type: token.ELSE, Literal: "else"},
		TryExpression:  tryExpr,
		ElseExpression: expression,
	}
}

func (p *Parser) parseElseBlockExpression(tryExpr ast.Expression) ast.Expression {
	expression := &ast.ElseBlockExpression{
		Token:         token.Token{Type: token.ELSE, Literal: "else"},
		TryExpression: tryExpr,
	}

	p.blockStack.Push(token.ELSE)

	expression.ElseBlock = p.parseBlockStatement()

	if p.errorExist() {
		return nil
	}

	p.blockStack.Pop()

	p.stackTrace.Pop()
	return expression
}
