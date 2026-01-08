package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

func (p *Parser) parseArrayLiteral() ast.Expression {
	p.stackTrace.Push(ARRAY)
	array := &ast.ArrayLiteral{Token: p.curToken}

	p.nextToken() // consume '[' token

	array.Elements = p.parseExpressionList(token.RBRACKET)

	if p.errorExist() {
		return nil
	}

	// ']' token already consumed in parseExpressionList
	p.stackTrace.Pop()
	return array
}
