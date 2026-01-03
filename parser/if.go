package parser

import (
	"ego/ast"
	"ego/token"
)

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if p.peekTokenIs(token.EOL) || p.peekTokenIs(token.EOF) {
		return nil
	}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.EOL) {
		return nil
	}

	p.blockStack.Push(token.IF)
	expression.Consequence = p.parseBlockStatement()
	p.blockStack.Pop()

	if p.curTokenIs(token.ELSE) {
		p.nextToken()

		if !p.curTokenIs(token.EOL) {
			return nil
		}

		p.blockStack.Push(token.ELSE)
		expression.Alternative = p.parseBlockStatement()
		p.blockStack.Pop()
	}

	return expression
}
