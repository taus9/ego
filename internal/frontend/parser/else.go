package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

func (p *Parser) parseElseExpression(tryExpr ast.Expression) ast.Expression {
	p.stackTrace.Push(ELSE)

	elseToken := p.curToken
	p.nextToken() // consume ELSE

	if p.curTokenIs(token.EOL) {
		return p.parseElseBlockExpression(tryExpr, elseToken)
	}

	expression := p.parseExpression(LOWEST)

	if p.errorExist() {
		return nil
	}

	p.stackTrace.Pop()

	return &ast.ElseInlineExpression{
		Token:          elseToken,
		TryExpression:  tryExpr,
		ElseExpression: expression,
	}
}

func (p *Parser) parseElseBlockExpression(tryExpr ast.Expression, elseToken token.Token) ast.Expression {
	expression := &ast.ElseBlockExpression{
		Token:         elseToken,
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
