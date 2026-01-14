package evaluator

import "testing"

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"ret 10", 10},
		{"ret 10 \n 9", 10},
		{"ret 2 * 5 \n 9", 10},
		{"9 \n ret 2 * 5 \n 9", 10},
		{
			`if 10 > 1
				if 10 > 1
					ret 10
				;
				ret 1
			;`,
			10,
		},
		{
			`f = :> x
				ret x
				x + 10
			;
			f(10)`,
			10,
		},
		{
			`f = :> x
				result = x + 10
				ret result
				ret 10
			;
			f(10)`,
			20,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}
