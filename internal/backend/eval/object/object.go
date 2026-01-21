package object

import (
	"bytes"
	"ego/internal/frontend/ast"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
)

type ObjectType string

const (
	RETURN_VALUE_OBJ    = "return_value"
	INTEGER_OBJ         = "integer"
	FLOAT_OBJ           = "float"
	BOOLEAN_OBJ         = "boolean"
	NIL_OBJ             = "nil"
	UNHANDLED_ERROR_OBJ = "unhandled_error"
	ERROR_OBJ           = "error"
	FUNCTION_OBJ        = "function"
	STRING_OBJ          = "string"
	BUILTIN_OBJ         = "builtin"
	ARRAY_OBJ           = "array"
	MAP_OBJ             = "map"
	BREAK_OBJ           = "break"
	AGAIN_OBJ           = "again"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

type Float struct {
	Value float64
}

func (f *Float) Inspect() string {
	s := strconv.FormatFloat(f.Value, 'g', -1, 64)

	// If it looks like an integer, force ".0"
	if !strings.ContainsAny(s, ".eE") {
		s += ".0"
	}

	return s
}
func (f *Float) Type() ObjectType { return FLOAT_OBJ }

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

type Nil struct{}

func (n *Nil) Inspect() string  { return "nil" }
func (n *Nil) Type() ObjectType { return NIL_OBJ }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

type UnhandledError struct {
	Message          string
	FormattedMessage string
}

func (ue *UnhandledError) Inspect() string  { return ue.Message }
func (ue *UnhandledError) Type() ObjectType { return UNHANDLED_ERROR_OBJ }

type Error struct {
	Pairs map[HashKey]MapPair
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range e.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) Exists(name string) bool {
	_, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Exists(name)
	}
	return ok
}

func (e *Environment) Delete(name string) {
	delete(e.store, name)
}

func (e *Environment) GetStore() map[string]Object {
	return e.store
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(":>")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(") \n")
	out.WriteString(f.Body.String())
	out.WriteString("\n;")

	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Boolean) HasKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HasKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HasKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type MapPair struct {
	Key   Object
	Value Object
}

type Map struct {
	Pairs map[HashKey]MapPair
}

func (m *Map) Type() ObjectType { return ERROR_OBJ }
func (e *Map) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range e.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type Hashable interface {
	HasKey() HashKey
}

type Break struct{}

func (b *Break) Type() ObjectType { return BREAK_OBJ }
func (b *Break) Inspect() string  { return "break" }

type Again struct{}

func (a *Again) Type() ObjectType { return AGAIN_OBJ }
func (a *Again) Inspect() string  { return "again" }
