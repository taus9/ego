package evaluator

import (
	"bytes"
	"ego/internal/backend/eval/object"
	"ego/internal/frontend/ast"
	"ego/internal/frontend/lexer"
	"ego/internal/frontend/parser"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ============================================================
// Import alias + path resolving
// ============================================================

func aliasForUse(u *ast.UseStatement) string {
	if u.Alias != nil {
		return u.Alias.Value
	}
	return defaultAliasFromSpec(u.Module.Value)
}

func defaultAliasFromSpec(spec string) string {
	// spec examples:
	// "math"
	// "dir/dir/file"
	// "./internal/utils"
	// "~/lib/utils"
	// "/home/user/lib/utils"
	s := strings.TrimSuffix(spec, ".ego")
	s = filepath.ToSlash(s)
	if i := strings.LastIndex(s, "/"); i >= 0 {
		s = s[i+1:]
	}
	return s
}

type Resolver struct {
	StdlibRoot string // e.g. /path/to/ego/stdlib
	HomeDir    string // e.g. /home/user
}

// ResolveImport returns an absolute, cleaned path to the module file.
// Rules:
// - no prefix => stdlib: <StdlibRoot>/<spec>.ego
// - "./" or "../" => relative to importing file directory
// - "~/" => relative to HomeDir
// - "/" => absolute path
func (r Resolver) ResolveImport(importerFile string, spec string) (string, error) {
	addExt := func(p string) string {
		if strings.HasSuffix(p, ".ego") {
			return p
		}
		return p + ".ego"
	}

	spec = filepath.ToSlash(spec)

	var p string
	switch {
	case strings.HasPrefix(spec, "./") || strings.HasPrefix(spec, "../") || spec == "." || spec == "..":
		// relative to the importing file's directory
		p = filepath.Join(filepath.Dir(importerFile), spec)

	case strings.HasPrefix(spec, "~/"):
		if r.HomeDir == "" {
			return "", fmt.Errorf("cannot resolve %q: home dir unknown", spec)
		}
		p = filepath.Join(r.HomeDir, spec[2:])

	case strings.HasPrefix(spec, "/"):
		p = spec

	default:
		// stdlib import
		if r.StdlibRoot == "" {
			return "", fmt.Errorf("cannot resolve %q: stdlib root not configured", spec)
		}
		p = filepath.Join(r.StdlibRoot, spec)
	}

	p = filepath.Clean(addExt(p))

	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(abs); err != nil {
		return "", fmt.Errorf("cannot import %q (resolved to %s): %w", spec, abs, err)
	}

	return abs, nil
}

type State int

const (
	Unseen State = iota
	Loading
	Loaded
)

type Loader struct {
	state      map[string]State
	parseCache map[string]*ast.Program
	envCache   map[string]*object.Environment
	GlobalEnv  *object.Environment
	resolver   Resolver
}

func NewLoader(res Resolver) *Loader {
	return &Loader{
		state:      map[string]State{},
		parseCache: map[string]*ast.Program{},
		envCache:   map[string]*object.Environment{},
		resolver:   res,
	}
}

// Load takes an absolute module path (recommended) and returns its evaluated Env.
// - If already Loaded: returns cached env (SKIP)
// - If Loading: circular dependency error
func (ld *Loader) Load(path string) (*object.Environment, error) {
	// Normalize path for stable cache keys
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	path = filepath.Clean(abs)

	switch ld.state[path] {
	case Loaded:
		return ld.envCache[path], nil
	case Loading:
		return nil, fmt.Errorf("circular dependency detected at %s", path)
	}

	ld.state[path] = Loading

	// Parse once (or reuse cached AST)
	program, ok := ld.parseCache[path]
	if !ok {

		src, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("error reading %s: %w", path, err)
		}

		l := lexer.New(string(src))
		p := parser.New(l)
		program = p.ParseProgram()

		//program = unwrapProtoMain(program)

		if p.ParseError() != nil {
			return nil, buildParserError(p.ParseError(), path)
		}

		ld.parseCache[path] = program
	}

	// Get all use statements in this module
	uses := ld.getUseStatements(program)

	// depsEnv holds imported module bindings for THIS module (math, m, utils, etc.)
	var depsEnv *object.Environment
	if ld.GlobalEnv != nil {
		// Use global env as outer if provided
		depsEnv = object.NewEnclosedEnvironment(ld.GlobalEnv)
	} else {
		// Or create a fresh one
		depsEnv = object.NewEnvironment()
	}

	// moduleEnv executes the module, with deps visible through outer chain
	moduleEnv := object.NewEnclosedEnvironment(depsEnv)

	// alias conflict detection inside THIS module
	aliasToResolved := map[string]string{}

	for _, u := range uses {
		spec := u.Module.Value
		alias := aliasForUse(u)

		resolved, err := ld.resolver.ResolveImport(path, spec)
		if err != nil {
			return nil, err
		}

		if prev, exists := aliasToResolved[alias]; exists && prev != resolved {
			return nil, fmt.Errorf("import alias conflict: %q refers to both %s and %s", alias, prev, resolved)
		}
		aliasToResolved[alias] = resolved

		// Load dependency (cached / skipped if already loaded)
		depModuleEnv, err := ld.Load(resolved)
		if err != nil {
			return nil, err
		}

		// Bind alias -> module object into *this module's* depsEnv
		modObj := &object.Module{Name: alias, Env: depModuleEnv}

		if !depsEnv.Declare(alias, modObj) {
			// This would only happen if the same alias appears twice in the same module.
			// If you want duplicates to be allowed when they resolve to same file, you can just continue.
			return nil, fmt.Errorf("import alias conflict: %q already declared in module %s", alias, path)
		}
	}

	// Evaluate module in moduleEnv (NO deps param needed)
	result := Eval(program, moduleEnv)
	if result != nil && result.Type() == object.UNHANDLED_ERROR_OBJ {
		return nil, buildRuntimeError(result.(*object.UnhandledError), path)
	}

	ld.envCache[path] = moduleEnv
	ld.state[path] = Loaded
	return moduleEnv, nil
}

