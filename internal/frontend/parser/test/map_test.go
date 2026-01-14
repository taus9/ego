package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/lexer"
	"ego/internal/frontend/parser"
	"testing"
)

func TestParsingMapLiteralStringKeys(t *testing.T) {
	input := "{'one': 1, 'two': 2, 'three': 3}"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	mapLiteral, ok := stmt.Expression.(*ast.MapLiteral)
	if !ok {
		t.Fatalf("exp is not ast.MapLiteral. got=%T", stmt.Expression)
	}

	if len(mapLiteral.Pairs) != 3 {
		t.Errorf("mapLiteral.Pairs has wrong length. got=%d", len(mapLiteral.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range mapLiteral.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}

		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyMapLiteral(t *testing.T) {
	input := "{}"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	mapLiteral, ok := stmt.Expression.(*ast.MapLiteral)
	if !ok {
		t.Fatalf("exp is not ast.MapLiteral. got=%T", stmt.Expression)
	}

	if len(mapLiteral.Pairs) != 0 {
		t.Errorf("mapLiteral.Pairs has wrong length. got=%d", len(mapLiteral.Pairs))
	}
}

func TestParsingMapLiteralWithExpressions(t *testing.T) {
	input := "{'one': 0 + 1, 'two': 10 - 8, 'three': 15 / 5}"

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	mapLiteral, ok := stmt.Expression.(*ast.MapLiteral)
	if !ok {
		t.Fatalf("exp is not ast.MapLiteral. got=%T", stmt.Expression)
	}

	if len(mapLiteral.Pairs) != 3 {
		t.Errorf("mapLiteral.Pairs has wrong length. got=%d", len(mapLiteral.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range mapLiteral.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}

		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}

		testFunc(value)
	}
}
