package parser

import (
	"ego/ast"
	"ego/token"
)

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if p.peekTokenIs(token.EOL) || p.peekTokenIs(token.EOF) {
		msg := "expected expression after IF"
		p.errors = append(p.errors, msg)
		return nil
	}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.EOL) {
		msg := "expected end of line after IF condition"
		p.errors = append(p.errors, msg)
		return nil
	}

	p.blockStack.Push(token.IF)
	expression.Consequence = p.parseBlockStatement()
	p.blockStack.Pop()

	if p.curTokenIs(token.ELSE) {
		p.nextToken()

		if !p.curTokenIs(token.EOL) {
			msg := "expected end of line after ELSE"
			p.errors = append(p.errors, msg)
			return nil
		}

		p.blockStack.Push(token.ELSE)
		expression.Alternative = p.parseBlockStatement()
		p.blockStack.Pop()
	}

	return expression
}
