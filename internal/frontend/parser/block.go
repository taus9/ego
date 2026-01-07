package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
	"fmt"
)

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	p.stackTrace.Push(BLOCK)
	block := &ast.BlockStatement{
		Token: token.Token{Type: token.BEGIN_INTERNAL,
			Literal: token.BEGIN_INTERNAL,
		},
	}
	block.Statements = []ast.Statement{}

	if p.curTokenIs(token.END_BLOCK) {
		p.nextToken() // consume END BLOCK token
		if !p.curTokenIs(token.EOL) && !p.curTokenIs(token.EOF) {
			p.createErrorMessage("token after END BLOCK must be EOL or EOF")
			return nil
		}
		p.stackTrace.Pop()
		return block
	}

	p.nextToken()

	if p.curTokenIs(token.EOF) {
		p.createErrorMessage(fmt.Sprintf("missing END BLOCK for %s", p.blockStack.Peek()))
		return nil
	}

	for !p.curTokenIs(token.END_BLOCK) && !p.curTokenIs(token.ELSE) {

		stmt := p.parseStatement()

		if p.errorExist() {
			return nil
		}

		block.Statements = append(block.Statements, stmt)

		p.nextToken()

		if p.curTokenIs(token.ELSE) && p.blockStack.Peek() != token.IF {
			p.createErrorMessage("unexpected ELSE without matching IF")
			return nil
		}

		if p.curTokenIs(token.EOF) {
			p.createErrorMessage(fmt.Sprintf("missing END BLOCK for %s", p.blockStack.Peek()))
			return nil
		}
	}

	p.nextToken() // consume END BLOCK or ELSE token

	if !p.curTokenIs(token.EOL) && !p.curTokenIs(token.EOF) {
		p.createErrorMessage("token after END BLOCK must be EOL or EOF")
		return nil
	}

	p.stackTrace.Pop()
	return block
}

func (p *Parser) parseEndBlockStatement() ast.Expression {
	// this function currently just prevents nonsense parsing errors
	// END_BLOCK token already consumed in parseBlockStatement
	// so if another END_BLOCK is encountered here, it's unexpected
	p.createErrorMessage("unexpected END BLOCK without matching BEGIN BLOCK")
	return nil
}

func (p *Parser) parseElseExpression() ast.Expression {
	// this function currently just prevents nonsense parsing errors
	// ELSE token already consumed in parseBlockStatement
	// so if another ELSE is encountered here, it's unexpected
	p.createErrorMessage("unexpected ELSE without matching IF")
	return nil
}
