package evaluator

import (
	"ego/internal/backend/eval/object"
	"fmt"
	"strings"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to 'len' not supported, got %s", args[0].Type())
			}
		},
	},

	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("first argument to 'push' must be ARRAY, got %s", args[0].Type())
			}

			array := args[0].(*object.Array)
			length := len(array.Elements)

			newElements := make([]object.Object, length+1)
			copy(newElements, array.Elements)
			newElements[length] = args[1]

			return &object.Array{Elements: newElements}
		},
	},

	"put": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NIL
		},
	},

	"str": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			return &object.String{Value: args[0].Inspect()}
		},
	},

	"type": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			return &object.String{Value: string(args[0].Type())}
		},
	},

	"join": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("first argument to 'join' must be ARRAY, got %s", args[0].Type())
			}

			if args[1].Type() != object.STRING_OBJ {
				return newError("second argument to 'join' must be STRING, got %s", args[1].Type())
			}

			array := args[0].(*object.Array)
			sep := args[1].(*object.String).Value

			strElements := make([]string, len(array.Elements))
			for i, elem := range array.Elements {
				strElements[i] = elem.Inspect()
			}

			var result strings.Builder
			for i, strElem := range strElements {
				result.WriteString(strElem)
				if i < len(strElements)-1 {
					result.WriteString(sep)
				}
			}

			return &object.String{Value: result.String()}
		},
	},

	"change": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 3 {
				return newError("wrong number of arguments. got=%d, want=3", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("first argument to 'change' must be ARRAY, got %s", args[0].Type())
			}

			array := args[0].(*object.Array)

			if args[1].Type() != object.INTEGER_OBJ {
				return newError("second argument to 'change' must be INTEGER, got %s", args[1].Type())
			}

			index := args[1].(*object.Integer).Value
			if index > int64(len(array.Elements)-1) || index < 0 {
				return newError("index out of bounds: %d for array of length %d", index, len(array.Elements))
			}

			newElements := make([]object.Object, len(array.Elements))
			copy(newElements, array.Elements)
			newElements[index] = args[2]

			return &object.Array{Elements: newElements}
		},
	},
}
