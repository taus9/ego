package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

func (p *Parser) parseArrayLiteral() ast.Expression {
	p.stackTrace.Push(ARRAY)
	array := &ast.ArrayLiteral{Token: p.curToken}

	array.Elements = p.parseExpressionList(token.RBRACKET)

	if p.errorExist() {
		return nil
	}

	p.stackTrace.Pop()
	return array
}
