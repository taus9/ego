package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/lexer"
	"ego/internal/frontend/parser"
	"testing"
)

func TestParseUseStatement(t *testing.T) {
	input := "use 'math' as m\n"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.UseStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.UseStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Module.Value != "math" {
		t.Errorf("stmt.Module.Value is not 'math'. got=%q", stmt.Module.Value)
	}

	if stmt.Alias == nil {
		t.Errorf("stmt.Alias is nil, expected Identifier with value 'm'")
	} else if stmt.Alias.Value != "m" {
		t.Errorf("stmt.Alias.Value is not 'm'. got=%q", stmt.Alias.Value)
	}
}
