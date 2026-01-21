package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
	"fmt"
)

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // consume current expression
		elseExpr := p.parseElseExpression(stmt.Expression)
		if p.errorExist() {
			return nil
		}

		stmt.Expression = elseExpr
	}

	if p.peekTokenIs(token.EOL) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	if p.errorExist() {
		return nil
	}

	for !p.peekTokenIs(token.EOL) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
		if p.errorExist() {
			return nil
		}
	}

	if p.errorExist() {
		return nil
	}
	return leftExp
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.curTokenIs(end) {
		p.nextToken() // consume end token
		return list
	}

	list = append(list, p.parseExpression(LOWEST))

	if p.errorExist() {
		return nil
	}

	p.nextToken() // move to next token

	for p.curTokenIs(token.COMMA) {
		p.nextToken() // consume ',' token

		list = append(list, p.parseExpression(LOWEST))

		if p.errorExist() {
			return nil
		}

		p.nextToken() // move to next token
	}

	if !p.curTokenIs(end) {
		p.createErrorMessage(fmt.Sprintf("expected next token to be %s, got %s instead", end, p.peekToken.Type))
		return nil
	}

	return list
}
