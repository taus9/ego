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
		fmt.Println("No file provided.")
		os.Exit(1)
	}

	filename := userArgs[0]
	// fmt.Printf("Running Ego file: %s\n", filename)

	source, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// fmt.Println("file loaded ok, beginning parsing...")

	l := lexer.New(string(source))
	p := parser.New(l)

	program := p.ParseProgram()

	if p.ParseError() != nil {
		printParserErrors(os.Stdout, p.ParseError())
		os.Exit(1)
	}

	env := object.NewEnvironment()
	evalRes := evaluator.Eval(program, env)

	if evalRes != nil && evalRes.Type() == object.ERROR_OBJ {
		fmt.Println(evalRes.Inspect())
	}
	// fmt.Println("Program executed successfully.")
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
			msg = "\n\tStack trace:    -->"
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
	case parser.LET:
		return "LET"
	case parser.RETURN:
		return "RETURN"
	case parser.MAP:
		return "MAP"
	case parser.CALL_ARGS:
		return "FUNCTION CALL ARGUMENTS"
	case parser.FUNC:
		return "FUNCTION"
	default:
		return "UNKNOWN"
	}
}
