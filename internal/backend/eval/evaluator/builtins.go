package evaluator

import (
	"ego/internal/backend/eval/object"
	"fmt"
	"strings"
)

var builtins = map[string]*object.Builtin{
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				if arg.Value < 0 {
					return &object.Integer{Value: -arg.Value}
				}
				return arg
			case *object.Float:
				if arg.Value < 0 {
					return &object.Float{Value: -arg.Value}
				}
				return arg
			default:
				return newError("argument to 'abs' not supported, got %s", args[0].Type())
			}
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

	"ok": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch args[0].Type() {
			case object.ERROR_OBJ:
				return FALSE
			default:
				return TRUE
			}
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

	"float": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				return &object.Float{Value: float64(arg.Value)}
			case *object.Float:
				return arg
			case *object.String:
				var floatValue float64
				_, err := fmt.Sscanf(arg.Value, "%f", &floatValue)
				if err != nil {
					return newError("cannot convert string to float: %s", arg.Value)
				}
				return &object.Float{Value: floatValue}
			default:
				return newError("argument to 'float' not supported, got %s", args[0].Type())
			}
		},
	},
	"int": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				return arg
			case *object.Float:
				return &object.Integer{Value: int64(arg.Value)}
			case *object.String:
				var intValue int64
				_, err := fmt.Sscanf(arg.Value, "%d", &intValue)
				if err != nil {
					return newError("cannot convert string to int: %s", arg.Value)
				}
				return &object.Integer{Value: intValue}
			default:
				return newError("argument to 'int' not supported, got %s", args[0].Type())
			}
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
}
