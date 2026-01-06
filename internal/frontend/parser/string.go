package parser

import "ego/internal/frontend/ast"

func (p *Parser) parseStringLiteral() ast.Expression {
	lit := &ast.StringLiteral{Token: p.curToken}

	lit.Value = p.curToken.Literal

	return lit
}
