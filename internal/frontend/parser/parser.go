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
	DECLARE
	RETURN
	MAP
	CALL_ARGS
	FUNC
	FOR
	ASSIGN
)

const (
	_ int = iota
	LOWEST
	OR
	AND
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

var precedences = map[token.TokenType]int{
	token.OR:       OR,
	token.AND:      AND,
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LT_EQ:    LESSGREATER,
	token.GT_EQ:    LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.MOD:      PRODUCT,
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

	currentExpression *ast.Expression

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
	p.registerPrefix(token.CONST, p.parseConstant)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.ANON_FUNCTION, p.parseAnonymousFunction)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseMapLiteral)

	// unexpected tokens
	p.registerPrefix(token.EOL, p.parseUnexpectedTokenError)
	p.registerPrefix(token.EOF, p.parseUnexpectedTokenError)
	p.registerPrefix(token.ASSIGN, p.parseUnexpectedTokenError)
	p.registerPrefix(token.DECLARE, p.parseUnexpectedTokenError)
	p.registerPrefix(token.PLUS, p.parseUnexpectedTokenError)
	p.registerPrefix(token.ASTERISK, p.parseUnexpectedTokenError)
	p.registerPrefix(token.SLASH, p.parseUnexpectedTokenError)
	p.registerPrefix(token.MOD, p.parseUnexpectedTokenError)
	p.registerPrefix(token.LT, p.parseUnexpectedTokenError)
	p.registerPrefix(token.GT, p.parseUnexpectedTokenError)
	p.registerPrefix(token.EQ, p.parseUnexpectedTokenError)
	p.registerPrefix(token.NOT_EQ, p.parseUnexpectedTokenError)
	p.registerPrefix(token.LT_EQ, p.parseUnexpectedTokenError)
	p.registerPrefix(token.GT_EQ, p.parseUnexpectedTokenError)
	p.registerPrefix(token.COMMA, p.parseUnexpectedTokenError)
	p.registerPrefix(token.SINGLE_QUOTE, p.parseUnexpectedTokenError)
	p.registerPrefix(token.RBRACKET, p.parseUnexpectedTokenError)
	p.registerPrefix(token.RPAREN, p.parseUnexpectedTokenError)
	p.registerPrefix(token.RBRACE, p.parseUnexpectedTokenError)
	p.registerPrefix(token.COLON, p.parseUnexpectedTokenError)
	p.registerPrefix(token.FOR, p.parseUnexpectedTokenError)
	p.registerPrefix(token.RETURN, p.parseUnexpectedTokenError)
	p.registerPrefix(token.TO, p.parseUnexpectedTokenError)
	p.registerPrefix(token.AND, p.parseUnexpectedTokenError)
	p.registerPrefix(token.OR, p.parseUnexpectedTokenError)

	p.registerPrefix(token.ILLEGAL, p.parseIllegalTokenError)

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
	p.registerInfix(token.MOD, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.LT_EQ, p.parseInfixExpression)
	p.registerInfix(token.GT_EQ, p.parseInfixExpression)

	// Read two tokens, so curToken and peekToken are both set
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
		if p.peekTokenIs(token.DECLARE) {
			return p.parseDeclareStatement()
		}

		exp := p.parseExpressionStatement()
		if p.peekTokenIs(token.ASSIGN) && p.isAssignableExpression(exp) {
			p.nextToken() // jump to '=' token
			return p.parseAssignStatement(exp.Expression)
		}

		return exp
	case token.COLON:
		return p.parseFunctionStatement()

	case token.RETURN:
		return p.parseReturnStatement()

	case token.FOR:
		return p.parseForStatement()

	case token.BREAK:
		if !p.stackTrace.Has(FOR) {
			p.createErrorMessage("break statement not within a for loop")
			return nil
		}
		return p.parseBreakStatement()

	case token.AGAIN:
		if !p.stackTrace.Has(FOR) {
			p.createErrorMessage("again statement not within a for loop")
			return nil
		}
		return p.parseAgainStatement()

	case token.EOL:
		return nil

	case token.UNTERMINATED_STRING:
		p.createErrorMessage("unterminated string")
		return nil

	case token.UNNAMED_CONST:
		p.createErrorMessage("constant '$' must be followed by an identifier")
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

func (p *Parser) parseIllegalTokenError() ast.Expression {
	p.createErrorMessage(fmt.Sprintf("illegal token: %s", p.curToken.Literal))
	return nil
}

func (p *Parser) errorExist() bool {
	return p.parseError != nil
}

func (p *Parser) isAssignableExpression(exp *ast.ExpressionStatement) bool {
	// currently only support for one dimensional index expressions
	switch node := exp.Expression.(type) {
	case *ast.Identifier:
		return true
	case *ast.IndexExpression:
		if _, ok := node.Left.(*ast.Identifier); !ok {
			return false
		}
		return true
	default:
		return false
	}
}
