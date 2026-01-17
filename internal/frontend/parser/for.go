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

func (p *Parser) parseForStatement() ast.Statement {
	p.stackTrace.Push(FOR)
	p.nextToken() // consume 'for' token

	if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.DECLARE) {
		return p.parseForToStatement()
	}

	if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.COMMA) {
		return p.parseForInStatement()
	}

	return p.parseForWhileStatement()
}

func (p *Parser) parseForWhileStatement() *ast.ForWhileStatement {
	fws := &ast.ForWhileStatement{Token: token.Token{Type: token.FOR, Literal: "for"}}

	condition := p.parseExpression(LOWEST)

	if p.errorExist() {
		return nil
	}

	fws.Condition = condition
	p.nextToken()

	if !p.curTokenIs(token.EOL) {
		p.createErrorMessage("expected EOL after for while condition expression")
		return nil
	}

	p.blockStack.Push(token.FOR)
	fws.Body = p.parseBlockStatement()

	if p.errorExist() {
		return nil
	}

	p.blockStack.Pop()

	//END BLOCK token already consumed in parseBlockStatement
	p.stackTrace.Pop()
	return fws
}

func (p *Parser) parseForToStatement() *ast.ForToStatement {
	fts := &ast.ForToStatement{Token: token.Token{Type: token.FOR, Literal: "for"}}
	fts.Iterator = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // consume identifier

	// redundant check
	if !p.curTokenIs(token.DECLARE) {
		p.createErrorMessage("expected ':=' after for index identifier")
		return nil
	}

	p.nextToken() // consume ':=' token
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

func (p *Parser) parseForInStatement() *ast.ForInStatement {
	fis := &ast.ForInStatement{Token: token.Token{Type: token.FOR, Literal: "for"}}
	fis.Index = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // consume index identifier
	p.nextToken() // consume ',' token

	if !p.curTokenIs(token.IDENT) {
		p.createErrorMessage("expected identifier after ',' in for in statement")
		return nil
	}

	fis.Value = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // consume value identifier

	if !p.curTokenIs(token.IN) {
		p.createErrorMessage("expected 'in' after for in value identifier")
		return nil
	}

	p.nextToken() // consume 'in' token

	fis.Iterable = p.parseExpression(LOWEST)

	if p.errorExist() {
		return nil
	}

	p.nextToken()

	if !p.curTokenIs(token.EOL) {
		p.createErrorMessage("expected EOL after for in iterable expression")
		return nil
	}

	p.blockStack.Push(token.FOR)
	fis.Body = p.parseBlockStatement()

	if p.errorExist() {
		return nil
	}

	p.blockStack.Pop()

	//END BLOCK token already consumed in parseBlockStatement
	p.stackTrace.Pop()
	return fis
}
