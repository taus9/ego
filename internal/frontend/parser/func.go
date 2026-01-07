package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
	p.stackTrace.Push(FUNC)
	fs := &ast.FunctionStatement{Token: p.curToken}

	p.nextToken() // consume ':' token

	if !p.curTokenIs(token.IDENT) {
		p.createErrorMessage("expected valid function name")
		return nil
	}

	// Function name
	fs.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // consume function name

	fs.Parameters = p.parseFunctionParameters()

	if p.errorExist() {
		return nil
	}

	p.blockStack.Push(token.COLON)
	fs.Body = p.parseBlockStatement()
	p.blockStack.Pop()

	if p.errorExist() {
		return nil
	}

	//END BLOCK token already consumed in parseBlockStatement
	p.stackTrace.Pop()
	return fs
}
