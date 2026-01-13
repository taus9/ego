package token

import "fmt"

type Span struct {
	Line   int
	Column int
}

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Span    Span
}

const (
	ILLEGAL             = "ILLEGAL"
	UNTERMINATED_STRING = "UNTERMINATED_STRING"

	EOL = "EOL"
	EOF = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // 1343
	STRING = "STRING" // "foobar"

	// Operators

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT     = "<"
	GT     = ">"
	EQ     = "=="
	NOT_EQ = "!="

	COMMA        = ","
	SINGLE_QUOTE = "'"
	LBRACKET     = "["
	RBRACKET     = "]"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	COLON         = ":"
	ANON_FUNCTION = ":>"
	END_BLOCK     = ";"
	LET           = ":="

	TRUE   = "TRUE"
	FALSE  = "FALSE"
	IF     = "IF"
	ELSE   = "ELSE"
	RETURN = "RET"
	FOR    = "FOR"

	// used internally by the parsers
	LET_INTERNAL   = "$LET"
	BEGIN_INTERNAL = "$BEGIN"
)

var keywords = map[string]TokenType{
	"true":  TRUE,
	"false": FALSE,
	"if":    IF,
	"else":  ELSE,
	"ret":   RETURN,
	"for":   FOR,
}

func LoopupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func (t Token) String() string {
	return fmt.Sprintf("Type: %s, Literal: %s, Line: %d, Column: %d", t.Type, t.Literal, t.Span.Line, t.Span.Column)
}
