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

	EOL = "EOL" // - p
	EOF = "EOF" // - p

	// Identifiers + literals
	IDENT  = "IDENT"  // add, foobar, x, y, ... - p
	INT    = "INT"    // 1343 - p
	FLOAT  = "FLOAT"  // 3.14 - p
	STRING = "STRING" // "foobar" - p

	// Operators

	ASSIGN   = "=" // - p
	PLUS     = "+" // - p
	MINUS    = "-" // - p
	BANG     = "!" // - p
	ASTERISK = "*" // - p
	SLASH    = "/" // - p
	MOD      = "%" // - p

	LT     = "<"  // - p
	GT     = ">"  // - p
	EQ     = "==" // - p
	NOT_EQ = "!=" // - p
	LT_EQ  = "<=" // - p
	GT_EQ  = ">=" // - p

	COMMA        = "," // - p
	SINGLE_QUOTE = "'" // - p
	LBRACKET     = "[" // - p
	RBRACKET     = "]" // - p

	LPAREN = "(" // - p
	RPAREN = ")" // - p
	LBRACE = "{" // - p
	RBRACE = "}" // - p

	COLON         = ":"  // - p
	ANON_FUNCTION = ":>" // - p
	END_BLOCK     = ";"  // - p
	LET           = ":="

	TRUE   = "TRUE"  // - p
	FALSE  = "FALSE" // - p
	IF     = "IF"    // - p
	ELSE   = "ELSE"  // - p
	RETURN = "RET"   // - p
	FOR    = "FOR"   // - p
	TO     = "TO"    // - p
	AND    = "AND"   // - p
	OR     = "OR"    // - p
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
	"to":    TO,
	"and":   AND,
	"or":    OR,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func (t Token) String() string {
	return fmt.Sprintf("Type: %s, Literal: %s, Line: %d, Column: %d", t.Type, t.Literal, t.Span.Line, t.Span.Column)
}
