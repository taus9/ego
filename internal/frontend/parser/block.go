package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
	"fmt"
)

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: token.Token{Type: token.BEGIN_INTERNAL,
			Literal: token.BEGIN_INTERNAL,
		},
	}
	block.Statements = []ast.Statement{}

	p.nextToken()

	if p.curTokenIs(token.EOF) {
		p.createErrorMessage(fmt.Sprintf("missing END BLOCK for %s", p.blockStack.Peek()))
		return nil
	}

	for !p.curTokenIs(token.END_BLOCK) && !p.curTokenIs(token.ELSE) {

		stmt := p.parseStatement()

		if p.errorsExist() {
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

	return block
}

// func (p *Parser) parseEndBlockStatement() ast.Expression {
// 	// this function currently just prevents nonsense parsing errors
// 	// when an END_BLOCK is encountered unexpectedly
// 	p.nextToken()
// 	if !p.curTokenIs(token.EOL) && !p.curTokenIs(token.EOF) {
// 		p.createErrorMessage("token after END BLOCK must be EOL or EOF")
// 	}

// 	if len(p.blockStack.elements) == 0 {
// 		p.createErrorMessage("unexpected END BLOCK without matching BEGIN BLOCK")
// 	} else {
// 		p.blockStack.Pop()
// 	}
// 	return nil
// }

// func (p *Parser) parseElseExpression() ast.Expression {
// 	// this function currently just prevents nonsense parsing errors
// 	// when an ELSE is encountered unexpectedly
// 	p.nextToken()

// 	if !p.curTokenIs(token.EOL) {
// 		p.createErrorMessage("token after ELSE must be EOL")
// 	}

// 	if len(p.blockStack.elements) == 0 || p.blockStack.Peek() != token.IF {
// 		p.createErrorMessage("unexpected ELSE without matching IF")
// 	} else {
// 		p.blockStack.Pop()
// 	}
// 	return nil
// }
