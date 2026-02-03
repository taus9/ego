package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/lexer"
	"ego/internal/frontend/parser"
	"testing"
)

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestNestedCallExpressionParsing(t *testing.T) {
	input := "a(b(c()))"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "a") {
		return
	}

	if len(exp.Arguments) != 1 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	nestedCall, ok := exp.Arguments[0].(*ast.CallExpression)
	if !ok {
		t.Fatalf("exp.Arguments[0] is not ast.CallExpression. got=%T",
			exp.Arguments[0])
	}

	if !testIdentifier(t, nestedCall.Function, "b") {
		return
	}

	if len(nestedCall.Arguments) != 1 {
		t.Fatalf("wrong length of arguments. got=%d", len(nestedCall.Arguments))
	}

	nestedNestedCall, ok := nestedCall.Arguments[0].(*ast.CallExpression)
	if !ok {
		t.Fatalf("nestedCall.Arguments[0] is not ast.CallExpression. got=%T",
			nestedCall.Arguments[0])
	}

	if !testIdentifier(t, nestedNestedCall.Function, "c") {
		return
	}

	if len(nestedNestedCall.Arguments) != 0 {
		t.Fatalf("wrong length of arguments. got=%d", len(nestedNestedCall.Arguments))
	}
}

func TestMultiArgNestedCall(t *testing.T) {
	input := "a(b(1, 2), c(3, 4))"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "a") {
		return
	}

	if len(exp.Arguments) != 2 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	firstArg, ok := exp.Arguments[0].(*ast.CallExpression)
	if !ok {
		t.Fatalf("exp.Arguments[0] is not ast.CallExpression. got=%T",
			exp.Arguments[0])
	}

	if !testIdentifier(t, firstArg.Function, "b") {
		return
	}

	if len(firstArg.Arguments) != 2 {
		t.Fatalf("wrong length of arguments. got=%d", len(firstArg.Arguments))
	}

	testLiteralExpression(t, firstArg.Arguments[0], 1)
	testLiteralExpression(t, firstArg.Arguments[1], 2)

	secondArg, ok := exp.Arguments[1].(*ast.CallExpression)
	if !ok {
		t.Fatalf("exp.Arguments[1] is not ast.CallExpression. got=%T",
			exp.Arguments[1])
	}

	if !testIdentifier(t, secondArg.Function, "c") {
		return
	}

	if len(secondArg.Arguments) != 2 {
		t.Fatalf("wrong length of arguments. got=%d", len(secondArg.Arguments))
	}

	testLiteralExpression(t, secondArg.Arguments[0], 3)
	testLiteralExpression(t, secondArg.Arguments[1], 4)
}

func TestEmptyArgMultiNestedCall(t *testing.T) {
	input := "a(b(), c())"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "a") {
		return
	}

	if len(exp.Arguments) != 2 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	firstArg, ok := exp.Arguments[0].(*ast.CallExpression)
	if !ok {
		t.Fatalf("exp.Arguments[0] is not ast.CallExpression. got=%T",
			exp.Arguments[0])
	}

	if !testIdentifier(t, firstArg.Function, "b") {
		return
	}

	if len(firstArg.Arguments) != 0 {
		t.Fatalf("wrong length of arguments. got=%d", len(firstArg.Arguments))
	}

	secondArg, ok := exp.Arguments[1].(*ast.CallExpression)
	if !ok {
		t.Fatalf("exp.Arguments[1] is not ast.CallExpression. got=%T",
			exp.Arguments[1])
	}

	if !testIdentifier(t, secondArg.Function, "c") {
		return
	}

	if len(secondArg.Arguments) != 0 {
		t.Fatalf("wrong length of arguments. got=%d", len(secondArg.Arguments))
	}
}
