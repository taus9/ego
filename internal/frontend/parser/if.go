package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

func (p *Parser) parseIfExpression() ast.Expression {
	p.stackTrace.Push(IF)
	expression := &ast.IfExpression{Token: p.curToken}

	if p.peekTokenIs(token.EOL) || p.peekTokenIs(token.EOF) {
		p.createErrorMessage("expected expression after IF")
		return nil
	}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if p.errorExist() {
		return nil
	}

	p.blockStack.Push(token.IF)

	p.nextToken()

	if !p.curTokenIs(token.EOL) && !p.errorExist() {
		p.createErrorMessage("expected EOL after IF condition")
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.errorExist() {
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

		if p.errorExist() {
			return nil
		}

		p.blockStack.Pop()
	}

	if !p.peekTokenIs(token.EOL) && !p.peekTokenIs(token.EOF) {
		p.createErrorMessage("token after END BLOCK must be EOL or EOF")
		return nil
	}

	p.stackTrace.Pop()
	return expression
}
