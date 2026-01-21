package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/lexer"
	"ego/internal/frontend/parser"
	"testing"
)

func TestIdentifierAssignStatement(t *testing.T) {
	input := "a = 42"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.AssignStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Name.Value != "a" {
		t.Errorf("assign name is not 'a'. got=%s", stmt.Name.Value)
	}

	if stmt.Index != nil {
		t.Errorf("assign index is not nil. got=%T", stmt.Index)
	}

	if !testLiteralExpression(t, stmt.Value, 42) {
		return
	}
}

func TestIndexAssignStatement(t *testing.T) {
	input := "arr[1 + 2] = 42"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.AssignStatement. got=%T",
			program.Statements[0])
	}

	indexExp, ok := stmt.Index.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("stmt.Index is not ast.IndexExpression. got=%T",
			stmt.Index)
	}

	if !testIdentifier(t, indexExp.Left, "arr") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 2) {
		return
	}

	if !testLiteralExpression(t, stmt.Value, 42) {
		return
	}
}

// func TestMultipleIndexAssignStatement(t *testing.T) {
// 	input := "matrix[1][2] = 99"

// 	l := lexer.New(input)
// 	p := parser.New(l)
// 	program := p.ParseProgram()
// 	checkParserErrors(t, p)

// 	if len(program.Statements) != 1 {
// 		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
// 			1, len(program.Statements))
// 	}

// 	stmt, ok := program.Statements[0].(*ast.AssignStatement)
// 	if !ok {
// 		t.Fatalf("program.Statements[0] is not ast.AssignStatement. got=%T",
// 			program.Statements[0])
// 	}

// 	firstIndexExp, ok := stmt.Index.(*ast.IndexExpression)
// 	if !ok {
// 		t.Fatalf("stmt.Index is not ast.IndexExpression. got=%T",
// 			stmt.Index)
// 	}

// 	if !testIdentifier(t, firstIndexExp.Left, "matrix") {
// 		return
// 	}

// }
