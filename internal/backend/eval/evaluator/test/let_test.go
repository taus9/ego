package evaluator

import "testing"

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"a = 5 \n a", 5},
		{"a = 5 * 5 \n a", 25},
		{"a = 5 \n b = a \n b", 5},
		{"a = 5 \n b = a \n c = a + b + 5 \n c", 15},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}
