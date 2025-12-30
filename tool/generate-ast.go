package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// package com.craftinginterpreters.tool;

// import java.io.IOException;
// import java.io.PrintWriter;
// import java.util.Arrays;
// import java.util.List;

// public class GenerateAst {
//   public static void main(String[] args) throws IOException {
//     if (args.length != 1) {
//       System.err.println("Usage: generate_ast <output directory>");
//       System.exit(64);
//     }
//     String outputDir = args[0];
//     defineAst(outputDir, "Expr", Arrays.asList(
//       "Binary   : Expr left, Token operator, Expr right",
//       "Grouping : Expr expression",
//       "Literal  : Object value",
//       "Unary    : Token operator, Expr right"
//     ));
//   }
// }

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: generate-ast <output path without .go>")
		os.Exit(64)
	}
	outputPath := os.Args[1] + ".go"
	types, err := parseTypes([]string{
		"Binary   : Left Expr, Operator Token, Right Expr",
		"Grouping : Expression Expr",
		"Literal  : Value any",
		"Unary    : Operator Token, Right Expr",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(65)
	}
	err = defineAst(outputPath, "Expr", types)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(65)
	}
}

//  private static void defineAst(
//      String outputDir, String baseName, List<String> types)
//      throws IOException {
//    String path = outputDir + "/" + baseName + ".java";
//    PrintWriter writer = new PrintWriter(path, "UTF-8");
//
//    writer.println("package com.craftinginterpreters.lox;");
//    writer.println();
//    writer.println("import java.util.List;");
//    writer.println();
//    writer.println("abstract class " + baseName + " {");
//
//    defineVisitor(writer, baseName, types);
//
//    // The AST classes.
//    for (String type : types) {
//      String className = type.split(":")[0].trim();
//      String fields = type.split(":")[1].trim();
//      defineType(writer, baseName, className, fields);
//    }
//    writer.println("}");
//    writer.close();
//  }

type Type struct {
	name   string
	fields []Field
}

type Field struct {
	name     string
	typeName string
}

func parseTypes(lines []string) ([]Type, error) {
	types := make([]Type, len(lines))
	for i, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed type line %d: %q", i+1, line)
		}
		typeName := strings.TrimSpace(parts[0])
		fieldLines := strings.Split(parts[1], ",")
		fields := make([]Field, len(fieldLines))
		for j, f := range fieldLines {
			parts := strings.SplitN(strings.TrimSpace(f), " ", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("malformed field  %d: %q", i+1, f)
			}
			fields[j] = Field{name: parts[0], typeName: parts[1]}
		}
		types[i] = Type{name: typeName, fields: fields}
	}
	return types, nil
}

func defineAst(path string, baseName string, types []Type) error {
	outf, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating %s: %w", path, err)
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

//  private static void defineVisitor(
//      PrintWriter writer, String baseName, List<String> types) {
//    writer.println("  interface Visitor<R> {");
//
//    for (String type : types) {
//      String typeName = type.split(":")[0].trim();
//      writer.println("    R visit" + typeName + baseName + "(" +
//          typeName + " " + baseName.toLowerCase() + ");");
//    }
//
//    writer.println("  }");
//  }

func defineVisitor(w io.Writer, baseName string, types []Type) {
	fmt.Fprintln(w, "type Visitor[R any] interface {")
	for _, t := range types {
		fmt.Fprintf(w, "\tvisit%s%s(%s %s) R\n",
			t.name, baseName, strings.ToLower(baseName), t.name)
	}
	fmt.Fprintln(w, "}")
}

//  private static void defineType(
//      PrintWriter writer, String baseName,
//      String className, String fieldList) {
//    writer.println("  static class " + className + " extends " +
//        baseName + " {");
//
//    // Constructor.
//    writer.println("    " + className + "(" + fieldList + ") {");
//
//    // Store parameters in fields.
//    String[] fields = fieldList.split(", ");
//    for (String field : fields) {
//      String name = field.split(" ")[1];
//      writer.println("      this." + name + " = " + name + ";");
//    }
//
//    writer.println("    }");
//
//    // Fields.
//    writer.println();
//    for (String field : fields) {
//      writer.println("    final " + field + ";");
//    }
//
//    writer.println("  }");
//  }

func defineType(w io.Writer, baseName, typeName string, fields []Field) {
	fmt.Fprintf(w, "type %s struct {\n", typeName)
	for _, field := range fields {
		fmt.Fprintf(w, "\t%s %s\n", field.name, field.typeName)
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "func (%s) kind() string {\n", typeName)
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

func (%[2]s %[1]sAcceptor[R]) accept(v Visitor[R]) R {
	return v.visit%[1]s%[3]s(%[1]s(%[2]s))
}
`, typeName, strings.ToLower(typeName[:1]), baseName)
	fmt.Fprintln(w)
}

func defineAsAcceptor(w io.Writer, baseName string, types []Type) {
	varName := strings.ToLower(baseName)
	fmt.Fprintf(w, "func asAcceptor[R any](%s %s) Acceptor[R] {\n", varName, baseName)
	fmt.Fprintf(w, "\tswitch e := %s.(type) {\n", varName)
	for _, t := range types {
		fmt.Fprintf(w, "\tcase %s:\n", t.name)
		fmt.Fprintf(w, "\t\treturn %sAcceptor[R](e)\n", t.name)
	}
	fmt.Fprintln(w, "\t}")
	fmt.Fprintf(w, "\tpanic(fmt.Errorf(\"no acceptor for %s %%s\", %s.kind()))\n", varName, varName)
	fmt.Fprintln(w, "}")
}
