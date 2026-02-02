package evaluator

import (
	"ego/internal/backend/eval/object"
	"fmt"
	"strconv"
)

func newError(format string, a ...any) *object.UnhandledError {
	return &object.UnhandledError{
		Message: fmt.Sprintf(format, a...),
		Token:   currentToken,
	}
}

func newExecError(message string, statusCode int) *object.UnhandledError {
	return &object.UnhandledError{
		Message:    message,
		Token:      currentToken,
		ErrorType:  "exec",
		StatusCode: statusCode,
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.UNHANDLED_ERROR_OBJ
	}
	return false
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

func createErrorObject(unhandledError *object.UnhandledError) *object.Error {
	pairs := make(map[object.HashKey]object.MapPair)

	var errorType string
	if unhandledError.ErrorType != "" {
		errorType = unhandledError.ErrorType
	} else {
		errorType = "runtime"
	}

	typeKey := &object.String{Value: "type"}
	typeValue := &object.String{Value: errorType}

	hashed := typeKey.HasKey()
	pairs[hashed] = object.MapPair{Key: typeKey, Value: typeValue}

	messageKey := &object.String{Value: "message"}
	messageValue := &object.String{Value: unhandledError.Message}

	hashed = messageKey.HasKey()
	pairs[hashed] = object.MapPair{Key: messageKey, Value: messageValue}

	lineKey := &object.String{Value: "line"}
	lineValue := &object.String{Value: strconv.Itoa(currentToken.Span.Line)}

	hashed = lineKey.HasKey()
	pairs[hashed] = object.MapPair{Key: lineKey, Value: lineValue}

	columnKey := &object.String{Value: "column"}
	columnValue := &object.String{Value: strconv.Itoa(currentToken.Span.Column)}

	hashed = columnKey.HasKey()
	pairs[hashed] = object.MapPair{Key: columnKey, Value: columnValue}

	if unhandledError.StatusCode != 0 {
		statusKey := &object.String{Value: "code"}
		statusValue := &object.Integer{Value: int64(unhandledError.StatusCode)}

		hashed = statusKey.HasKey()
		pairs[hashed] = object.MapPair{Key: statusKey, Value: statusValue}
	}

	return &object.Error{Pairs: pairs}
}
