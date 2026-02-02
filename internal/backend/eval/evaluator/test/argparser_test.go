package evaluator_test

import (
	"testing"

	"ego/internal/backend/eval/evaluator"
)

func TestArgParser(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantCmd string
		wantArg []string
		wantErr bool
	}{
		{
			name:    "basic and quoted",
			input:   `arg1 arg2 "arg 3"`,
			wantCmd: "arg1",
			wantArg: []string{"arg2", "arg 3"}, // quoted args returned WITHOUT quotes
		},
		{
			name:    "quoted then unquoted",
			input:   `"quoted arg" unquoted`,
			wantCmd: "quoted arg",
			wantArg: []string{"unquoted"},
		},
		{
			name:    "leading whitespace",
			input:   `   leadingWhitespace`,
			wantCmd: "leadingWhitespace",
			wantArg: []string{},
		},
		{
			name:    "unterminated quote",
			input:   `"unterminated quote`,
			wantErr: true,
		},
		{
			name:    "brace group single token",
			input:   `{interpolated arg}`,
			wantCmd: `{interpolated arg}`,
			wantArg: []string{},
		},
		{
			name:    "nested braces preserved",
			input:   `{nested {interpolation}}`,
			wantCmd: `{nested {interpolation}}`,
			wantArg: []string{},
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
			name:    "mixed spacing",
			input:   `arg1   "arg  two"   {arg3}`,
			wantCmd: "arg1",
			wantArg: []string{"arg  two", "{arg3}"},
		},
		{
			name:    "empty",
			input:   ``,
			wantErr: true,
		},
		{
			name:    "inline brace group in unquoted token",
			input:   `--path={ some file }`,
			wantCmd: `--path={ some file }`,
			wantArg: []string{},
		},
		{
			name:    "brace group mid-token",
			input:   `hi{ ident }there`,
			wantCmd: `hi{ ident }there`,
			wantArg: []string{},
		},
		{
			name:    "multiple tokens with inline brace group",
			input:   `cmd --flag={ a b } tail`,
			wantCmd: `cmd`,
			wantArg: []string{`--flag={ a b }`, `tail`},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			parser := evaluator.NewArgumentParser(tt.input)
			ec, err := parser.ParseArguments()

			if tt.wantErr {
				if err == nil {
					t.Fatalf("input: %q, expected error, got none (cmd=%+v)", tt.input, ec)
				}
				return
			}

			if err != nil {
				t.Fatalf("input: %q, unexpected error: %v", tt.input, err)
			}

			if ec == nil {
				t.Fatalf("input: %q, expected ExecCommand, got nil", tt.input)
			}

			if ec.Command != tt.wantCmd {
				t.Fatalf("input: %q, command: expected %q, got %q", tt.input, tt.wantCmd, ec.Command)
			}

			if len(ec.Args) != len(tt.wantArg) {
				t.Fatalf("input: %q, expected %d args, got %d args (%v)",
					tt.input, len(tt.wantArg), len(ec.Args), ec.Args,
				)
			}

			for i := range tt.wantArg {
				if ec.Args[i] != tt.wantArg[i] {
					t.Fatalf("input: %q, arg[%d]: expected %q, got %q",
						tt.input, i, tt.wantArg[i], ec.Args[i],
					)
				}
			}
		})
	}
}
