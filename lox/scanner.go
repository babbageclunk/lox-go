package lox

import "strconv"

// import java.util.ArrayList;
// import java.util.HashMap;
// import java.util.List;
// import java.util.Map;

// import static com.craftinginterpreters.lox.TokenType.*;

// class Scanner {
//   private final String source;
//   private final List<Token> tokens = new ArrayList<>();
//   private int start = 0;
//   private int current = 0;
//   private int line = 1;

//   Scanner(String source) {
//     this.source = source;
//   }
// }

type Scanner struct {
	source  []rune
	tokens  []Token
	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{source: []rune(source), line: 1}
}

// List<Token> scanTokens() {
//   while (!isAtEnd()) {
//     // We are at the beginning of the next lexeme.
//     start = current;
//     scanToken();
//   }

//   tokens.add(new Token(EOF, "", null, line));
//   return tokens;
// }

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		// We're at the beginning of the next lexeme.
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, NewToken(TokenEof, "", nil, s.line))
	return s.tokens
}

//	private boolean isAtEnd() {
//	  return current >= source.length();
//	}
func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

// private void scanToken() {
//   char c = advance();
//   switch (c) {
//     case '(': addToken(LEFT_PAREN); break;
//     case ')': addToken(RIGHT_PAREN); break;
//     case '{': addToken(LEFT_BRACE); break;
//     case '}': addToken(RIGHT_BRACE); break;
//     case ',': addToken(COMMA); break;
//     case '.': addToken(DOT); break;
//     case '-': addToken(MINUS); break;
//     case '+': addToken(PLUS); break;
//     case ';': addToken(SEMICOLON); break;
//     case '*': addToken(STAR); break;
//     case '!':
//       addToken(match('=') ? BANG_EQUAL : BANG);
//       break;
//     case '=':
//       addToken(match('=') ? EQUAL_EQUAL : EQUAL);
//       break;
//     case '<':
//       addToken(match('=') ? LESS_EQUAL : LESS);
//       break;
//     case '>':
//       addToken(match('=') ? GREATER_EQUAL : GREATER);
//       break;
//     case '/':
//       if (match('/')) {
//         // A comment goes until the end of the line.
//         while (peek() != '\n' && !isAtEnd()) advance();
//       } else {
//         addToken(SLASH);
//       }
//       break;
//
//     case ' ':
//     case '\r':
//     case '\t':
//       // Ignore whitespace.
//       break;
//
//     case '\n':
//       line++;
//       break;
//
//     case '"': string(); break;
//
//     default:
//       if (isDigit(c)) {
//         number();
//       } else if (isAlpha(c)) {
//         identifier();
//       } else {
//         Lox.error(line, "Unexpected character.");
//       }
//       break;
//   }
// }

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(TokenLeftParen)
	case ')':
		s.addToken(TokenRightParen)
	case '{':
		s.addToken(TokenLeftBrace)
	case '}':
		s.addToken(TokenRightBrace)
	case ',':
		s.addToken(TokenComma)
	case '.':
		s.addToken(TokenDot)
	case '-':
		s.addToken(TokenMinus)
	case '+':
		s.addToken(TokenPlus)
	case ';':
		s.addToken(TokenSemicolon)
	case '*':
		s.addToken(TokenStar)

	case '!':
		tt := TokenBang
		if s.match('=') {
			tt = TokenBangEqual
		}
		s.addToken(tt)
	case '=':
		tt := TokenEqual
		if s.match('=') {
			tt = TokenEqualEqual
		}
		s.addToken(tt)
	case '<':
		tt := TokenLess
		if s.match('=') {
			tt = TokenLessEqual
		}
		s.addToken(tt)
	case '>':
		tt := TokenGreater
		if s.match('=') {
			tt = TokenGreaterEqual
		}
		s.addToken(tt)

	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(TokenSlash)
		}

	case ' ', '\r', '\t':
		// Ignore whitespace.
	case '\n':
		s.line++

	case '"':
		s.string()

	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			report(s.line, "Unexpected character.")
		}
	}
}

// private void identifier() {
//   while (isAlphaNumeric(peek())) advance();
//   String text = source.substring(start, current);
//   TokenType type = keywords.get(text);
//   if (type == null) type = IDENTIFIER;
//   addToken(type);
// }

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	tt, found := keywords[string(text)]
	if !found {
		tt = TokenIdentifier
	}
	s.addToken(tt)
}

// private void number() {
//   while (isDigit(peek())) advance();

//   // Look for a fractional part.
//   if (peek() == '.' && isDigit(peekNext())) {
//     // Consume the "."
//     advance();

//     while (isDigit(peek())) advance();
//   }

//   addToken(NUMBER,
//       Double.parseDouble(source.substring(start, current)));
// }

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part.
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	value, err := strconv.ParseFloat(string(s.source[s.start:s.current]), 64)
	if err != nil {
		report(s.line, "parsing float: "+err.Error())
		return
	}
	s.addTokenLit(TokenNumber, value)
}

//   private void string() {
//     while (peek() != '"' && !isAtEnd()) {
//       if (peek() == '\n') line++;
//       advance();
//     }

//     if (isAtEnd()) {
//       Lox.error(line, "Unterminated string.");
//       return;
//     }

//     // The closing ".
//     advance();

//	  // Trim the surrounding quotes.
//	  String value = source.substring(start + 1, current - 1);
//	  addToken(STRING, value);
//	}
func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		report(s.line, "Unterminated string.")
		return
	}

	// The closing ".
	s.advance()

	// Trim the surrounding quotes/
	value := s.source[s.start+1 : s.current-1]
	s.addTokenLit(TokenString, string(value))
}

// private boolean match(char expected) {
//   if (isAtEnd()) return false;
//   if (source.charAt(current) != expected) return false;

//   current++;
//   return true;
// }

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

// private char peek() {
//   if (isAtEnd()) return '\0';
//   return source.charAt(current);
// }

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

//	private char peekNext() {
//	  if (current + 1 >= source.length()) return '\0';
//	  return source.charAt(current + 1);
//	}
func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

// private boolean isAlpha(char c) {
//   return (c >= 'a' && c <= 'z') ||
//          (c >= 'A' && c <= 'Z') ||
//           c == '_';
// }

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

// private boolean isAlphaNumeric(char c) {
//   return isAlpha(c) || isDigit(c);
// }

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

// private boolean isDigit(char c) {
//   return c >= '0' && c <= '9';
// }

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

// private char advance() {
//   return source.charAt(current++);
// }

func (s *Scanner) advance() rune {
	r := s.source[s.current]
	s.current++
	return r
}

// private void addToken(TokenType type) {
//   addToken(type, null);
// }

func (s *Scanner) addToken(tt TokenType) {
	s.addTokenLit(tt, nil)
}

// private void addToken(TokenType type, Object literal) {
//   String text = source.substring(start, current);
//   tokens.add(new Token(type, text, literal, line));
// }

func (s *Scanner) addTokenLit(tt TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tt, string(text), literal, s.line))
}
