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

func TestUnterminatedStringToken(t *testing.T) {
	input := "'this is unterminated"

	l := lexer.New(input)

	tok := l.NextToken()
	if tok.Type != token.UNTERMINATED_STRING {
		t.Fatalf("token type wrong. expected=%q, got=%q",
			token.UNTERMINATED_STRING, tok.Type)
	}
}

func TestExecToken(t *testing.T) {
	input := "`ls -la`"

	l := lexer.New(input)

	tok := l.NextToken()
	if tok.Type != token.EXEC {
		t.Fatalf("token type wrong. expected=%q, got=%q",
			token.EXEC, tok.Type)
	}

	if tok.Literal != "ls -la" {
		t.Fatalf("literal wrong. expected=%q, got=%q",
			"ls -la", tok.Literal)
	}
}

func TestUnterminatedExecToken(t *testing.T) {
	input := "`this is unterminated"

	l := lexer.New(input)

	tok := l.NextToken()
	if tok.Type != token.UNTERMINATED_STRING {
		t.Fatalf("token type wrong. expected=%q, got=%q",
			token.UNTERMINATED_STRING, tok.Type)
	}
}
