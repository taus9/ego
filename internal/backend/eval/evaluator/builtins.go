package evaluator

import (
	"ego/internal/backend/eval/object"
	"fmt"
	"maps"
	"os"
)

var builtins = map[string]*object.Builtin{

	"append": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("first argument must be ARRAY, got %s", args[0].Type())
			}

			array := args[0].(*object.Array)
			length := len(array.Elements)

			newElements := make([]object.Object, length+1)
			copy(newElements, array.Elements)
			newElements[length] = args[1]

			return &object.Array{Elements: newElements}
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
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			switch dstObject := args[0].(type) {
			case *object.Array:
				srcArray, ok := args[1].(*object.Array)
				if !ok {
					return newError("second argument must be ARRAY, got %s", args[1].Type())
				}

				newElements := make([]object.Object, len(srcArray.Elements))
				copy(newElements, srcArray.Elements)
				dstObject.Elements = newElements

				return &object.Integer{Value: int64(len(newElements))}

			case *object.Map:
				srcMap, ok := args[1].(*object.Map)
				if !ok {
					return newError("second argument must be MAP, got %s", args[1].Type())
				}

				clear(dstObject.Pairs)
				maps.Copy(dstObject.Pairs, srcMap.Pairs)

				return &object.Integer{Value: int64(len(dstObject.Pairs))}

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
			_, ok = mapObj.Pairs[hashableKey.HasKey()]
			delete(mapObj.Pairs, hashableKey.HasKey())

			return &object.Boolean{Value: ok}
		},
	},

	"error": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to must be STRING, got %s", args[0].Type())
			}

			strObj := args[0].(*object.String)
			pairs := make(map[object.HashKey]object.MapPair)

			messageKey := &object.String{Value: "message"}
			messageValue := &object.String{Value: strObj.Value}

			hashed := messageKey.HasKey()
			pairs[hashed] = object.MapPair{Key: messageKey, Value: messageValue}

			lineKey := &object.String{Value: "line"}
			lineValue := NIL

			hashed = lineKey.HasKey()
			pairs[hashed] = object.MapPair{Key: lineKey, Value: lineValue}

			columnKey := &object.String{Value: "column"}
			columnValue := NIL

			hashed = columnKey.HasKey()
			pairs[hashed] = object.MapPair{Key: columnKey, Value: columnValue}

			return &object.Error{Pairs: pairs}
		},
	},

	"exit": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) > 1 {
				return newError("wrong number of arguments. got=%d, want=0 or 1", len(args))
			}

			exitCode := 0

			if len(args) == 1 {
				if args[0].Type() != object.INTEGER_OBJ {
					return newError("argument to 'exit' must be INTEGER, got %s", args[0].Type())
				}
				intObj := args[0].(*object.Integer)
				exitCode = int(intObj.Value)
			}

			os.Exit(exitCode)
			return NIL
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

	"input": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) > 1 {
				return newError("wrong number of arguments. got=%d, want=0 or 1", len(args))
			}

			if len(args) == 1 {
				print(args[0].Inspect())
			}

			var input string
			_, err := fmt.Scanln(&input)
			if err != nil {
				return newError("error reading input: %s", err.Error())
			}

			return &object.String{Value: input}
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

	"put": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				println(arg.Inspect())
			}
			return NIL
		},
	},

	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				print(arg.Inspect())
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
}
