package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/lexer"
	"ego/internal/frontend/token"
	"fmt"
)

const (
	_ int = iota
	PROGRAM
	ANON_FUNCTION
	ARRAY
	BLOCK
	IF
	LET
	RETURN
	MAP
	CALL_ARGS
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	parseError *ParseError

	curToken  token.Token
	peekToken token.Token

	blockStack *Stack
	stackTrace *Stack

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:          l,
		parseError: nil,
		blockStack: NewStack(),
		stackTrace: NewStack(),
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.ANON_FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseMapLiteral)

	// unexpected tokens
	p.registerPrefix(token.EOL, p.parseUnexpectedTokenError)
	p.registerPrefix(token.EOF, p.parseUnexpectedTokenError)

	p.registerPrefix(token.END_BLOCK, p.parseEndBlockStatement)
	p.registerPrefix(token.ELSE, p.parseElseExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) ParseError() *ParseError {
	return p.parseError
}

func (p *Parser) peekError(t token.TokenType) {
	p.createErrorMessage(fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type))
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	p.stackTrace.Push(PROGRAM)

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()

		// Stop parsing on the first error encountered
		if p.errorExist() {
			return program
		}

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.IDENT:
		if p.peekTokenIs(token.ASSIGN) {
			return p.parseLetStatement()
		}
		return p.parseExpressionStatement()

	case token.RETURN:
		return p.parseReturnStatement()

	case token.EOL:
		return nil

	case token.UNTERMINATED_STRING:
		p.createErrorMessage("unterminated string")
		return nil

	case token.ILLEGAL:
		p.createErrorMessage(fmt.Sprintf("illegal token: %s", p.curToken.Literal))
		return nil

	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) createErrorMessage(messsage string) {
	p.parseError = NewParseError(messsage, p.stackTrace, p.curToken)
}

func (p *Parser) parseUnexpectedTokenError() ast.Expression {
	p.createErrorMessage(fmt.Sprintf("unexpected token: %s", p.curToken.Type))
	return nil
}

func (p *Parser) errorExist() bool {
	return p.parseError != nil
}
