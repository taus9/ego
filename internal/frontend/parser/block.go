package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: token.Token{Type: token.BEGIN_INTERNAL,
			Literal: token.BEGIN_INTERNAL,
		},
	}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.END_BLOCK) && !p.curTokenIs(token.ELSE) {

		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.nextToken()

		if p.curTokenIs(token.ELSE) && p.blockStack.Peek() != token.IF {
			msg := "unexpected ELSE without matching IF"
			p.errors = append(p.errors, msg)
			return nil
		}

		if p.curTokenIs(token.EOF) {
			msg := "unexpected end of file, missing end block"
			p.errors = append(p.errors, msg)
			return nil
		}
	}

	return block
}
