package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
	"fmt"
)

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	p.createErrorMessage(fmt.Sprintf("no prefix parse function for %s found", t))
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	if p.errorExist() {
		return nil
	}

	return expression
}
