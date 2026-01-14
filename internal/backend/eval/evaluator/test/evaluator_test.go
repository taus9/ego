package evaluator

import (
	"ego/internal/backend/eval/evaluator"
	"ego/internal/backend/eval/object"
	"ego/internal/frontend/lexer"
	"ego/internal/frontend/parser"
	"testing"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return evaluator.Eval(program, env)
}

func testNilObject(t *testing.T, obj object.Object) bool {
	if obj != evaluator.NIL {
		t.Errorf("object is not Nil. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