func (ld *Loader) getUseStatements(program *ast.Program) []*ast.UseStatement {
	var uses []*ast.UseStatement
	for _, stmt := range program.Statements {
		if useStmt, ok := stmt.(*ast.UseStatement); ok {
			uses = append(uses, useStmt)
		} else {
			// Stop at the first non-use statement
			break
		}
	}

	return uses
}

func unwrapProtoMain(program *ast.Program) *ast.Program {
	protoMain, ok := program.Statements[0].(*ast.FunctionStatement)
	if !ok {
		return program
	}

	newProgram := &ast.Program{
		Statements: protoMain.Body.Statements,
	}

	return newProgram
}

func buildParserError(parseError *parser.ParseError, file string) error {
	var buf bytes.Buffer
	buf.WriteString("\tNope!\n")
	fmt.Fprintf(&buf, "\t%s", file)
	fmt.Fprintf(&buf, "\n\tParser Error:   %s", parseError.Message)

	span := parseError.Token.Span
	fmt.Fprintf(&buf, "\n\tToken Location: line %d, column %d", span.Line, span.Column)

	stack := parseError.StackTrace
	size := stack.Size()
	for i := size - 1; i >= 0; i-- {
		element := stack.Elements()[i]
		var msg string
		switch i {
		case size - 1:
			msg = "\n\tParser Trace:   -->"
		default:
			msg = "\t\t   "
		}

		buf.WriteString("\t")
		fmt.Fprintf(&buf, "%s %s", msg, stackTraceItemToString(element.(int)))
		buf.WriteString("\n")
	}

	return errors.New(buf.String())
}

func buildRuntimeError(runtimeError *object.UnhandledError, path string) error {
	var buf bytes.Buffer
	buf.WriteString("\tYikes!\n")
	fmt.Fprintf(&buf, "\t%s\n", path)
	span := runtimeError.Token.Span
	fmt.Fprintf(&buf, "\tToken Location: line %d, column %d\n", span.Line, span.Column)
	fmt.Fprintf(&buf, "\tRuntime Error:  %s", runtimeError.Message)
	return errors.New(buf.String())
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
	case parser.DECLARE:
		return "DECLARE"
	case parser.RETURN:
		return "RETURN"
	case parser.MAP:
		return "MAP"
	case parser.CALL_ARGS:
		return "FUNCTION CALL ARGUMENTS"
	case parser.FUNC:
		return "FUNCTION"
	case parser.FOR:
		return "FOR LOOP"
	case parser.ASSIGN:
		return "ASSIGNMENT"
	case parser.WHEN:
		return "WHEN"
	case parser.IS:
		return "IS"
	case parser.ELSE:
		return "ELSE"
	case parser.USE:
		return "USE"
	default:
		return "UNKNOWN"
	}
}
