package evaluator

import (
	"ego/internal/backend/eval/object"
	"testing"
)

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len('')`, 0},
		{`len('hello')`, 5},
		{`len('hello world')`, 11},
		{`len(1)`, "argument to 'len' not supported, got INTEGER"},
		{`len('one', 'two')`, "wrong number of arguments. got=2, want=1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.UnhandledError)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. got=%q, want=%q", errObj.Message, expected)
			}
		}
	}
}

func TestBuiltinFunctionPush(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"a = [] \n a = push(a, 1) \n a",
			&object.Array{Elements: []object.Object{
				&object.Integer{Value: 1},
			}},
		},
		{
			"push([1, 2, 3], 4)",
			&object.Array{Elements: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
				&object.Integer{Value: 4},
			}},
		},
		{
			"push([], 1)",
			&object.Array{Elements: []object.Object{
				&object.Integer{Value: 1},
			}},
		},
		{
			"push(1, 2)",
			"first argument to 'push' must be ARRAY, got INTEGER",
		},
		{
			"push([1, 2], 3, 4)",
			"wrong number of arguments. got=3, want=2",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case *object.Array:
			result, ok := evaluated.(*object.Array)
			if !ok {
				t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
			}
			if len(result.Elements) != len(expected.Elements) {
				t.Fatalf("array has wrong number of elements. got=%d, want=%d", len(result.Elements), len(expected.Elements))
			}
			for i, expectedElem := range expected.Elements {
				testIntegerObject(t, result.Elements[i], expectedElem.(*object.Integer).Value)
			}
		case string:
			errObj, ok := evaluated.(*object.UnhandledError)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. got=%q, want=%q", errObj.Message, expected)
			}
		}
	}
}
