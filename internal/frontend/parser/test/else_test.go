package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/lexer"
	"ego/internal/frontend/parser"
	"testing"
)

func TestElseInlineParsing(t *testing.T) {
	input := "5 / 0 else 0"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	elseExpr, ok := stmt.Expression.(*ast.ElseInlineExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ElseInlineExpression. got=%T",
			stmt.Expression)
	}

	if elseExpr.TryExpression.String() != "(5 / 0)" {
		t.Errorf("elseExpr.TryExpression is not '5 / 0'. got=%q",
			elseExpr.TryExpression.String())
	}

	if elseExpr.ElseExpression.String() != "0" {
		t.Errorf("elseExpr.ElseExpression is not '0'. got=%q",
			elseExpr.ElseExpression.String())
	}
}

func TestElseBlockParsing(t *testing.T) {
	input := "5 / 0 else \n ret 0 \n ;"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	elseExpr, ok := stmt.Expression.(*ast.ElseBlockExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ElseBlockExpression. got=%T",
			stmt.Expression)
	}

	if elseExpr.TryExpression.String() != "(5 / 0)" {
		t.Errorf("elseExpr.TryExpression is not '5 / 0'. got=%q",
			elseExpr.TryExpression.String())
	}

	if len(elseExpr.ElseBlock.Statements) != 1 {
		t.Fatalf("elseExpr.ElseBlock.Statements does not contain 1 statement. got=%d",
			len(elseExpr.ElseBlock.Statements))
	}

	retStmt, ok := elseExpr.ElseBlock.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("elseExpr.ElseBlock.Statements[0] is not ast.ReturnStatement. got=%T",
			elseExpr.ElseBlock.Statements[0])
	}

	if retStmt.ReturnValue.String() != "0" {
		t.Errorf("retStmt.ReturnValue is not '0'. got=%q",
			retStmt.ReturnValue.String())
	}
}
