package evaluator

import (
	"ego/internal/backend/eval/object"
	"testing"
)

func TestAnonFunctionObject(t *testing.T) {
	input := ":> x \n x + 2 \n ;"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong number of parameters. got=%d", len(fn.Parameters))
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"identity = :> x \n x \n ; \n identity(5)", 5},
		{"identity = :> x \n ret x \n ; \n identity(5)", 5},
		{"double = :> x \n x * 2 \n ; \n double(5)", 10},
		{"add = :> x, y \n x + y \n ; \n add(5, 5)", 10},
		{"add = :> x, y \n x + y \n ; \n add(5 + 5, add(5, 5))", 20},
		{":> x \n x \n ;(5)", 5},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}
