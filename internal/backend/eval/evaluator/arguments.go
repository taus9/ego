package evaluator

import "fmt"

type ArgParser struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewArgumentParser(input string) *ArgParser {
	ap := &ArgParser{input: input}
	ap.readChar()
	return ap
}

func (ap *ArgParser) readChar() {
	if ap.readPosition >= len(ap.input) {
		ap.ch = 0
	} else {
		ap.ch = ap.input[ap.readPosition]
	}
	ap.position = ap.readPosition
	ap.readPosition += 1
}

func (ap *ArgParser) NextArgument() string {
	ap.skipWhitespace()

	switch ap.ch {
	case '"':
		return ap.readQuotedArgument()
	case 0:
		return ""
	default:
		return ap.readUnquotedArgument()
	}

}

func (ap *ArgParser) readQuotedArgument() string {
	position := ap.position + 1
	for {
		ap.readChar()
		if ap.ch == '"' || ap.ch == 0 {
			break
		}
	}
	arg := ap.input[position:ap.position]
	ap.readChar()

	return fmt.Sprintf("\"%s\"", arg)
}

func (ap *ArgParser) readUnquotedArgument() string {
	position := ap.position
	for {
		if ap.ch == ' ' || ap.ch == '\t' || ap.ch == 0 {
			break
		}
		ap.readChar()
	}
	return ap.input[position:ap.position]
}

func (ap *ArgParser) skipWhitespace() {
	for ap.ch == ' ' || ap.ch == '\t' || ap.ch == '\r' {
		ap.readChar()
	}
}
