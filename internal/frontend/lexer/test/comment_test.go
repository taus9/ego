package lexer

import (
	"ego/internal/frontend/lexer"
	"ego/internal/frontend/token"
	"testing"
)

func TestSkipComment(t *testing.T) {
	input := "a = 5 # This is a comment"

	lexer := lexer.New(input)

	tok := lexer.NextToken()
	if tok.Type != token.IDENT || tok.Literal != "a" {
		t.Errorf("Expected IDENT token with literal 'a', got %v with literal %q", tok.Type, tok.Literal)
	}

	tok = lexer.NextToken()
	if tok.Type != token.ASSIGN || tok.Literal != "=" {
		t.Errorf("Expected ASSIGN token with literal '=', got %v with literal %q", tok.Type, tok.Literal)
	}

	tok = lexer.NextToken()
	if tok.Type != token.INT || tok.Literal != "5" {
		t.Errorf("Expected INT token with literal '5', got %v with literal %q", tok.Type, tok.Literal)
	}

	tok = lexer.NextToken()
	if tok.Type != token.EOF {
		t.Errorf("Expected EOF token after skipping comment, got %v", tok.Type)
	}
}
