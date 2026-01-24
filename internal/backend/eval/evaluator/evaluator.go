package evaluator

import (
	"ego/internal/backend/eval/object"
	"ego/internal/frontend/ast"
	"ego/internal/frontend/token"
	"fmt"
	"slices"
	"strconv"
)

var reservedWords = []string{
	"nil",
	"$E",
}

var currentToken token.Token

var (
	NIL          = &object.Nil{}
	TRUE         = &object.Boolean{Value: true}
	FALSE        = &object.Boolean{Value: false}
	ZERO_INT     = &object.Integer{Value: 0}
	ZERO_FLOAT   = &object.Float{Value: 0.0}
	EMPTY_STRING = &object.String{Value: ""}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.BlockStatement:
		currentToken = node.Token
		return evalBlockStatement(node, env)

	case *ast.IfBlockExpression:
		currentToken = node.Token
		return evalIfBlockExpression(node, env)

	case *ast.IfTernaryExpression:
		currentToken = node.Token
		return evalIfTernaryExpression(node, env)

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

		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		if node.Name.Value == "_" {
			// discard value assigned to underscore
			return NIL
		}
		if !env.Declare(node.Name.Value, val) {
			return newError("identifier already declared in current scope: %s", node.Name.Value)
		}

	case *ast.AssignStatement:
		currentToken = node.Token
		var assignName string
		if node.Name != nil {
			// simple variable assignment
			assignName = node.Name.Value

			if !node.Name.Mutable {
				return newError("cannot assign to constant: %s", assignName)
			}

			val := Eval(node.Value, env)
			if isError(val) {
				return val
			}
			if assignName == "_" {
				// discard value assigned to underscore
				return NIL
			}
			if !env.Set(assignName, val) {
				return newError("identifier not declared: %s", assignName)
			}

		} else if node.Index != nil {
			// array index assignment
			indexExp, ok := node.Index.(*ast.IndexExpression)
			if !ok {
				return newError("invalid assignment target")
			}
			assignName = indexExp.Left.(*ast.Identifier).Value

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
				if !env.Set(assignName, iteratorObj) {
					return newError("identifier not declared: %s", assignName)
				}

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
				if !env.Set(assignName, iteratorObj) {
					return newError("identifier not declared: %s", assignName)
				}

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
			evalResult := evalBlockStatement(node.Body, env)
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

		// the iterator must be set in the env outside the for block
		// so that it can be constant through out the whole loop
		// once we are done with the loop it should be deleted
		defer env.Delete(node.Iterator.Value)

		if startInt.Value <= endInt.Value {
			for i := startInt.Value; i <= endInt.Value; i++ {
				env.Set(node.Iterator.Value, &object.Integer{Value: i})
				evalResult := evalBlockStatement(node.Body, env)
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
				evalResult := evalBlockStatement(node.Body, env)
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

		// just like with the for to loop declared iterators must
		// be removed from the env outside the for block
		// once the loop is finished
		defer env.Delete(node.Index.Value)
		defer env.Delete(node.Value.Value)

		switch iterable := iterableObj.(type) {
		case *object.String:
			for idx, ch := range iterable.Value {
				if node.Index.Value != "_" {
					env.Set(node.Index.Value, &object.Integer{Value: int64(idx)})
				}
				if node.Value.Value != "_" {
					env.Set(node.Value.Value, &object.String{Value: string(ch)})
				}

				evalResult := evalBlockStatement(node.Body, env)

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
		case *object.Array:
			for idx, element := range iterable.Elements {
				if node.Index.Value != "_" {
					env.Set(node.Index.Value, &object.Integer{Value: int64(idx)})
				}
				if node.Value.Value != "_" {
					env.Set(node.Value.Value, element)
				}

				evalResult := evalBlockStatement(node.Body, env)

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
				if node.Index.Value != "_" {
					env.Set(node.Index.Value, pair.Key)
				}
				if node.Value.Value != "_" {
					env.Set(node.Value.Value, pair.Value)
				}

				evalResult := evalBlockStatement(node.Body, env)

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
		if !env.Declare(node.Name.Value, function) {
			return newError("function already declared in current scope: %s", node.Name.Value)
		}

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

	case *ast.ElseInlineExpression:
		currentToken = node.Token
		tryObj := Eval(node.TryExpression, env)

		if isError(tryObj) {
			unhandledError, _ := tryObj.(*object.UnhandledError)
			errorObj := createErrorObject(unhandledError)

			blockEnv := object.NewEnclosedEnvironment(env)
			blockEnv.Declare("$E", errorObj)

			result := Eval(node.ElseExpression, blockEnv)

			return result
		}

		return tryObj

	case *ast.ElseBlockExpression:
		currentToken = node.Token
		tryObj := Eval(node.TryExpression, env)

		if isError(tryObj) {
			unhandledError, _ := tryObj.(*object.UnhandledError)
			errorObj := createErrorObject(unhandledError)

			env.Declare("$E", errorObj)
			defer env.Delete("$E")

			result := evalBlockStatement(node.ElseBlock, env)

			return result
		}

		return tryObj

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
	env.Declare("nil", &object.Nil{})
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
		case *object.UnhandledError:
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
	intObj, ok := right.(*object.Integer)
	if ok {
		return nativeBoolToBooleanObject(intObj.Value == 0)
	}

	floatObj, ok := right.(*object.Float)
	if ok {
		return nativeBoolToBooleanObject(floatObj.Value == 0.0)
	}

	strObj, ok := right.(*object.String)
	if ok {
		return nativeBoolToBooleanObject(strObj.Value == "")
	}

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
	if right.Type() == object.INTEGER_OBJ {
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	}
	if right.Type() == object.FLOAT_OBJ {
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	}
	return newError("unknown operator: -%s", right.Type())
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case operator == "and":
		return nativeBoolToBooleanObject(isTruthy(left) && isTruthy(right))

	case operator == "or":
		return nativeBoolToBooleanObject(isTruthy(left) || isTruthy(right))

	case left.Type() == object.FLOAT_OBJ || right.Type() == object.FLOAT_OBJ:
		if left.Type() != object.FLOAT_OBJ && left.Type() != object.INTEGER_OBJ {
			return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
		}
		if right.Type() != object.FLOAT_OBJ && right.Type() != object.INTEGER_OBJ {
			return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
		}
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

func evalIfBlockExpression(ie *ast.IfBlockExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return evalBlockStatement(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return evalBlockStatement(ie.Alternative, env)
	} else {
		return NIL
	}
}

func evalIfTernaryExpression(ie *ast.IfTernaryExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		conditionEnv := object.NewEnclosedEnvironment(env)
		return Eval(ie.Consequence, conditionEnv)
	} else if ie.Alternative != nil {
		conditionEnv := object.NewEnclosedEnvironment(env)
		return Eval(ie.Alternative, conditionEnv)
	} else {
		return NIL
	}
}

func isTruthy(obj object.Object) bool {
	intObj, ok := obj.(*object.Integer)
	if ok {
		return intObj.Value != 0
	}

	floatObj, ok := obj.(*object.Float)
	if ok {
		return floatObj.Value != 0.0
	}

	strObj, ok := obj.(*object.String)
	if ok {
		return strObj.Value != ""
	}

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
	blockEnv := object.NewEnclosedEnvironment(env)
	for _, statement := range block.Statements {
		result = Eval(statement, blockEnv)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ ||
				rt == object.UNHANDLED_ERROR_OBJ ||
				rt == object.BREAK_OBJ ||
				rt == object.AGAIN_OBJ {
				return result
			}
		}
	}

	return result
}

func newError(format string, a ...interface{}) *object.UnhandledError {
	ferr := fmt.Sprintf(format, a...)
	loc := fmt.Sprintf("line %d, column %d", currentToken.Span.Line, currentToken.Span.Column)

	msg := "\tYikes!\n\tRuntime Error:  " + ferr + "\n\tError Location: " + loc

	return &object.UnhandledError{Message: ferr, FormattedMessage: msg}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.UNHANDLED_ERROR_OBJ
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
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index)
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.MAP_OBJ:
		return evalMapIndexExpression(left, index)
	case left.Type() == object.ERROR_OBJ:
		return evalErrorIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalErrorIndexExpression(errObj, index object.Object) object.Object {
	errorObject := errObj.(*object.Error)

	key, ok := index.(*object.String)
	if !ok {
		return newError("$E index must be a string")
	}

	hashed := key.HasKey()
	pair, ok := errorObject.Pairs[hashed]
	if !ok {
		return NIL
	}

	return pair.Value
}

func evalStringIndexExpression(str, index object.Object) object.Object {
	stringObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(stringObject.Value) - 1)

	if idx < 0 || idx > max {
		return newError("string index out of bounds: %d", idx)
	}

	return &object.String{Value: string(stringObject.Value[idx])}
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

func createErrorObject(unhandledError *object.UnhandledError) *object.Error {
	pairs := make(map[object.HashKey]object.MapPair)

	messageKey := &object.String{Value: "message"}
	messageValue := &object.String{Value: unhandledError.Message}

	hashed := messageKey.HasKey()
	pairs[hashed] = object.MapPair{Key: messageKey, Value: messageValue}

	lineKey := &object.String{Value: "line"}
	lineValue := &object.String{Value: strconv.Itoa(currentToken.Span.Line)}

	hashed = lineKey.HasKey()
	pairs[hashed] = object.MapPair{Key: lineKey, Value: lineValue}

	columnKey := &object.String{Value: "column"}
	columnValue := &object.String{Value: strconv.Itoa(currentToken.Span.Column)}

	hashed = columnKey.HasKey()
	pairs[hashed] = object.MapPair{Key: columnKey, Value: columnValue}

	return &object.Error{Pairs: pairs}
}
