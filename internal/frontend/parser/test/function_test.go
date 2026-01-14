package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/lexer"
	"ego/internal/frontend/parser"
	"testing"
)

func TestFunctionStatementNoArgs(t *testing.T) {
	input := ":add \n ;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.FunctionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FunctionStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Name.Value != "add" {
		t.Errorf("function name is not 'add'. got=%s", stmt.Name.Value)
	}

	if len(stmt.Parameters) != 0 {
		t.Errorf("function parameters length is not 0. got=%d",
			len(stmt.Parameters))
	}
}

func TestFunctionStatmentWithOneArg(t *testing.T) {
	input := ":add x \n ;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.FunctionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FunctionStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Name.Value != "add" {
		t.Errorf("function name is not 'add'. got=%s", stmt.Name.Value)
	}

	if len(stmt.Parameters) != 1 {
		t.Errorf("function parameters length is not 1. got=%d",
			len(stmt.Parameters))
	}

	if stmt.Parameters[0].Value != "x" {
		t.Errorf("function parameter is not 'x'. got=%s",
			stmt.Parameters[0].Value)
	}
}

func TestFunctionStatementWithMultipleArgs(t *testing.T) {
	input := ":add x, y, z \n ;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.FunctionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FunctionStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Name.Value != "add" {
		t.Errorf("function name is not 'add'. got=%s", stmt.Name.Value)
	}

	if len(stmt.Parameters) != 3 {
		t.Errorf("function parameters length is not 3. got=%d",
			len(stmt.Parameters))
	}

	expectedParams := []string{"x", "y", "z"}
	for i, param := range expectedParams {
		if stmt.Parameters[i].Value != param {
			t.Errorf("function parameter %d is not '%s'. got=%s",
				i, param, stmt.Parameters[i].Value)
		}
	}
}
