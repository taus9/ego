package lexer

import (
	"ego/internal/frontend/lexer"
	"ego/internal/frontend/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
	five = 5
	ten = 10
	add = :> x, y
		x + y
	;
	result = add five, ten
	!-/*5
	5 < 10 > 5
	if 5 > 10
		ret true
	;
	ret false
	10 == 10
	10 != 9
	: sub = [0] - [1]
	'foobar'
	'foo bar'
	[1, 2]
	{'foo': 'bar'}
	for i = 0 to 9
		reminder = i % 2
	;
	x = 3.14
	y = 0.001
	z = 10.0
	x and y or z
	x <= y
	y >= z
	a := 42
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.EOL, "\n"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.EOL, "\n"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.EOL, "\n"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.ANON_FUNCTION, ":>"},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.EOL, "\n"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.EOL, "\n"},
		{token.END_BLOCK, ";"},
		{token.EOL, "\n"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.EOL, "\n"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.EOL, "\n"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.EOL, "\n"},
		{token.IF, "if"},
		{token.INT, "5"},
		{token.GT, ">"},
		{token.INT, "10"},
		{token.EOL, "\n"},
		{token.RETURN, "ret"},
		{token.TRUE, "true"},
		{token.EOL, "\n"},
		{token.END_BLOCK, ";"},
		{token.EOL, "\n"},
		{token.RETURN, "ret"},
		{token.FALSE, "false"},
		{token.EOL, "\n"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.EOL, "\n"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.EOL, "\n"},
		{token.COLON, ":"},
		{token.IDENT, "sub"},
		{token.ASSIGN, "="},
		{token.LBRACKET, "["},
		{token.INT, "0"},
		{token.RBRACKET, "]"},
		{token.MINUS, "-"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.RBRACKET, "]"},
		{token.EOL, "\n"},
		{token.STRING, "foobar"},
		{token.EOL, "\n"},
		{token.STRING, "foo bar"},
		{token.EOL, "\n"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.EOL, "\n"},
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.EOL, "\n"},
		{token.FOR, "for"},
		{token.IDENT, "i"},
		{token.ASSIGN, "="},
		{token.INT, "0"},
		{token.TO, "to"},
		{token.INT, "9"},
		{token.EOL, "\n"},
		{token.IDENT, "reminder"},
		{token.ASSIGN, "="},
		{token.IDENT, "i"},
		{token.MOD, "%"},
		{token.INT, "2"},
		{token.EOL, "\n"},
		{token.END_BLOCK, ";"},
		{token.EOL, "\n"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.FLOAT, "3.14"},
		{token.EOL, "\n"},
		{token.IDENT, "y"},
		{token.ASSIGN, "="},
		{token.FLOAT, "0.001"},
		{token.EOL, "\n"},
		{token.IDENT, "z"},
		{token.ASSIGN, "="},
		{token.FLOAT, "10.0"},
		{token.EOL, "\n"},
		{token.IDENT, "x"},
		{token.AND, "and"},
		{token.IDENT, "y"},
		{token.OR, "or"},
		{token.IDENT, "z"},
		{token.EOL, "\n"},
		{token.IDENT, "x"},
		{token.LT_EQ, "<="},
		{token.IDENT, "y"},
		{token.EOL, "\n"},
		{token.IDENT, "y"},
		{token.GT_EQ, ">="},
		{token.IDENT, "z"},
		{token.EOL, "\n"},
		{token.IDENT, "a"},
		{token.DECLARE, ":="},
		{token.INT, "42"},
		{token.EOL, "\n"},
		{token.EOF, ""},
	}

	l := lexer.New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
