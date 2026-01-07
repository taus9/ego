package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if p.peekTokenIs(token.EOL) || p.peekTokenIs(token.EOF) {
		p.createErrorMessage("expected expression after IF")
		return nil
	}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if p.errorsExist() {
		return nil
	}

	p.blockStack.Push(token.IF)

	p.nextToken()

	if !p.curTokenIs(token.EOL) && !p.errorsExist() {
		p.createErrorMessage("expected EOL after IF condition")
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.errorsExist() {
		return nil
	}

	p.blockStack.Pop()

	if p.curTokenIs(token.ELSE) {
		p.nextToken()

		if !p.curTokenIs(token.EOL) {
			p.createErrorMessage("expected EOL after ELSE")
			return nil
		}

		p.blockStack.Push(token.ELSE)

		expression.Alternative = p.parseBlockStatement()

		if p.errorsExist() {
			return nil
		}

		p.blockStack.Pop()
	}

	return expression
}
