package evaluator

import (
	"ego/internal/backend/eval/object"
	"fmt"
	"strconv"
	"strings"
)

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

func (ap *ArgParser) NextArgument() (string, error) {
	ap.skipWhitespace()

	switch ap.ch {
	case '"':
		return ap.readQuotedArgument()
	case '{':
		return ap.readInterpolatedArgument()
	case 0:
		return "", nil
	default:
		return ap.readUnquotedArgument()
	}
}

func (ap *ArgParser) readQuotedArgument() (string, error) {
	// return raw value, without surrounding quotes.
	start := ap.position + 1

	for {
		ap.readChar() // skip opening quote

		if ap.ch == 0 {
			return "", fmt.Errorf("unterminated quoted argument")
		}
		if ap.ch == '"' {
			break
		}

		// possibly handle escape sequences here in the future
	}

	arg := ap.input[start:ap.position]
	ap.readChar() // skip closing quote
	return arg, nil
}

func (ap *ArgParser) readInterpolatedArgument() (string, error) {
	return ap.readBraceGroup()
}

func (ap *ArgParser) readUnquotedArgument() (string, error) {
	start := ap.position
	var b strings.Builder
	b.Grow(32)

	for {
		if ap.ch == 0 || ap.ch == ' ' || ap.ch == '\t' || ap.ch == '\r' || ap.ch == '\n' {
			break
		}

		if ap.ch == '{' {
			group, err := ap.readBraceGroup()
			if err != nil {
				return "", err
			}
			b.WriteString(group)
			continue
		}

		b.WriteByte(ap.ch)
		ap.readChar()
	}

	if b.Len() == 0 {
		return ap.input[start:ap.position], nil
	}

	return b.String(), nil
}

func (ap *ArgParser) readBraceGroup() (string, error) {
	start := ap.position

	braces := 0

	for {
		if ap.ch == 0 {
			return "", fmt.Errorf("unterminated interpolated argument")
		}

		switch ap.ch {
		case '{':
			braces++
		case '}':
			braces--
			if braces == 0 {
				// include the closing brace in the slice.
				end := ap.position + 1
				ap.readChar() // skip closing brace
				return ap.input[start:end], nil
			}
		}

		ap.readChar() // continue reading inside braces
	}
}

func (ap *ArgParser) skipWhitespace() {
	for ap.ch == ' ' || ap.ch == '\t' || ap.ch == '\r' || ap.ch == '\n' {
		ap.readChar()
	}
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
	default:
		return "", fmt.Errorf("cannot interpolate %s into exec command", obj.Type())
	}
}
