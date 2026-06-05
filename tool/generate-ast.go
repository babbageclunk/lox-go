package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rotisserie/eris"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: generate-ast <output path without .go>")
		os.Exit(64)
	}
	baseFilename := os.Args[1]
	err := defineAst(baseFilename+"-expr.go", "Expr", []string{
		"Assign   : Name Token, Value Expr",
		"Binary   : Left Expr, Operator Token, Right Expr",
		"Call     : Callee Expr, Paren Token, Arguments []Expr",
		"Grouping : Expression Expr",
		"Literal  : Value any",
		"Logical  : Left Expr, Operator Token, Right Expr",
		"Unary    : Operator Token, Right Expr",
		"Variable : Name Token",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(65)
	}
	err = defineAst(baseFilename+"-stmt.go", "Stmt", []string{
		"Block      : Statements []Stmt",
		"Break      : ",
		"Expression : Expression Expr",
		"Function   : Name Token, Params []Token, Body []Stmt",
		"If         : Condition Expr, ThenBranch Stmt, ElseBranch Stmt",
		"Print      : Expression Expr",
		"Var        : Name Token, Initializer Expr",
		"While      : Condition Expr, Body Stmt",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(65)
	}
}

type Type struct {
	name   string
	fields []Field
}

type Field struct {
	name     string
	typeName string
}

func parseTypes(baseName string, lines []string) ([]Type, error) {
	types := make([]Type, len(lines))
	for i, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, eris.Errorf("malformed type line %d: %q", i+1, line)
		}
		typeName := strings.TrimSpace(parts[0]) + baseName
		fieldStr := strings.TrimSpace(parts[1])
		if len(fieldStr) == 0 {
			types[i] = Type{name: typeName}
			continue
		}

		fieldLines := strings.Split(fieldStr, ",")
		fields := make([]Field, len(fieldLines))
		for j, f := range fieldLines {
			parts := strings.SplitN(strings.TrimSpace(f), " ", 2)
			if len(parts) != 2 {
				return nil, eris.Errorf("malformed field  %d: %q", i+1, f)
			}
			fields[j] = Field{name: parts[0], typeName: parts[1]}
		}
		types[i] = Type{name: typeName, fields: fields}
	}
	return types, nil
}

func defineAst(path string, baseName string, lines []string) error {
	types, err := parseTypes(baseName, lines)
	if err != nil {
		return eris.Wrapf(err, "parsing %s types", baseName)
	}

	outf, err := os.Create(path)
	if err != nil {
		return eris.Wrapf(err, "creating %s", err)
	}
	defer outf.Close()
	w := bufio.NewWriter(outf)
	defer w.Flush()
	fmt.Fprintln(w, "package lox")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "import \"fmt\"")
	fmt.Fprintln(w)
	defineVisitor(w, baseName, types)
	fmt.Fprintln(w)
	for _, t := range types {
		defineType(w, baseName, t.name, t.fields)
	}
	defineAsAcceptor(w, baseName, types)
	return nil
}

func defineVisitor(w io.Writer, baseName string, types []Type) {
	fmt.Fprintf(w, "type %sVisitor[R any] interface {\n", baseName)
	for _, t := range types {
		fmt.Fprintf(w, "\tVisit%s(%s %s) (R, error)\n",
			t.name, strings.ToLower(baseName), t.name)
	}
	fmt.Fprintln(w, "}")
}

func defineType(w io.Writer, baseName, typeName string, fields []Field) {
	fmt.Fprintf(w, "type %s struct {\n", typeName)
	for _, field := range fields {
		fmt.Fprintf(w, "\t%s %s\n", field.name, field.typeName)
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "func (%s) %sKind() string {\n", typeName, strings.ToLower(baseName[:1]))
	fmt.Fprintf(w, "\treturn \"%s\"\n", typeName)
	fmt.Fprintln(w, "}")

	// This is a deviation from the Java implementation. Because Go
	// doesn't allow methods to introduce new type parameters, we need
	// to define a new generic type wrapping each node type to
	// implement the generic Acceptor interface, and then the generic
	// asAcceptor function that converts a given Expr into the
	// acceptor for that Expr type and return type (see
	// defineAsAcceptor).
	fmt.Fprintf(w, `
type %[1]sAcceptor[R any] %[1]s

func (%[2]s %[1]sAcceptor[R]) accept(vis %[3]sVisitor[R]) (R, error) {
	return vis.Visit%[1]s(%[1]s(%[2]s))
}
`, typeName, strings.ToLower(typeName[:1]), baseName)
	fmt.Fprintln(w)
}

func defineAsAcceptor(w io.Writer, baseName string, types []Type) {
	varName := strings.ToLower(baseName)
	fmt.Fprintf(w, "func as%[2]sAcceptor[R any](%[1]s %[2]s) %[2]sAcceptor[R] {\n", varName, baseName)
	fmt.Fprintf(w, "\tswitch e := %s.(type) {\n", varName)
	for _, t := range types {
		fmt.Fprintf(w, "\tcase %s:\n", t.name)
		fmt.Fprintf(w, "\t\treturn %sAcceptor[R](e)\n", t.name)
	}
	fmt.Fprintln(w, "\t}")
	fmt.Fprintf(w, "\tpanic(fmt.Errorf(\"no acceptor for %s %%s\", %s.%sKind()))\n", varName, varName, strings.ToLower(baseName[:1]))
	fmt.Fprintln(w, "}")
}
