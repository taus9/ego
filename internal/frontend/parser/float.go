package parser

import (
	"ego/internal/frontend/ast"
	"fmt"
	"strconv"
)

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.createErrorMessage(fmt.Sprintf("could not parse %s as float", p.curToken.Literal))
		return nil
	}
	lit.Value = value
	return lit
}
