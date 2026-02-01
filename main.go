package main

import (
	"ego/internal/backend/eval/evaluator"
	"ego/internal/backend/eval/object"
	"fmt"
	"os"
)

func main() {
	userArgs := os.Args[1:]
	if len(userArgs) == 0 {
		showUsage()
		os.Exit(0)
	}

	switch userArgs[0] {

	case "version":
		fmt.Println("ego version 0.1.0")
		return

	case "help":
		showUsage()
		return

	default:
		filename := userArgs[0]
		runFile(filename, userArgs[1:])
	}
}

func runFile(entryPath string, args []string) {
	globalEnv := object.NewEnvironment()
	evaluator.InitReservedValues(globalEnv)

	initUserArgs(globalEnv, args)

	res := evaluator.Resolver{
		StdlibRoot: "stdlib",
		HomeDir:    os.Getenv("HOME"),
	}

	ld := evaluator.NewLoader(res)
	ld.GlobalEnv = globalEnv

	_, err := ld.Load(entryPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func initUserArgs(env *object.Environment, args []string) {
	for i, arg := range args {
		argName := fmt.Sprintf("$%d", i+1)
		argObj := &object.String{Value: arg}
		env.Declare(argName, argObj)
	}
	argc := len(args)
	env.Declare("$a", &object.Integer{Value: int64(argc)})
}

func showUsage() {
	fmt.Println(`ego programming language

Usage:
  ego <arguments>

Arguments:
  version   Show the version of Ego
  help      Show this help message
  <file>    Path to the Ego source file to execute

Example:
  ego my_program.ego`)
}
