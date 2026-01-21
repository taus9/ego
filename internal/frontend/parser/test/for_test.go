package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/lexer"
	"ego/internal/frontend/parser"
	"testing"
)

func TestParsingForWhileStatement(t *testing.T) {
	input := "for i < 5 \n i := i + 1 \n ;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ForWhileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ForWhileStatement. got=%T",
			program.Statements[0])
	}

	if !testInfixExpression(t, stmt.Condition, "i", "<", 5) {
		return
	}

	if len(stmt.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(stmt.Body.Statements))
	}

	bodyStmt, ok := stmt.Body.Statements[0].(*ast.DeclareStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.LetStatement. got=%T",
			stmt.Body.Statements[0])
	}

	if bodyStmt.Name.Value != "i" {
		t.Errorf("let name is not 'i'. got=%s", bodyStmt.Name.Value)
	}

	if !testInfixExpression(t, bodyStmt.Value, "i", "+", 1) {
		return
	}
}

func TestParsingForToStatement(t *testing.T) {
	input := "for i := 0 + 2 to a \n ret i \n ;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ForToStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ForToStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Iterator.Value != "i" {
		t.Errorf("iterator is not 'i'. got=%s", stmt.Iterator.Value)
	}

	if !testInfixExpression(t, stmt.Start, 0, "+", 2) {
		return
	}

	if !testIdentifier(t, stmt.End, "a") {
		return
	}

	if len(stmt.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(stmt.Body.Statements))
	}

	bodyStmt, ok := stmt.Body.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ReturnStatement. got=%T",
			stmt.Body.Statements[0])
	}

	if !testIdentifier(t, bodyStmt.ReturnValue, "i") {
		return
	}
}

func TestParsingForInStatement(t *testing.T) {
	input := "for i, v in myArray \n ret v \n ;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ForInStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ForInStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Index.Value != "i" {
		t.Errorf("index is not 'i'. got=%s", stmt.Index.Value)
	}

	if stmt.Value.Value != "v" {
		t.Errorf("value is not 'v'. got=%s", stmt.Value.Value)
	}

	iterable, ok := stmt.Iterable.(*ast.Identifier)
	if !ok {
		t.Fatalf("stmt.Iterable is not ast.Identifier. got=%T", stmt.Iterable)
	}

	if iterable.Value != "myArray" {
		t.Errorf("iterable is not 'myArray'. got=%s", iterable.Value)
	}

	if len(stmt.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(stmt.Body.Statements))
	}

	bodyStmt, ok := stmt.Body.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ReturnStatement. got=%T",
			stmt.Body.Statements[0])
	}

	if !testIdentifier(t, bodyStmt.ReturnValue, "v") {
		return
	}
}
