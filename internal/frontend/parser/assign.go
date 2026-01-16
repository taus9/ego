package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

func (p *Parser) parseAssignStatement(exp ast.Expression) *ast.AssignStatement {
	p.stackTrace.Push(ASSIGN)
	assignStmt := &ast.AssignStatement{Token: token.Token{Type: token.ASSIGN, Literal: token.ASSIGN_INTERNAL}}

	switch node := exp.(type) {
	case *ast.Identifier:
		assignStmt.Name = node
		assignStmt.Index = nil
	case *ast.IndexExpression:
		assignStmt.Name = nil
		assignStmt.Index = node
	default:
		p.createErrorMessage("invalid assignment target")
		return nil
	}

	p.nextToken() // consume '='
	assignStmt.Value = p.parseExpression(LOWEST)
	if p.errorExist() {
		return nil
	}

	if p.peekTokenIs(token.EOL) {
		p.nextToken()
	}

	p.stackTrace.Pop()
	return assignStmt
}
