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
	err := defineAst(outputPath, []string{
		"Binary   : left Expr, operator Token, right Expr",
		"Grouping : expression Expr",
		"Literal  : value any",
		"Unary    : operator Token, right Expr",
	})
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
//    // The AST classes.
//    for (String type : types) {
//      String className = type.split(":")[0].trim();
//      String fields = type.split(":")[1].trim();
//      defineType(writer, baseName, className, fields);
//    }
//    writer.println("}");
//    writer.close();
//  }

func defineAst(path string, types []string) error {
	outf, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating %s: %w", path, err)
	}
	defer outf.Close()
	w := bufio.NewWriter(outf)
	defer w.Flush()
	fmt.Fprintln(w, "package lox")
	fmt.Fprintln(w)
	for i, line := range types {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("malformed type line %d: %q", i+1, line)
		}
		typeName := strings.TrimSpace(parts[0])
		fields := strings.Split(strings.TrimSpace(parts[1]), ",")
		if err := defineType(w, typeName, fields); err != nil {
			return fmt.Errorf("defining type %s: %w", typeName, err)
		}
	}
	return nil
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

func defineType(w io.Writer, typeName string, fields []string) error {
	fmt.Fprintf(w, "type %s struct {\n", typeName)
	for _, field := range fields {
		fmt.Fprintf(w, "\t%s\n", strings.TrimSpace(field))
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "func (%s) kind() string {\n", typeName)
	fmt.Fprintf(w, "\treturn \"%s\"\n", typeName)
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
	return nil
}
