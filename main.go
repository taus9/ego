package main

import (
	"bytes"
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
		// assume it's a filename
		filename := userArgs[0]
		source, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}
		run(string(source))
	}
}

func run(source string) {
	l := lexer.New(string(source))
	p := parser.New(l)

	program := p.ParseProgram()

	if p.ParseError() != nil {
		printParserErrors(os.Stdout, p.ParseError())
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evaluator.InitReservedValues(env)
	evalRes := evaluator.Eval(program, env)

	if evalRes != nil && evalRes.Type() == object.UNHANDLED_ERROR_OBJ {
		fmt.Println(evalRes.Inspect())
	}
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

func printParserErrors(out io.Writer, parseError *parser.ParseError) {
	var buf bytes.Buffer
	buf.WriteString("\tNope!\n")
	fmt.Fprintf(&buf, "\tParser Error:   %s", parseError.Message)

	span := parseError.Token.Span
	fmt.Fprintf(&buf, "\n\tToken Location: line %d, column %d", span.Line, span.Column)

	stack := parseError.StackTrace
	size := stack.Size()
	for i := size - 1; i >= 0; i-- {
		element := stack.Elements()[i]
		var msg string
		switch i {
		case size - 1:
			msg = "\n\tParser Trace:   -->"
		default:
			msg = "\t\t   "
		}

		buf.WriteString("\t")
		fmt.Fprintf(&buf, "%s %s", msg, stackTraceItemToString(element.(int)))
		buf.WriteString("\n")
	}

	io.WriteString(out, buf.String())
}

func stackTraceItemToString(item int) string {
	switch item {
	case parser.PROGRAM:
		return "PROGRAM"
	case parser.ANON_FUNCTION:
		return "ANONYMOUS FUNCTION"
	case parser.ARRAY:
		return "ARRAY"
	case parser.BLOCK:
		return "BLOCK"
	case parser.IF:
		return "IF"
	case parser.DECLARE:
		return "DECLARE"
	case parser.RETURN:
		return "RETURN"
	case parser.MAP:
		return "MAP"
	case parser.CALL_ARGS:
		return "FUNCTION CALL ARGUMENTS"
	case parser.FUNC:
		return "FUNCTION"
	case parser.FOR:
		return "FOR LOOP"
	case parser.ASSIGN:
		return "ASSIGNMENT"
	default:
		return "UNKNOWN"
	}
}
