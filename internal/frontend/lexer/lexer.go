package lexer

import (
	"ego/internal/frontend/token"
	"errors"
	"strings"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte

	line   int
	column int
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {

	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
	l.column += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {

	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
			tok.Span = token.Span{Line: l.line, Column: l.column - len(tok.Literal)}
		} else {
			tok = l.newToken(token.ASSIGN, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
			tok.Span = token.Span{Line: l.line, Column: l.column - len(tok.Literal)}
		} else {
			tok = l.newToken(token.BANG, l.ch)
		}
	case ':':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.ANON_FUNCTION, Literal: literal}
			tok.Span = token.Span{Line: l.line, Column: l.column - len(tok.Literal)}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DECLARE, Literal: literal}
			tok.Span = token.Span{Line: l.line, Column: l.column - len(tok.Literal)}
		} else {
			tok = l.newToken(token.COLON, l.ch)
		}
	case '\'':
		str, err := l.readString()
		if err != nil {
			tok.Type = token.UNTERMINATED_STRING
		} else {
			tok.Type = token.STRING
		}
		tok.Literal = str
		// column -1 because when we call readString we advance past the opening quote
		tok.Span = token.Span{Line: l.line, Column: l.column - len(tok.Literal) - 1}
	case ';':
		tok = l.newToken(token.END_BLOCK, l.ch)
	case '(':
		tok = l.newToken(token.LPAREN, l.ch)
	case ')':
		tok = l.newToken(token.RPAREN, l.ch)
	case ',':
		tok = l.newToken(token.COMMA, l.ch)
	case '+':
		tok = l.newToken(token.PLUS, l.ch)
	case '{':
		tok = l.newToken(token.LBRACE, l.ch)
	case '}':
		tok = l.newToken(token.RBRACE, l.ch)
	case '-':
		tok = l.newToken(token.MINUS, l.ch)
	case '/':
		tok = l.newToken(token.SLASH, l.ch)
	case '%':
		tok = l.newToken(token.MOD, l.ch)
	case '*':
		tok = l.newToken(token.ASTERISK, l.ch)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LT_EQ, Literal: literal}
			tok.Span = token.Span{Line: l.line, Column: l.column - len(tok.Literal)}
		} else {
			tok = l.newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GT_EQ, Literal: literal}
			tok.Span = token.Span{Line: l.line, Column: l.column - len(tok.Literal)}
		} else {
			tok = l.newToken(token.GT, l.ch)
		}
	case '[':
		tok = l.newToken(token.LBRACKET, l.ch)
	case ']':
		tok = l.newToken(token.RBRACKET, l.ch)
	case '\n':
		tok = l.newToken(token.EOL, l.ch)
		l.line++
		l.column = 0
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Span = token.Span{Line: l.line, Column: l.column}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Span = token.Span{Line: l.line, Column: l.column - len(tok.Literal)}
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			if strings.Contains(tok.Literal, ".") {
				tok.Type = token.FLOAT
			} else {
				tok.Type = token.INT
			}

			tok.Span = token.Span{Line: l.line, Column: l.column - len(tok.Literal)}
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	position := l.position
	hasDot := false
	for isDigit(l.ch) || (l.ch == '.' && !hasDot) {
		if l.ch == '.' {
			hasDot = true
		}
		l.readChar()
	}

	number := l.input[position:l.position]
	// If the number ends with a dot, backtrack one position
	if hasDot && number[len(number)-1] == '.' {
		l.position--
		l.readPosition--
		l.column--
		number = l.input[position:l.position]
	}

	return number
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
		Span:    token.Span{Line: l.line, Column: l.column},
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readString() (string, error) {
	position := l.position + 1
	err := error(nil)
	for {
		l.readChar()
		if l.ch == '\'' {
			break
		}
		if l.ch == '\n' || l.ch == 0 {
			err = errors.New(token.UNTERMINATED_STRING)
			break
		}
	}

	return l.input[position:l.position], err
}
