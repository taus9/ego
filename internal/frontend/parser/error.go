package parser

import "ego/internal/frontend/token"

type ParseError struct {
	Message    string
	StackTrace *Stack
	Token      token.Token
}

func NewParseError(message string, stackTrace *Stack, tok token.Token) *ParseError {
	return &ParseError{
		Message:    message,
		StackTrace: stackTrace,
		Token:      tok,
	}
}

func (e *ParseError) Error() string {
	return e.Message
}
