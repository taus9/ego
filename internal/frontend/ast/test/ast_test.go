package ast

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.DeclareStatement{
				Token: token.Token{Type: token.DECLARE_INTERNAL, Literal: "$LET"},
				Name: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "$LET myVar := anotherVar" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
