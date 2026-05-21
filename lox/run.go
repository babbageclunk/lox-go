package lox

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/rotisserie/eris"
)

// private static void runFile(String path) throws IOException {
//   byte[] bytes = Files.readAllBytes(Paths.get(path));
//   run(new String(bytes, Charset.defaultCharset()));
//
//   // Indicate an error in the exit code.
//   if (hadError) System.exit(65);
// }

func RunFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}
	scanner := NewScanner(string(bytes))
	tokens, err := scanner.ScanTokens()
	if err != nil {
		runtimeError(err)
		os.Exit(65)
	}
	parser := NewParser(tokens)
	statements, err := parser.parse()
	// Stop if there was a syntax error.
	if err != nil {
		runtimeError(err)
		os.Exit(65)
	}

	err = interpreter.Interpret(statements)
	if err != nil {
		runtimeError(err)
		os.Exit(70)
	}
	return nil
}

// private static void runPrompt() throws IOException {
//   InputStreamReader input = new InputStreamReader(System.in);
//   BufferedReader reader = new BufferedReader(input);

//   for (;;) {
//     System.out.print("> ");
//     String line = reader.readLine();
//     if (line == null) break;
//     run(line);
//     hadError = false;
//   }
// }

func RunPrompt() error {
	input := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		bytes, _, err := input.ReadLine()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("reading line: %w", err)
		}
		run(string(bytes))
	}
}

// private static void run(String source) {
//   Scanner scanner = new Scanner(source);
//   List<Token> tokens = scanner.scanTokens();
//   Parser parser = new Parser(tokens);
//   Expr expression = parser.parse();
//   // Stop if there was a syntax error.
//   if (hadError) return;
//
//   System.out.println(new AstPrinter().print(expression));
// }

var interpreter = NewInterpreter()

func run(source string) {
	scanner := NewScanner(source)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	parser := NewParser(tokens)
	expr, ok := parser.safeExpression()
	if ok {
		val, err := interpreter.Evaluate(expr)
		if err != nil {
			runtimeError(err)
			return
		}
		fmt.Println(stringify(val))
		return
	}
	statements, err := parser.parse()
	// Stop if there was a syntax error.
	if err != nil {
		runtimeError(err)
		return
	}

	err = interpreter.Interpret(statements)
	if err != nil {
		runtimeError(err)
	}
}

// static void error(int line, String message) {
//   report(line, "", message);
// }

func reportError(line int, where, message string) error {
	return fmt.Errorf("[line %d] Error%s: %s", line, where, message)
}

func runtimeError(err error) {
	var tErr tokenError
	if eris.As(err, &tErr) {
		fmt.Fprintf(os.Stderr, "[line %d] Error: %s\n", tErr.token.Line, err.Error())
	} else {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

// static void error(Token token, String message) {
//   if (token.type == TokenType.EOF) {
//     report(token.line, " at end", message);
//   } else {
//     report(token.line, " at '" + token.lexeme + "'", message);
//   }
// }

func loxError(token Token, message string) error {
	var where string
	if token.Type == TokenEof {
		where = " at end"
	} else {
		where = fmt.Sprintf(" at %q", token.Lexeme)
	}
	return reportError(token.Line, where, message)
}
