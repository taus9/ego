package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
)

// I'm proud of this one
func (p *Parser) parseWhenExpression() ast.Expression {
	p.stackTrace.Push(WHEN)

	expression := &ast.WhenExpression{Token: p.curToken}

	p.nextToken() // consume when token

	if p.curTokenIs(token.EOL) {
		expression.Selector = nil
	} else {
		expression.Selector = p.parseExpression(LOWEST)
	}

	if p.errorExist() {
		return nil
	}

	p.nextToken()
	hasElse := false

	for !p.curTokenIs(token.EOF) && !p.curTokenIs(token.END_BLOCK) {
		switch p.curToken.Type {
		case token.IS:
			p.stackTrace.Push(IS)
			whenCase := &ast.WhenCase{Token: p.curToken}

			p.nextToken()
			whenCase.Condition = p.parseExpression(LOWEST)

			if p.errorExist() {
				return nil
			}

			// IS single line statement
			if p.peekTokenIs(token.COLON) {
				p.nextToken() // advance to COLON

				whenCase.Block = p.parseSingleLineWhenCase(whenCase.Token)

				if p.errorExist() {
					return nil
				}

				expression.Cases = append(expression.Cases, whenCase)

				p.nextToken()
				p.stackTrace.Pop()
				continue
			}

			// IS block statement
			if p.peekTokenIs(token.EOL) {
				p.nextToken()
				whenCase.Block = p.parseBlockStatement()

				if p.errorExist() {
					return nil
				}

				expression.Cases = append(expression.Cases, whenCase)

				p.nextToken()
				p.stackTrace.Pop()
				continue
			}

			p.createErrorMessage("only { or EOL allowed after IS condition")
			return nil

		case token.ELSE:
			p.stackTrace.Push(ELSE)
			if hasElse {
				p.createErrorMessage("ELSE block already defined for WHEN expression")
				return nil
			}

			elseToken := p.curToken

			if p.peekTokenIs(token.COLON) {
				p.nextToken() // advance to COLON

				expression.ElseBlock = p.parseSingleLineWhenCase(elseToken)

				if p.errorExist() {
					return nil
				}

				p.nextToken()
				hasElse = true
				p.stackTrace.Pop()
				continue
			}

			if p.peekTokenIs(token.EOL) {
				p.nextToken()
				expression.ElseBlock = p.parseBlockStatement()

				if p.errorExist() {
					return nil
				}

				p.nextToken()
				hasElse = true
				p.stackTrace.Pop()
				continue
			}

			p.createErrorMessage("only : or EOL allowed after IS condition")
			return nil

		case token.EOL:
			p.nextToken()

		default:
			p.createErrorMessage("expected IS or ELSE in WHEN expression")
			return nil
		}
	}

	if p.curTokenIs(token.EOF) {
		p.createErrorMessage("missing close ; for when expression")
		return nil
	}

	if !p.peekTokenIs(token.EOL) && !p.peekTokenIs(token.EOF) {
		p.createErrorMessage("expected EOL or EOF after END BLOCK")
		return nil
	}

	p.stackTrace.Pop()
	return expression
}

func (p *Parser) parseSingleLineWhenCase(tok token.Token) *ast.BlockStatement {
	p.stackTrace.Push(BLOCK)
	block := &ast.BlockStatement{
		Token: token.Token{Type: token.BEGIN_INTERNAL,
			Literal: token.BEGIN_INTERNAL,
		},
	}
	block.Statements = []ast.Statement{}

	p.nextToken() // consume :

	expression := p.parseExpression(LOWEST)

	if p.errorExist() {
		return nil
	}

	if !p.peekTokenIs(token.EOL) {
		p.createErrorMessage("Expected EOL after single line WHEN case")
		return nil
	}

	p.nextToken() // advance to EOL

	stmt := &ast.ExpressionStatement{Token: tok, Expression: expression}
	block.Statements = append(block.Statements, stmt)

	p.stackTrace.Pop()
	return block
}
