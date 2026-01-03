package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOL     = "EOL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 1343

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

	FUNCTION      = ":"
	ANON_FUNCTION = ":>"
	END_BLOCK     = ";"
	LET           = ":="

	TRUE   = "TRUE"
	FALSE  = "FALSE"
	IF     = "IF"
	RETURN = "RET"

	// used internally by the parsers
	LET_INTERNAL   = "$LET"
	BEGIN_INTERNAL = "$BEGIN"
)

var keywords = map[string]TokenType{
	"true":  TRUE,
	"false": FALSE,
	"if":    IF,
	"ret":   RETURN,
}

func LoopupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
