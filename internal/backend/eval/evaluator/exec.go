package evaluator

import (
	"ego/internal/backend/eval/object"
	"ego/internal/frontend/ast"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func evalExecLiteralExpression(node *ast.ExecLiteral, env *object.Environment) object.Object {
	ap := NewArgumentParser(node.Command)

	rawCommand, err := ap.ParseArguments()
	if err != nil {
		return newError("%s", err.Error())
	}

	exeCommand, err := resolveExecCommandTokens(rawCommand, env)
	if err != nil {
		return newError("%s", err.Error())
	}

	command := exeCommand.Command
	args := exeCommand.Args

	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {

		var exitCode int
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 127
		}

		return newExecError(string(output), exitCode)
	}

	return createExecObject(command, args, string(output))
}

func evalExecIndexExpression(execObj, index object.Object) object.Object {
	execObject := execObj.(*object.Exec)

	key, ok := index.(*object.String)
	if !ok {
		return newError("exec index must be a string")
	}

	hashed := key.HasKey()
	pair, ok := execObject.Pairs[hashed]
	if !ok {
		return NIL
	}

	return pair.Value
}

func createExecObject(command string, args []string, output string) *object.Exec {
	pairs := make(map[object.HashKey]object.MapPair)

	commandKey := &object.String{Value: "command"}
	commandValue := &object.String{Value: command}

	hashed := commandKey.HasKey()
	pairs[hashed] = object.MapPair{Key: commandKey, Value: commandValue}

	argsKey := &object.String{Value: "args"}
	argsElements := make([]object.Object, len(args))
	for i, arg := range args {
		argsElements[i] = &object.String{Value: arg}
	}
	argsValue := &object.Array{Elements: argsElements}

	hashed = argsKey.HasKey()
	pairs[hashed] = object.MapPair{Key: argsKey, Value: argsValue}

	outputKey := &object.String{Value: "output"}
	outputValue := &object.String{Value: output}

	hashed = outputKey.HasKey()
	pairs[hashed] = object.MapPair{Key: outputKey, Value: outputValue}

	return &object.Exec{Pairs: pairs}
}

func resolveExecCommandTokens(ec *ExeCommand, env *object.Environment) (*ExeCommand, error) {

	cmd, err := resolveToken(ec.Command, env)
	if err != nil {
		return nil, err
	}
	if cmd == "" {
		return nil, fmt.Errorf("empty command after interpolation")
	}

	args := make([]string, 0, len(ec.Args))
	for _, tok := range ec.Args {
		s, err := resolveToken(tok, env)
		if err != nil {
			return nil, err
		}
		args = append(args, s)
	}

	return &ExeCommand{Command: cmd, Args: args}, nil
}

func resolveToken(token string, env *object.Environment) (string, error) {
	var b strings.Builder
	b.Grow(len(token))

	for i := 0; i < len(token); {
		if token[i] != '{' {
			b.WriteByte(token[i])
			i++
			continue
		}
		// found '{', look for closing '}'
		j := i + 1
		for j < len(token) && token[j] != '}' {
			j++
		}
		if j >= len(token) {
			return "", fmt.Errorf("unterminated identifier in exec token")
		}

		inner := strings.TrimSpace(token[i+1 : j])
		if inner == "" {
			return "", fmt.Errorf("empty identifier in exec token")
		}

		obj := evalIdentifier(&ast.Identifier{Value: inner}, env)
		if isError(obj) {
			return "", fmt.Errorf("identifier not found: %s", inner)
		}

		s, err := toArgString(obj)
		if err != nil {
			return "", err
		}

		b.WriteString(s)
		i = j + 1
	}

	return b.String(), nil
}

func toArgString(obj object.Object) (string, error) {
	switch v := obj.(type) {
	case *object.String:
		return v.Value, nil
	case *object.Integer:
		return strconv.FormatInt(v.Value, 10), nil
	case *object.Float:
		return strconv.FormatFloat(v.Value, 'f', -1, 64), nil
	case *object.Boolean:
		if v.Value {
			return "true", nil
		}
		return "false", nil
	case *object.Nil:
		return "", fmt.Errorf("cannot interpolate nil into exec command")
	default:
		return "", fmt.Errorf("cannot interpolate %s into exec command", obj.Type())
	}
}
