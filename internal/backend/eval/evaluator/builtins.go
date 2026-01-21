package evaluator

import (
	"ego/internal/backend/eval/object"
	"fmt"
	"maps"
	"strings"
)

var builtins = map[string]*object.Builtin{
	"abs": {

		// returns the absolute value of an integer or float
		// runtime error for any other type

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

	"clear": {

		// clears all elements from an array or map
		// runtime error for any other type

		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				arg.Elements = arg.Elements[:0]
			case *object.Map:
				clear(arg.Pairs)
			default:
				return newError("invalid argument type, got=%s", args[0].Type())
			}
			return NIL
		},
	},

	"copy": {

		// returns a shallow copy of an array or map
		// runtime error for any other type

		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				newElements := make([]object.Object, len(arg.Elements))
				copy(newElements, arg.Elements)
				return &object.Array{Elements: newElements}
			case *object.Map:
				newPairs := make(map[object.HashKey]object.MapPair)
				maps.Copy(newPairs, arg.Pairs)
				return &object.Map{Pairs: newPairs}
			default:
				return newError("invalid argument type, got=%s", args[0].Type())
			}
		},
	},

	"delete": {

		// deletes a key from a map
		// runtime error if first arg is not a map or if key is not hashable

		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			if args[0].Type() != object.MAP_OBJ {
				return newError("first argument to 'delete' must be MAP, got %s", args[0].Type())
			}

			mapObj := args[0].(*object.Map)

			hashableKey, ok := args[1].(object.Hashable)
			if !ok {
				return newError("unusable as map key: %s", args[1].Type())
			}

			delete(mapObj.Pairs, hashableKey.HasKey())
			return NIL
		},
	},

	"max": {

		// returns the maximum value from a list of numbers
		// runtime error for unsupported types or no arguments

		Fn: func(args ...object.Object) object.Object {
			if len(args) == 0 {
				return newError("max() requires at least one argument")
			}

			var maxValue object.Object = args[0]

			for _, arg := range args[1:] {
				switch max := maxValue.(type) {
				case *object.Integer:
					if arg.Type() != object.INTEGER_OBJ && arg.Type() != object.FLOAT_OBJ {
						return newError("argument to 'max' not supported, got %s", arg.Type())
					}
					if argInt, ok := arg.(*object.Integer); ok && argInt.Value > max.Value {
						maxValue = arg
					}
					if argFloat, ok := arg.(*object.Float); ok && argFloat.Value > float64(max.Value) {
						maxValue = arg
					}
				case *object.Float:
					if arg.Type() != object.FLOAT_OBJ && arg.Type() != object.INTEGER_OBJ {
						return newError("argument to 'max' not supported, got %s", arg.Type())
					}
					if argFloat, ok := arg.(*object.Float); ok && argFloat.Value > max.Value {
						maxValue = arg
					}
					if argInt, ok := arg.(*object.Integer); ok && float64(argInt.Value) > max.Value {
						maxValue = arg
					}
				default:
					return newError("argument to 'max' not supported, got %s", arg.Type())
				}
			}

			return maxValue
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

		// returns a BOOLEAN of TRUE if arg is an ERROR_OBJ
		// returns FALSE for any other type

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
			case *object.Map:
				return &object.Integer{Value: int64(len(arg.Pairs))}
			default:
				return newError("argument to 'len' not supported, got %s", args[0].Type())
			}
		},
	},

	"append": {
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
