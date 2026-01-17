package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/lexer"
	"ego/internal/frontend/parser"
	"testing"
)

func TestParsingBreakStatement(t *testing.T) {
	input := "for true \n break \n ;"
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ForWhileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ForWhileStatement. got=%T",
			program.Statements[0])
	}

	breakStmt, ok := stmt.Body.Statements[0].(*ast.BreakStatement)
	if !ok {
		t.Fatalf("stmt.Body.Statements[0] is not ast.BreakStatement. got=%T",
			stmt.Body.Statements[0])
	}

	if breakStmt.String() != "break" {
		t.Errorf("stmt.TokenLiteral not 'break', got=%q",
			breakStmt.String())
	}

}

func TestParsingAgainStatement(t *testing.T) {
	// terrible example lol
	input := "for true \n again \n ;"
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ForWhileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ForWhileStatement. got=%T",
			program.Statements[0])
	}

	againStmt, ok := stmt.Body.Statements[0].(*ast.AgainStatement)
	if !ok {
		t.Fatalf("stmt.Body.Statements[0] is not ast.AgainStatement. got=%T",
			stmt.Body.Statements[0])
	}

	if againStmt.String() != "again" {
		t.Errorf("stmt.TokenLiteral not 'again', got=%q",
			againStmt.String())
	}
}
