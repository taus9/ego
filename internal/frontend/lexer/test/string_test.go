package lexer

import (
	"ego/internal/frontend/lexer"
	"ego/internal/frontend/token"
	"testing"
)

func TestStringToken(t *testing.T) {
	input := "'testing'"

	l := lexer.New(input)

	tok := l.NextToken()
	if tok.Type != token.STRING {
		t.Fatalf("token type wrong. expected=%q, got=%q",
			token.STRING, tok.Type)
	}

	if tok.Literal != "testing" {
		t.Fatalf("literal wrong. expected=%q, got=%q",
			"testing", tok.Literal)
	}
}
