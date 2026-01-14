package evaluator

import "testing"

func TestClosures(t *testing.T) {
	input := `
		newAdder = :> x
			:> y
				x + y
			;
		;
		addTwo = newAdder(2)
		addTwo(2)
	`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 4)
}
