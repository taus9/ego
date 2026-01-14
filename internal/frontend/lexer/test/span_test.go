package lexer

import (
	"ego/internal/frontend/lexer"
	"testing"
)

func TestSpanTracking(t *testing.T) {
	input := "five = 5\nten = 10"

	l := lexer.New(input)

	tok := l.NextToken() // five
	if tok.Span.Line != 1 || tok.Span.Column != 1 {
		t.Fatalf("token span wrong. expected=(1,1), got=(%d,%d)",
			tok.Span.Line, tok.Span.Column)
	}

	tok = l.NextToken() // =
	if tok.Span.Line != 1 || tok.Span.Column != 6 {
		t.Fatalf("token span wrong. expected=(1,6), got=(%d,%d)",
			tok.Span.Line, tok.Span.Column)
	}

	tok = l.NextToken() // 5
	if tok.Span.Line != 1 || tok.Span.Column != 8 {
		t.Fatalf("token span wrong. expected=(1,8), got=(%d,%d)",
			tok.Span.Line, tok.Span.Column)
	}

	tok = l.NextToken() // EOL
	if tok.Span.Line != 1 || tok.Span.Column != 9 {
		t.Fatalf("token span wrong. expected=(1,9), got=(%d,%d)",
			tok.Span.Line, tok.Span.Column)
	}

	tok = l.NextToken() // ten
	if tok.Span.Line != 2 || tok.Span.Column != 1 {
		t.Fatalf("token span wrong. expected=(2,1), got=(%d,%d)",
			tok.Span.Line, tok.Span.Column)
	}

	tok = l.NextToken() // =
	if tok.Span.Line != 2 || tok.Span.Column != 5 {
		t.Fatalf("token span wrong. expected=(2,5), got=(%d,%d)",
			tok.Span.Line, tok.Span.Column)
	}

	tok = l.NextToken() // 10
	if tok.Span.Line != 2 || tok.Span.Column != 7 {
		t.Fatalf("token span wrong. expected=(2,7), got=(%d,%d)",
			tok.Span.Line, tok.Span.Column)
	}
}
