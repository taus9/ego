package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/lexer"
	"ego/internal/frontend/parser"
	"testing"
)

func TestDeclareStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"x := 5", "x", 5},
		{"y := true", "y", true},
		{"foobar := y", "foobar", "y"},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testDeclareStatements(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.DeclareStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func testDeclareStatements(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "$DECLARE" {
		t.Errorf("s.TokenLiteral not '$DECLARE'. got=%q", s.TokenLiteral())
		return false
	}

	declareStmt, ok := s.(*ast.DeclareStatement)
	if !ok {
		t.Errorf("s not *ast.DeclareStatement. got=%T", s)
		return false
	}

	if declareStmt.Name.Value != name {
		t.Errorf("declareStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, declareStmt.Name.TokenLiteral())
		return false
	}

	return true
}
