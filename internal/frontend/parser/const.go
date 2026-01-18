package parser

import (
	"ego/internal/frontend/ast"
)

func (p *Parser) parseConstant() ast.Expression {
	return &ast.Constant{Token: p.curToken, Value: p.curToken.Literal}
}
