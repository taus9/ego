package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

func (p *Parser) parseIfExpression() ast.Expression {
	p.stackTrace.Push(IF)

	if p.peekTokenIs(token.EOL) || p.peekTokenIs(token.EOF) {
		p.createErrorMessage("expected expression after IF")
		return nil
	}

	p.nextToken() // consume IF

	condition := p.parseExpression(LOWEST)

	if p.errorExist() {
		return nil
	}

	if p.peekTokenIs(token.EOL) {
		p.nextToken()
		return p.parseIfBlockExpression(condition)
	}

	if p.peekTokenIs(token.LBRACE) {
		p.nextToken()
		return p.parseIfTernaryExpression(condition)
	}

	p.createErrorMessage("invalid if expression")
	return nil
}

func (p *Parser) parseIfBlockExpression(condition ast.Expression) ast.Expression {
	expression := &ast.IfBlockExpression{
		Token:     token.Token{Type: token.IF, Literal: "if"},
		Condition: condition,
	}

	p.blockStack.Push(token.IF)

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

func (p *Parser) parseIfTernaryExpression(condition ast.Expression) ast.Expression {
	expression := &ast.IfTernaryExpression{
		Token:     token.Token{Type: token.IF, Literal: "if"},
		Condition: condition,
	}

	p.nextToken() // consume '{'

	expression.Consequence = p.parseExpression(LOWEST)

	if p.errorExist() {
		return nil
	}

	if !p.peekTokenIs(token.RBRACE) {
		p.createErrorMessage("expected '}' after consequence expression")
		return nil
	}

	p.nextToken()

	if p.peekTokenIs(token.EOL) || p.peekTokenIs(token.EOF) {
		return expression
	}

	if !p.peekTokenIs(token.ELSE) {
		p.createErrorMessage("expected ELSE after consequence expression")
		return nil
	}

	p.nextToken()

	if !p.peekTokenIs(token.LBRACE) {
		p.createErrorMessage("expected '{' after ELSE")
		return nil
	}

	p.nextToken() // consume ELSE
	p.nextToken() // consume '{'

	expression.Alternative = p.parseExpression(LOWEST)

	if p.errorExist() {
		return nil
	}

	if !p.peekTokenIs(token.RBRACE) {
		p.createErrorMessage("expected '}' after alternative expression")
		return nil
	}

	p.nextToken() // consume '}'
	p.stackTrace.Pop()

	return expression
}
