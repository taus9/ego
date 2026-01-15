# Ego Programming Language

Ego is a small interpreted programming language implemented in Go. It includes a lexer, parser, AST, evaluator, REPL, and a simple CLI so you can write and run Ego source files.

This project is designed as a learning and exploration language, with a focus on clear architecture and readable implementation rather than raw performance.

## Features

- Expression-based language with integers, floats, booleans, strings, and `nil`
- Variables and immutable bindings via assignment (e.g. `x = 5`)
- Arithmetic and comparison operators: `+`, `-`, `*`, `/`, `%`, `==`, `!=`, `<`, `>`, `<=`, `>=`
- Logical operators: `and`, `or`, `!`
- `if` / `else` expressions
- `for` loops:
  - `for i = 0 to 10` style loops
  - `for condition` style while-loops
- Functions:
  - Named functions
  - Anonymous functions using `:>` syntax
  - First-class function values and function calls
- Arrays and maps with index access
- Simple standard builtins (e.g. `put(...)` for output)
- REPL with multi-line support and basic parse error reporting

## Examples

A simple FizzBuzz implementation in Ego (see `examples/fizzbuzz.ego`):

```ego
:fizzbuzz n
    for i = 1 to n 
        printed = false
        if i % 15 == 0
            put('FizzBuzz')
            printed = true
        ;
        if i % 3 == 0 and !printed
            put('Fizz')
            printed = true
        ;
        if i % 5 == 0 and !printed
            put('Buzz')
            printed = true
        ;
        if !printed
            put(i)
        ;
    ;
;

fizzbuzz(100)
```

More examples live in the `examples/` directory:

- `examples/fib.ego` – Fibonacci sequence
- `examples/fizzbuzz.ego` – FizzBuzz
- `examples/for.ego` – Loop constructs
- `examples/test.ego` – Misc test snippets

## Getting Started

### Prerequisites

- Go toolchain installed (1.22+ recommended)

### Build & Run

Clone the repository and run Ego programs via the CLI:

```bash
# Run an Ego file
go run main.go examples/fizzbuzz.ego

# Or point to any .ego file
go run main.go path/to/program.ego
```

You can also start an interactive REPL:

```bash
# From within the repo
go run main.go
```

(In this setup, the REPL is available as an internal package at `internal/repl`, and can be wired into a CLI entrypoint as needed.)

### Command-line usage

The CLI supports a couple of options:

```text
ego programming language

Usage:
  ego [options] <file>

Options:
  -v      Show version information
  -h      Show this help message

Arguments:
  <file>  Path to the Ego source file to execute
```

## Architecture Overview

The project is organized into clear frontend and backend layers under `internal/`:

- `internal/frontend/token` – token types and span information
- `internal/frontend/lexer` – converts source text into a stream of tokens, tracking line/column spans
- `internal/frontend/ast` – abstract syntax tree node definitions for expressions, statements, and programs
- `internal/frontend/parser` – recursive-descent Pratt parser that builds ASTs, with operator precedence and a small stack-trace mechanism for better parse errors
- `internal/backend/eval/object` – runtime object system (integers, floats, booleans, strings, arrays, maps, functions, errors, environments, etc.)
- `internal/backend/eval/evaluator` – tree-walking evaluator that executes AST nodes in an environment, including control flow, functions, and builtins
- `internal/repl` – line-based REPL that keeps a shared environment, supports multi-line blocks, and shows parse errors

The top-level `main.go` wires everything together into a simple CLI: it reads a file, lexes/parses it into an AST, evaluates it in a fresh environment, and prints either the result or a detailed parse error with stack trace.

## Testing

The core components are covered by unit tests:

- Lexer tests under `internal/frontend/lexer/test` (tokenization, spans, strings)
- Parser tests under `internal/frontend/parser/test` (expressions, statements, precedence, functions, arrays/maps, errors)
- Evaluator tests under `internal/backend/eval/evaluator/test` (expressions, conditionals, loops, functions, errors, arrays/maps)
- Object tests under `internal/backend/eval/object/test` (e.g. string map keys)

Run the full test suite with:

```bash
go test ./...
```

## Background & Motivation

Ego started as an exploration of Thorston Ball’s book *Writing an Interpreter in Go* (really good).

The project is intended as a readable, educational codebase that demonstrates:

- Lexing and parsing with clear error reporting
- Building and walking an AST
- Implementing environments and lexical scoping
- Designing a small but expressive language with a minimal set of constructs

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
