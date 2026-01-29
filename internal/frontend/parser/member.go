package parser

import "ego/internal/frontend/ast"

func (p *Parser) parseMemberExpression(left ast.Expression) ast.Expression {
	expr := &ast.MemberExpression{Token: p.curToken}

	obj, ok := left.(*ast.Identifier)
	if !ok {
		p.createErrorMessage("invalid member expression object")
		return nil
	}
	expr.ModuleName = obj

	p.nextToken()

	property := p.parseExpression(LOWEST)
	if p.errorExist() {
		return nil
	}
	expr.ModuleProperty = property

	return expr
}
