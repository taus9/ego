package parser

import (
	"ego/internal/frontend/ast"
	"ego/internal/frontend/parser"
	"testing"
)

func checkParserErrors(t *testing.T, p *parser.Parser) {
	parseErrors := p.ParseError()
	if parseErrors == nil {
		return
	}

	t.Errorf("parser error: %s", parseErrors.Message)
	t.FailNow()
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case float64:
		return testFloatLiteral(t, exp, v)
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}
