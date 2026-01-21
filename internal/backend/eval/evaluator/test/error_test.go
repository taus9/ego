package evaluator

import (
	"ego/internal/backend/eval/object"
	"testing"
)

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true \n 5",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5 \n true + false \n 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if 10 > 1 \n  ret true + false \n ;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`if 10 > 1
				if 10 > 1
					ret true + false
				;
				ret 1
			;`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`'Hello' - 'World'`,
			"unknown operator: STRING - STRING",
		},
		{
			"{'name': 'ego'}[:> x \n x \n ;]",
			"unusable as map key: FUNCTION",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.UnhandledError)
		if !ok {
			t.Errorf("no error object returned. got=%T (%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. got=%q, want=%q", errObj.Message, tt.expectedMessage)
		}
	}
}
