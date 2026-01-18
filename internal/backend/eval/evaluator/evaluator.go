package evaluator

import (
	"ego/internal/backend/eval/object"
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
	"fmt"
	"slices"
)

var reservedWords = []string{
	"nil",
}

var currentToken token.Token

var (
	NIL   = &object.Nil{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.BlockStatement:
		currentToken = node.Token
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		currentToken = node.Token
		return evalIfExpression(node, env)

	case *ast.ReturnStatement:
		currentToken = node.Token
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.ExpressionStatement:
		currentToken = node.Token
		return Eval(node.Expression, env)

	case *ast.PrefixExpression:
		currentToken = node.Token
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		currentToken = node.Token
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.DeclareStatement:
		currentToken = node.Token
		if isReservedWord(node.Name.Value) {
			return newError("cannot use reserved word as identifier: %s", node.Name.Value)
		}
		if env.Exists(node.Name.Value) {
			return newError("identifier already declared in current scope: %s", node.Name.Value)
		}
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.AssignStatement:
		currentToken = node.Token
		var assignName string
		if node.Name != nil {
			// simple variable assignment
			assignName = node.Name.Value
			if !env.Exists(assignName) {
				return newError("identifier not declared: %s", assignName)
			}

			if !node.Name.Mutable {
				return newError("cannot assign to constant: %s", assignName)
			}

			val := Eval(node.Value, env)
			if isError(val) {
				return val
			}
			env.Set(assignName, val)

		} else if node.Index != nil {
			// array index assignment
			indexExp, ok := node.Index.(*ast.IndexExpression)
			if !ok {
				return newError("invalid assignment target")
			}
			assignName = indexExp.Left.(*ast.Identifier).Value

			if !env.Exists(assignName) {
				return newError("identifier not declared: %s", assignName)
			}

			indexObj := Eval(indexExp.Index, env)
			if isError(indexObj) {
				return indexObj
			}
			leftObj := Eval(indexExp.Left, env)
			if isError(leftObj) {
				return leftObj
			}

			switch iteratorObj := leftObj.(type) {
			case *object.Map:
				// map index assignment
				key := indexObj
				valueObj := Eval(node.Value, env)
				if isError(valueObj) {
					return valueObj
				}
				iteratorObj.Pairs[key.(object.Hashable).HasKey()] = object.MapPair{Key: key, Value: valueObj}
				env.Set(assignName, iteratorObj)

			case *object.Array:
				// array index assignment
				indexInt, ok := indexObj.(*object.Integer)
				if !ok {
					return newError("array index must be an integer")
				}
				if indexInt.Value < 0 || indexInt.Value >= int64(len(iteratorObj.Elements)) {
					return newError("array index out of bounds: %d", indexInt.Value)
				}
				valueObj := Eval(node.Value, env)
				if isError(valueObj) {
					return valueObj
				}
				iteratorObj.Elements[indexInt.Value] = valueObj
				env.Set(assignName, iteratorObj)

			default:
				return newError("index assignment not supported for type: %s", leftObj.Type())
			}
		}
	case *ast.ForWhileStatement:
		currentToken = node.Token
		conditionObj := Eval(node.Condition, env)
		if isError(conditionObj) {
			return conditionObj
		}

		for isTruthy(conditionObj) {
			evalResult := Eval(node.Body, env)
			if isError(evalResult) {
				return evalResult
			}
			if _, ok := evalResult.(*object.Break); ok {
				break
			}
			if _, ok := evalResult.(*object.Again); ok {
				continue
			}
			conditionObj = Eval(node.Condition, env)
			if isError(conditionObj) {
				return conditionObj
			}
		}

	case *ast.ForToStatement:
		currentToken = node.Token
		startObj := Eval(node.Start, env)
		if isError(startObj) {
			return startObj
		}
		endObj := Eval(node.End, env)
		if isError(endObj) {
			return endObj
		}

		startInt, ok1 := startObj.(*object.Integer)
		endInt, ok2 := endObj.(*object.Integer)
		if !ok1 || !ok2 {
			return newError("FOR loop bounds must be integers")
		}

		if startInt.Value <= endInt.Value {
			for i := startInt.Value; i <= endInt.Value; i++ {
				env.Set(node.Iterator.Value, &object.Integer{Value: i})
				evalResult := Eval(node.Body, env)
				if isError(evalResult) {
					return evalResult
				}
				if _, ok := evalResult.(*object.Break); ok {
					break
				}
				if _, ok := evalResult.(*object.Again); ok {
					continue
				}
			}
		} else {
			for i := startInt.Value; i >= endInt.Value; i-- {
				env.Set(node.Iterator.Value, &object.Integer{Value: i})
				evalResult := Eval(node.Body, env)
				if isError(evalResult) {
					return evalResult
				}
				if _, ok := evalResult.(*object.Break); ok {
					break
				}
				if _, ok := evalResult.(*object.Again); ok {
					continue
				}
			}
		}

	case *ast.ForInStatement:
		currentToken = node.Token
		iterableObj := Eval(node.Iterable, env)
		if isError(iterableObj) {
			return iterableObj
		}

		switch iterable := iterableObj.(type) {
		case *object.Array:
			for idx, element := range iterable.Elements {
				env.Set(node.Index.Value, &object.Integer{Value: int64(idx)})
				env.Set(node.Value.Value, element)
				evalResult := Eval(node.Body, env)
				if isError(evalResult) {
					return evalResult
				}
				if _, ok := evalResult.(*object.Break); ok {
					break
				}
				if _, ok := evalResult.(*object.Again); ok {
					continue
				}
			}
		case *object.Map:
			for _, pair := range iterable.Pairs {
				env.Set(node.Index.Value, pair.Key)
				env.Set(node.Value.Value, pair.Value)
				evalResult := Eval(node.Body, env)
				if isError(evalResult) {
					return evalResult
				}
				if _, ok := evalResult.(*object.Break); ok {
					break
				}
				if _, ok := evalResult.(*object.Again); ok {
					continue
				}
			}
		default:
			return newError("cannot iterate over type: %s", iterableObj.Type())
		}

	case *ast.BreakStatement:
		currentToken = node.Token
		return &object.Break{}

	case *ast.AgainStatement:
		currentToken = node.Token
		return &object.Again{}

	case *ast.FunctionStatement:
		currentToken = node.Token
		function := &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}
		env.Set(node.Name.Value, function)

	case *ast.AnonymousFunction:
		currentToken = node.Token
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}

	case *ast.CallExpression:
		currentToken = node.Token
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

	case *ast.IndexExpression:
		currentToken = node.Token
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)

	case *ast.MapLiteral:
		currentToken = node.Token
		return evalMapLiteral(node, env)

	case *ast.Identifier:
		currentToken = node.Token
		return evalIdentifier(node, env)

	case *ast.StringLiteral:
		currentToken = node.Token
		return &object.String{Value: node.Value}

	case *ast.FloatLiteral:
		currentToken = node.Token
		return &object.Float{Value: node.Value}

	case *ast.IntegerLiteral:
		currentToken = node.Token
		return &object.Integer{Value: node.Value}

	case *ast.ArrayLiteral:
		currentToken = node.Token
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.Boolean:
		currentToken = node.Token
		return nativeBoolToBooleanObject(node.Value)

	}

	return nil
}

