package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

/*

for <ident> = <exp> to <exp>
	<block>
;

*/

func (p *Parser) parseForToStatement() *ast.ForToStatement {
	p.stackTrace.Push(FOR)
	fts := &ast.ForToStatement{Token: p.curToken}

	p.nextToken() // consume 'for' token

	if !p.curTokenIs(token.IDENT) {
		p.createErrorMessage("expected valid identifier after 'for'")
		return nil
	}

	// Index identifier
	fts.Iterator = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // consume identifier

	if !p.curTokenIs(token.ASSIGN) {
		p.createErrorMessage("expected '=' after for index identifier")
		return nil
	}

	p.nextToken() // consume '=' token

	fts.Start = p.parseExpression(LOWEST)

	if p.errorExist() {
		return nil
	}

	p.nextToken()

	if !p.curTokenIs(token.TO) {
		p.createErrorMessage("expected 'to' after for start expression")
		return nil
	}

	p.nextToken() // consume 'to' token

	fts.End = p.parseExpression(LOWEST)

	if p.errorExist() {
		return nil
	}

	p.nextToken()

	if !p.curTokenIs(token.EOL) {
		p.createErrorMessage("expected EOL after for to expression")
		return nil
	}

	p.blockStack.Push(token.FOR)
	fts.Body = p.parseBlockStatement()

	if p.errorExist() {
		return nil
	}

	p.blockStack.Pop()

	//END BLOCK token already consumed in parseBlockStatement
	p.stackTrace.Pop()
	return fts
}
