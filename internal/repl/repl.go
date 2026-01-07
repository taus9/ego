package repl

import (
	"bufio"
	"ego/internal/backend/eval/evaluator"
	"ego/internal/backend/eval/object"
	"ego/internal/frontend/lexer"
	"ego/internal/frontend/parser"
	"fmt"
	"io"
	"strings"
)

const PROMT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	blockCode := []string{}
	for {
		fmt.Fprintf(out, PROMT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		allLines := ""
		line := scanner.Text() + "\n"

		if len(blockCode) > 0 {
			allLines = strings.Join(blockCode, "") + line
		} else {
			allLines = line
		}

		l := lexer.New(allLines)
		p := parser.New(l)

		program := p.ParseProgram()

		if containsOpenBlockError(p.ParseError().Message) {
			blockCode = append(blockCode, line)
			continue
		}

		if p.ParseError() != nil {
			printParserErrors(out, p.ParseError().Message)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

		blockCode = []string{}
	}
}

func printParserErrors(out io.Writer, message string) {
	io.WriteString(out, "\t"+message+"\n")
}

func containsOpenBlockError(message string) bool {
	if message == "unexpected end of file, missing end block" {
		return true
	}
	return false
}
