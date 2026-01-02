package parser

import (
	"ego/ast"
	"ego/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.IDENT:
		if p.peekTokenIs(token.ASSIGN) {
			return p.parseAssignStatement()
		}
		return nil
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}
