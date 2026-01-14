package evaluator

import (
	"ego/internal/backend/eval/evaluator"
	"ego/internal/backend/eval/object"
	"testing"
)

func TestHashLiterals(t *testing.T) {

	input := "two = 'two' \n {'one': 1, two: 2, 'three': 3, 4: 4, true: 5, false: 6}"

	evaluated := testEval(input)
	hash, ok := evaluated.(*object.Map)
	if !ok {
		t.Fatalf("Eval didn't return Map. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HasKey():   1,
		(&object.String{Value: "two"}).HasKey():   2,
		(&object.String{Value: "three"}).HasKey(): 3,
		(&object.Integer{Value: 4}).HasKey():      4,
		evaluator.TRUE.HasKey():                   5,
		evaluator.FALSE.HasKey():                  6,
	}

	if len(hash.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong number of pairs. got=%d, want=%d", len(hash.Pairs), len(expected))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := hash.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
			continue
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestMapIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"{'foo': 5}['foo']",
			5,
		},
		{
			"{'foo': 5}['bar']",
			nil,
		},
		{
			"key = 'foo' \n {'foo': 5}[key]",
			5,
		},
		{
			"{}['foo']",
			nil,
		},
		{
			"{5: 5}[5]",
			5,
		},
		{
			"{true: 5}[true]",
			5,
		},
		{
			"{false: 5}[false]",
			5,
		},
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
