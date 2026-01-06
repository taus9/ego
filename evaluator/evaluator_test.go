package evaluator

import (
	"ego/lexer"
	"ego/object"
	"ego/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 /3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evalutated := testEval(tt.input)
		testIntegerObject(t, evalutated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evalutated := testEval(tt.input)
		testBooleanObject(t, evalutated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evalutated := testEval(tt.input)
		testBooleanObject(t, evalutated, tt.expected)
	}
}

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

func testNilObject(t *testing.T, obj object.Object) bool {
	if obj != NIL {
		t.Errorf("object is not Nil. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

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
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T (%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. got=%q, want=%q", errObj.Message, tt.expectedMessage)
		}
	}
}

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

func TestStringLiteral(t *testing.T) {
	input := `'Hello, World!'`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello, World!" {
		t.Errorf("string has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `'Hello, ' + 'World!'`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello, World!" {
		t.Errorf("string has wrong value. got=%q", str.Value)
	}
}

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
			errObj, ok := evaluated.(*object.Error)
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

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	array, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("array has wrong number of elements. got=%d", len(array.Elements))
	}

	testIntegerObject(t, array.Elements[0], 1)
	testIntegerObject(t, array.Elements[1], 4)
	testIntegerObject(t, array.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"i = 0 \n [1][i]",
			1,
		},
		{
			"[1, 2, 3][1 + 1]",
			3,
		},
		{
			"myArray = [1, 2, 3] \n myArray[2]",
			3,
		},
		{
			"myArray = [1, 2, 3] \n myArray[0] + myArray[1] + myArray[2]",
			6,
		},
		{
			"myArray = [1, 2, 3] \n i = myArray[0] \n myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
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
			errObj, ok := evaluated.(*object.Error)
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
		TRUE.HasKey():                             5,
		FALSE.HasKey():                            6,
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
