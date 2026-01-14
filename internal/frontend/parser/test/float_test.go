package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/lexer"
	"ego/internal/frontend/parser"
	"strconv"
	"strings"
	"testing"
)

func TestFloatLiteralExpression(t *testing.T) {
	input := "5.5"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("exp not *ast.FloatLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5.5 {
		t.Errorf("literal.Value not %f. got=%f", 5.5, literal.Value)
	}
	if literal.TokenLiteral() != "5.5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5.5",
			literal.TokenLiteral())
	}
}

func testFloatLiteral(t *testing.T, exp ast.Expression, value float64) bool {
	lit, ok := exp.(*ast.FloatLiteral)
	if !ok {
		t.Errorf("exp not *ast.FloatLiteral. got=%T", exp)
		return false
	}
	if lit.Value != value {
		t.Errorf("lit.Value not %f. got=%f", value, lit.Value)
		return false
	}
	if lit.TokenLiteral() != formatFloatLiteral(value) {
		t.Errorf("lit.TokenLiteral not %f. got=%s",
			value, lit.TokenLiteral())
		return false
	}
	return true
}

func formatFloatLiteral(value float64) string {
	s := strconv.FormatFloat(value, 'g', -1, 64)

	// If it looks like an integer, force ".0"
	if !strings.ContainsAny(s, ".eE") {
		s += ".0"
	}

	return s
}
