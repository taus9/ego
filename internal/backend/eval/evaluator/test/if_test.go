package evaluator

import "testing"

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if 1 \n  10 \n ;", 10},
		{"if true \n  10 \n ;", 10},
		{"if false \n  10 \n ;", nil},
		{"if 1 < 2 \n  10 \n ;", 10},
		{"if 1 > 2 \n  10 \n ;", nil},
		{"if 1 > 2 \n  10 \n else \n  20 \n ;", 20},
		{"if 1 < 2 \n  10 \n else \n  20 \n ;", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNilObject(t, evaluated)
		}
	}
}
