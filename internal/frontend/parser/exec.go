package parser

import "ego/internal/frontend/ast"

func (p *Parser) parseExecLiteral() ast.Expression {
	lit := &ast.ExecLiteral{Token: p.curToken}

	lit.Command = p.curToken.Literal

	return lit
}