func InitReservedValues(env *object.Environment) {
	env.Set("nil", &object.Nil{})
}

func isReservedWord(word string) bool {
	_, ok := builtins[word]
	return slices.Contains(reservedWords, word) || ok
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NIL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case operator == "and":
		return nativeBoolToBooleanObject(isTruthy(left) && isTruthy(right))

	case operator == "or":
		return nativeBoolToBooleanObject(isTruthy(left) || isTruthy(right))

	case left.Type() == object.FLOAT_OBJ || right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)

	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)

	case operator == "==":
		return nativeBoolToBooleanObject(left == right)

	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	case left.Type() != right.Type():

		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := toFloatValue(left)
	rightVal := toFloatValue(right)

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func toFloatValue(obj object.Object) float64 {
	switch v := obj.(type) {
	case *object.Integer:
		return float64(v.Value)
	case *object.Float:
		return v.Value
	default:
		return 0.0
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NIL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NIL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ ||
				rt == object.ERROR_OBJ ||
				rt == object.BREAK_OBJ ||
				rt == object.AGAIN_OBJ {
				return result
			}
		}
	}

	return result
}

func newError(format string, a ...interface{}) *object.Error {
	ferr := fmt.Sprintf(format, a...)
	loc := fmt.Sprintf("line %d, column %d", currentToken.Span.Line, currentToken.Span.Column)

	msg := "\tYikes!\n\tRuntime Error:  " + ferr + "\n\tError Location: " + loc

	return &object.Error{Message: msg}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: %s", node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: STRING %s STRING", operator)
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	return &object.String{Value: leftVal + rightVal}
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.MAP_OBJ:
		return evalMapIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return newError("array index out of bounds: %d", idx)
	}

	return arrayObject.Elements[idx]
}

func evalMapLiteral(node *ast.MapLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.MapPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashableKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as map key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashableKey.HasKey()
		pairs[hashed] = object.MapPair{Key: key, Value: value}
	}

	return &object.Map{Pairs: pairs}
}

func evalMapIndexExpression(mapObj, index object.Object) object.Object {
	mapObject := mapObj.(*object.Map)

	hashableIndex, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as map key: %s", index.Type())
	}

	hashed := hashableIndex.HasKey()
	pair, ok := mapObject.Pairs[hashed]
	if !ok {
		return NIL
	}

	return pair.Value
}
