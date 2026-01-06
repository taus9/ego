package main

import (
	"ego/internal/backend/eval/evaluator"
	"ego/internal/backend/eval/object"
	"ego/internal/frontend/lexer"
	"ego/internal/frontend/parser"
	"fmt"
	"io"
	"os"
)

func main() {
	userArgs := os.Args[1:]
	if len(userArgs) == 0 {
		fmt.Println("No file provided.")
		os.Exit(1)
	}

	filename := userArgs[0]

	source, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	l := lexer.New(string(source))
	p := parser.New(l)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(os.Stdout, p.Errors())
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluator.Eval(program, env)
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
