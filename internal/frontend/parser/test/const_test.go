package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/lexer"
	"ego/internal/frontend/parser"
	"testing"
)

func TestParseConstant(t *testing.T) {
	input := "$PI"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	constant, ok := stmt.Expression.(*ast.Constant)
	if !ok {
		t.Fatalf("exp not *ast.Constant. got=%T", stmt.Expression)
	}
	if constant.Value != "$PI" {
		t.Errorf("constant.Value not %s. got=%s", "$PI", constant.Value)
	}
	if constant.TokenLiteral() != "$PI" {
		t.Errorf("constant.TokenLiteral not %s. got=%s", "$PI",
			constant.TokenLiteral())
	}
}
