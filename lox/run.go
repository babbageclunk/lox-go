package lox

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
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
	run(string(bytes))
	if hadError {
		os.Exit(65)
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
		hadError = false
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

func run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()
	parser := NewParser(tokens)
	expr := parser.parse()
	// Stop if there was a syntax error.
	if hadError {
		return
	}
	res, err := Print(expr)
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Println(res)
	}
}

// static void error(int line, String message) {
//   report(line, "", message);
// }

// private static void report(int line, String where,
//                            String message) {
//   System.err.println(
//       "[line " + line + "] Error" + where + ": " + message);
//   hadError = true;
// }

func report(line int, message string) {
	reportError(line, "", message)
}

var hadError = false

func reportError(line int, where, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
	hadError = true
}

// static void error(Token token, String message) {
//   if (token.type == TokenType.EOF) {
//     report(token.line, " at end", message);
//   } else {
//     report(token.line, " at '" + token.lexeme + "'", message);
//   }
// }

func loxError(token Token, message string) {
	if token.Type == TokenEof {
		reportError(token.Line, " at end", message)
	} else {
		reportError(token.Line, fmt.Sprintf(" at %q", token.Lexeme), message)
	}
}
