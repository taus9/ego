package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/lexer"
	"ego/internal/frontend/parser"
	"testing"
)

func TestExecLiteralExpression(t *testing.T) {
	input := "`ls -la`"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	execLiteral, ok := stmt.Expression.(*ast.ExecLiteral)
	if !ok {
		t.Fatalf("exp not *ast.Exec. got=%T", stmt.Expression)
	}

	if execLiteral.Command != "ls -la" {
		t.Errorf("execLiteral.Command not %q. got=%q", "ls -la", execLiteral.Command)
	}
}
