package evaluator_test

import (
	"testing"

	"ego/internal/backend/eval/evaluator"
)

func TestArgParser(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:  "basic and quoted",
			input: `arg1 arg2 "arg 3"`,
			// NOTE: quoted args are now returned WITHOUT quotes.
			want: []string{"arg1", "arg2", "arg 3"},
		},
		{
			name:  "quoted then unquoted",
			input: `"quoted arg" unquoted`,
			// NOTE: quoted args are now returned WITHOUT quotes.
			want: []string{"quoted arg", "unquoted"},
		},
		{
			name:  "leading whitespace",
			input: `   leadingWhitespace`,
			want:  []string{"leadingWhitespace"},
		},
		{
			name:    "unterminated quote",
			input:   `"unterminated quote`,
			wantErr: true,
		},
		{
			name:  "brace group single token",
			input: `{interpolated arg}`,
			want:  []string{`{interpolated arg}`},
		},
		{
			name:  "nested braces preserved",
			input: `{nested {interpolation}}`,
			want:  []string{`{nested {interpolation}}`},
		},
		{
			name:    "unterminated interpolation",
			input:   `{unterminated interpolation`,
			wantErr: true,
		},
		{
			name:    "unterminated nested interpolation",
			input:   `{unterminated nested {interpolation}`,
			wantErr: true,
		},
		{
			name:  "mixed spacing",
			input: `arg1   "arg  two"   {arg3}`,
			// NOTE: quoted args are now returned WITHOUT quotes.
			want: []string{"arg1", "arg  two", "{arg3}"},
		},
		{
			name:  "empty",
			input: ``,
			want:  []string{},
		},

		// New: braces inside an unquoted token should NOT split on whitespace inside braces.
		{
			name:  "inline brace group in unquoted token",
			input: `--path={ some file }`,
			want:  []string{`--path={ some file }`},
		},
		{
			name:  "brace group mid-token",
			input: `hi{ ident }there`,
			want:  []string{`hi{ ident }there`},
		},
		{
			name:  "multiple tokens with inline brace group",
			input: `cmd --flag={ a b } tail`,
			want:  []string{`cmd`, `--flag={ a b }`, `tail`},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			parser := evaluator.NewArgumentParser(tt.input)

			var got []string
			for {
				arg, err := parser.NextArgument()
				if err != nil {
					if tt.wantErr {
						return
					}
					t.Fatalf("input: %q, unexpected error: %v", tt.input, err)
				}
				if arg == "" {
					break
				}
				got = append(got, arg)
			}

			if tt.wantErr {
				t.Fatalf("input: %q, expected error, got none (args=%v)", tt.input, got)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("input: %q, expected %d args, got %d args (%v)",
					tt.input, len(tt.want), len(got), got,
				)
			}

			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Fatalf("input: %q, arg[%d]: expected %q, got %q",
						tt.input, i, tt.want[i], got[i],
					)
				}
			}
		})
	}
}
